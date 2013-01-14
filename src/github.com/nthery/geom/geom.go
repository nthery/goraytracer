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

// Package geom defines some geometry primitives
package geom

import (
	"fmt"
	"math"
)

// A Point2d lies on a plane orthogonal to the z-axis
type Point2d struct {
	X, Y float64
}

// A Plane2d is a bounded plane orthogonal to the z-axis
type Plane2d struct {
	Tl, Br Point2d // top-left and bottom-right corners
	Z      float64
}

func (p *Plane2d) Validate() error {
	if p.Tl.X >= p.Br.X || p.Tl.Y <= p.Br.Y {
		return fmt.Errorf("negative or null plane width or height: %#v", p)
	}
	return nil
}

func (p *Plane2d) Dx() float64 {
	return float64(p.Br.X - p.Tl.X)
}

func (p *Plane2d) Dy() float64 {
	return float64(p.Tl.Y - p.Br.Y)
}

// A Point is a 3-dimensional point
type Point struct {
	X, Y, Z float64
}

func PointsEqual(lhs, rhs Point, epsilon float64) bool {
	return FloatsEqual(lhs.X, rhs.X, epsilon) &&
		FloatsEqual(lhs.Y, rhs.Y, epsilon) &&
		FloatsEqual(lhs.Z, rhs.Z, epsilon)
}

var Origin = Point{0, 0, 0}

// A Line is a line segment
type Line [2]Point

type Vector struct {
	X, Y, Z float64
}

func MakeVector(head, tail Point) Vector {
	return Vector{head.X - tail.X, head.Y - tail.Y, head.Z - tail.Z}
}

func (v *Vector) UnitVector() Vector {
	m := v.Module()
	return Vector{v.X / m, v.Y / m, v.Z / m}
}

func (v *Vector) Module() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func DotProduct(lhs, rhs *Vector) float64 {
	return lhs.X*rhs.X + lhs.Y*rhs.Y + lhs.Z*rhs.Z
}

func VectorsEqual(lhs, rhs Vector, epsilon float64) bool {
	return FloatsEqual(lhs.X, rhs.X, epsilon) &&
		FloatsEqual(lhs.Y, rhs.Y, epsilon) &&
		FloatsEqual(lhs.Z, rhs.Z, epsilon)
}

type Sphere struct {
	Center Point
	Radius float64
}

func (s *Sphere) Validate() error {
	if s.Radius < 0 {
		return fmt.Errorf("invalid sphere: negative radius")
	}
	return nil
}

func (s *Sphere) NormalVectorAt(p *Point) Vector {
	v := MakeVector(*p, s.Center)
	return v.UnitVector()
}

// Return the point nearest from l[0] intersecting s and l.  Set ok to false if
// there is no intersection.  t is proportional to the distance between the
// intersection point and l[0].
//
// Formulas taken from:
// 	http://www.ccs.neu.edu/home/fell/CSU540/programs/RayTracingFormulas.htm
func SphereLineIntersection(s Sphere, l Line) (p Point, t float64, ok bool) {
	p = Origin
	t = math.MaxFloat64
	ok = false

	dx := l[1].X - l[0].X
	dy := l[1].Y - l[0].Y
	dz := l[1].Z - l[0].Z

	a := dx*dx + dy*dy + dz*dz

	b := 2*dx*(l[0].X-s.Center.X) +
		2*dy*(l[0].Y-s.Center.Y) +
		2*dz*(l[0].Z-s.Center.Z)

	c := s.Center.X*s.Center.X + s.Center.Y*s.Center.Y + s.Center.Z*s.Center.Z +
		l[0].X*l[0].X + l[0].Y*l[0].Y + l[0].Z*l[0].Z -
		2*(s.Center.X*l[0].X+s.Center.Y*l[0].Y+s.Center.Z*l[0].Z) -
		s.Radius*s.Radius

	d := b*b - 4*a*c

	if d < 0 {
		return
	}
	ok = true

	t = (-b - math.Sqrt(d)) / (2 * a)
	p = Point{l[0].X + t*dx, l[0].Y + t*dy, l[0].Z + t*dz}

	return
}

func FloatsEqual(lhs, rhs, epsilon float64) bool {
	return math.Abs(lhs-rhs) < epsilon
}
