/*
Copyright (c) 2013 Nicolas Thery <nthery@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

/*
Package raytracer implements a very simple ray tracer based the algorithms described in:

	http://www.ccs.neu.edu/home/fell/CSU540/programs/RayTracingFormulas.htm
*/
package raytracer

import (
	"fmt"
	"github.com/nthery/goraytracer/geom"
	"image"
	"image/color"
	"math"
)

// A Color is a red/green/blue triplet of color channels in [0..1] range
type Color struct {
	R, G, B float64
}

func (c *Color) Validate() error {
	if isColorChannelValid(c.R) && isColorChannelValid(c.G) && isColorChannelValid(c.B) {
		return nil
	}
	return fmt.Errorf("color out-of-range: %#v", c)
}

// toRGBA converts to standard 32bpp.
// The color.Color interface is not used for performance.
func (c *Color) toRGBA() color.RGBA {
	return color.RGBA{uint8(c.R * 255), uint8(c.G * 255), uint8(c.B * 255), 255}
}

func isColorChannelValid(c float64) bool {
	return 0 <= c && c <= 1
}

// Sphere objects are part of the scene to render.
type Sphere struct {
	// No embedding here for compatibility with json package
	Sphere geom.Sphere
	Color  Color
}

func (s *Sphere) Validate() error {
	if err := s.Sphere.Validate(); err != nil {
		return err
	}
	if err := s.Color.Validate(); err != nil {
		return fmt.Errorf("invalid sphere: %v", err)
	}
	return nil
}

// A Frustum is a pyramidal viewing frustum orthogonal to the z-axis.  The
// rendered scene is projected onto the near plane.  The size ratio between the
// near and far planes determines the field of view.
type Frustum struct {
	Near, Far geom.Plane2d // near and far viewing planes
}

func (f *Frustum) Validate() error {
	if err := f.Near.Validate(); err != nil {
		return fmt.Errorf("invalid near frustum plane: %v", err)
	}
	if err := f.Far.Validate(); err != nil {
		return fmt.Errorf("invalid far frustum plane: %v", err)
	}
	return nil
}

// The Scene to render.
type Scene struct {
	ViewFrustum Frustum
	Light       geom.Point // coordinate of light source
	Objects     []Sphere   // objects to render
	Bg          Color      // background color
	Kd          float64    // diffuse coefficient
}

func (s *Scene) Validate() error {
	for _, o := range s.Objects {
		if err := o.Validate(); err != nil {
			return fmt.Errorf("invalid scene object: %v", err)
		}
	}
	if err := s.Bg.Validate(); err != nil {
		return fmt.Errorf("invalid scene background: %v", err)
	}
	if s.Kd < 0 || s.Kd > 1 {
		return fmt.Errorf("invalid scene diffuse coefficient: %v", s.Kd)
	}
	if err := s.ViewFrustum.Validate(); err != nil {
		return fmt.Errorf("invalid scene frustum: %v", err)
	}
	return nil
}

// diffuseShading computes a color channel value taking into account the
// diffuse and ambiant coefficients and the angle between the pixel and light
// source.
func diffuseShading(factor, kd float64, channel float64) float64 {
	ka := 1 - kd
	return factor*kd*channel + factor*ka
}

// rayHitsObject returns whether the ray intersects one object in the scene.
func (s *Scene) rayHitsObject(ray geom.Line) bool {
	for i := range s.Objects {
		_, _, ok := geom.SphereLineIntersection(s.Objects[i].Sphere, ray)
		if ok {
			return true
		}
	}
	return false
}

// castRay finds the nearest intersection point between the ray and the scene
// objects.  On return, obj is nil if there is no intersection.
func (s *Scene) castRay(ray geom.Line) (obj *Sphere, intersection geom.Point) {
	var pmin geom.Point
	tmin := math.MaxFloat64
	imin := -1
	for i := range s.Objects {
		p, t, ok := geom.SphereLineIntersection(s.Objects[i].Sphere, ray)
		if ok {
			if t < tmin {
				tmin = t
				pmin = p
				imin = i
			}
		}
	}

	if imin == -1 {
		return nil, geom.Origin
	}

	return &s.Objects[imin], pmin
}

func (s *Scene) computeObjectColorAt(obj *Sphere, p geom.Point) Color {
	normal := obj.Sphere.NormalVectorAt(&p)
	light := geom.MakeVector(s.Light, p)
	light = light.UnitVector()
	dot := geom.DotProduct(&light, &normal)
	if dot < 0 {
		dot = 0
	}
	r := diffuseShading(dot, s.Kd, obj.Color.R)
	g := diffuseShading(dot, s.Kd, obj.Color.G)
	b := diffuseShading(dot, s.Kd, obj.Color.B)

	return Color{r, g, b}
}

func bgShadowPixel(c Color) Color {
	return Color{c.R / 2, c.G / 2, c.B / 2}
}

func (s *Scene) renderPixel(x, y float64) color.RGBA {
	xfar := x * s.ViewFrustum.Far.Dx() / s.ViewFrustum.Near.Dx()
	yfar := y * s.ViewFrustum.Far.Dy() / s.ViewFrustum.Near.Dx()
	ray := geom.Line{
		geom.Point{x, y, s.ViewFrustum.Near.Z},
		geom.Point{xfar, yfar, s.ViewFrustum.Far.Z},
	}

	obj, intersection := s.castRay(ray)

	var c Color
	if obj != nil {
		// Is intersection shadowed by another object?
		sray := geom.Line{s.Light, intersection}
		other, _ := s.castRay(sray)
		if other != nil && other != obj {
			c = Color{
				(1 - s.Kd) * obj.Color.R,
				(1 - s.Kd) * obj.Color.G,
				(1 - s.Kd) * obj.Color.B,
			}
		} else {
			c = s.computeObjectColorAt(obj, intersection)
		}
	} else {
		sray := geom.Line{
			geom.Point{xfar, yfar, s.ViewFrustum.Far.Z},
			s.Light,
		}
		if s.rayHitsObject(sray) {
			c = Color{s.Bg.R / 2, s.Bg.G / 2, s.Bg.B / 2}
		} else {
			c = s.Bg
		}
	}

	return c.toRGBA()
}

// Render validates the scene and runs the ray-tracing algorithm over it.  It
// generates an in-memory image containing the result.  The scene is divided in
// nstripes horizontal stripes that are processed concurrently.
func (s *Scene) Render(nstripes int) (*image.RGBA, error) {
	if nstripes < 1 {
		nstripes = 1
	}

	if err := s.Validate(); err != nil {
		return nil, err
	}

	vp := &s.ViewFrustum.Near
	w := int(vp.Dx())
	h := int(vp.Dy())
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	slice := h / nstripes
	ch := make(chan bool)
	for n := 0; n < nstripes; n++ {
		ystart := slice * n
		yend := ystart + slice
		go func() {
			for y := ystart; y < yend; y++ {
				for x := 0; x < w; x++ {
					c := s.renderPixel(float64(x)+vp.Tl.X, -vp.Br.Y-float64(y))
					img.SetRGBA(x, y, c)
				}
			}
			ch <- true
		}()
	}

	// block until all stripes processed
	for n := 0; n < nstripes; n++ {
		<-ch
	}
	return img, nil
}
