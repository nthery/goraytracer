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

package geom

import (
	"testing"
)

const epsilon = 0.01

var testData = [...]struct {
	l         Line
	s         Sphere
	intersect bool
	p         Point
}{
	{Line{Origin, Point{4, 2, 0}}, Sphere{Point{7, 5, 0}, 2},
		true, Point{6.27, 3.13, 0}},
	{Line{Point{1, 0.5, 0}, Point{4, 2, 0}}, Sphere{Point{7, 5, 0}, 2},
		true, Point{6.27, 3.13, 0}},
	{Line{Point{1, 0.5, 0}, Point{4, 2, 0}}, Sphere{Point{7, 5, 1}, 2},
		true, Point{6.62, 3.31, 0}},
	{Line{Point{1, 0, 0.5}, Point{4, 0, 2}}, Sphere{Point{7, 0, 5}, 2},
		true, Point{6.27, 0, 3.13}},
	{Line{Point{1, 0.5, 0}, Point{4, 2, 0}}, Sphere{Point{7, 5, 1}, 1},
		false, Origin},
}

func TestSphereLineIntersection(t *testing.T) {
	for _, td := range testData {
		i, _, ok := SphereLineIntersection(td.s, td.l)
		if td.intersect != ok {
			t.Fatalf("bad ok: exp: %v act: %v", td.intersect, ok)
		}
		if ok && !PointsEqual(i, td.p, epsilon) {
			t.Fatalf("bad intersection: exp: %v act: %v", td.p, i)
		}
	}
}

func TestSphereNormalVector(t *testing.T) {
	s := Sphere{Point{1, 1, 1}, 4}
	p := Point{5, 1, 1}
	n := s.NormalVectorAt(&p)
	exp := Vector{1, 0, 0}
	if !VectorsEqual(n, exp, epsilon) {
		t.Fatalf("exp: %v act: %v", exp, n)
	}
}

func TestDotProduct(t *testing.T) {
	v1 := Vector{2, 3, 4}
	v2 := Vector{3, 4, 5}
	dot := DotProduct(&v1, &v2)
	exp := float64(2*3 + 3*4 + 4*5)
	if dot != exp {
		t.Fatalf("exp: %v act: %v", exp, dot)
	}
}
