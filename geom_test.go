// Copyright 2012 Daniel Connelly.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"math"
	"testing"
)

const EPS = 0.000000001

func TestDist(t *testing.T) {
	p := Point{2, 3}
	q := Point{5, 6}
	dist := math.Sqrt(18)
	if d := p.dist(q); d != dist {
		t.Errorf("dist(%v, %v) = %v; expected %v", p, q, d, dist)
	}
}

func TestNewBBox(t *testing.T) {
	p := Point{-2.5, 3.0}
	q := Point{5.5, 4.5}
	lengths := []float64{8.0, 1.5}

	bbox, err := NewBBox(p, lengths[0], lengths[1])
	if err != nil {
		t.Errorf("Error on NewBBox(%v, %v): %v", p, lengths, err)
	}
	if d := p.dist(bbox.min); d > EPS {
		t.Errorf("Expected p == bbox.min")
	}
	if d := q.dist(bbox.max); d > EPS {
		t.Errorf("Expected q == bbox.max")
	}
}

func TestNewBBoxDistError(t *testing.T) {
	p := Point{-2.5, 3.0}
	lengths := []float64{-8.0, 1.5}
	_, err := NewBBox(p, lengths[0], lengths[1])
	if _, ok := err.(DistError); !ok {
		t.Errorf("Expected distError on NewBBox(%v, %v)", p, lengths)
	}
}

func TestRectSize(t *testing.T) {
	p := Point{-2.5, 3.0}
	lengths := []float64{8.0, 1.5}
	rect, _ := NewBBox(p, lengths[0], lengths[1])
	size := lengths[0] * lengths[1]
	actual := rect.size()
	if size != actual {
		t.Errorf("Expected %v.size() == %v, got %v", rect, size, actual)
	}
}

func TestRectMargin(t *testing.T) {
	p := Point{-2.5, 3.0}
	lengths := []float64{8.0, 1.5}
	rect, _ := NewBBox(p, lengths[0], lengths[1])
	size := 2*8.0 + 2*1.5
	actual := rect.margin()
	if size != actual {
		t.Errorf("Expected %v.margin() == %v, got %v", rect, size, actual)
	}
}

func TestContainsPoint(t *testing.T) {
	p := Point{-2.4, 0.0}
	lengths := []float64{1.1, 4.9}
	rect, _ := NewBBox(p, lengths[0], lengths[1])

	q := Point{-1.7, 4.8}
	if yes := rect.containsPoint(q); !yes {
		t.Errorf("Expected %v contains %v", rect, q)
	}
}

func TestDoesNotContainPoint(t *testing.T) {
	p := Point{-2.4, 0.0}
	lengths := []float64{1.1, 4.9}
	rect, _ := NewBBox(p, lengths[0], lengths[1])

	q := Point{-1.7, -3.2}
	if yes := rect.containsPoint(q); yes {
		t.Errorf("Expected %v doesn't contain %v", rect, q)
	}
}

func TestContainsBBox(t *testing.T) {
	p := Point{-2.4, 0.0}
	lengths1 := []float64{1.1, 4.9}
	rect1, _ := NewBBox(p, lengths1[0], lengths1[1])

	q := Point{-1.9, 1.0}
	lengths2 := []float64{0.6, 3.7}
	rect2, _ := NewBBox(q, lengths2[0], lengths2[1])

	if yes := rect1.containsBBox(rect2); !yes {
		t.Errorf("Expected %v.containsBBox(%v", rect1, rect2)
	}
}

func TestDoesNotContainRectOverlaps(t *testing.T) {
	p := Point{-2.4, 0.0}
	lengths1 := []float64{1.1, 4.9}
	rect1, _ := NewBBox(p, lengths1[0], lengths1[1])

	q := Point{-1.9, 1.0}
	lengths2 := []float64{1.4, 3.7}
	rect2, _ := NewBBox(q, lengths2[0], lengths2[1])

	if yes := rect1.containsBBox(rect2); yes {
		t.Errorf("Expected %v doesn't contain %v", rect1, rect2)
	}
}

func TestDoesNotContainRectDisjoint(t *testing.T) {
	p := Point{-2.4, 0.0}
	lengths1 := []float64{1.1, 4.9}
	rect1, _ := NewBBox(p, lengths1[0], lengths1[1])

	q := Point{-19.6, -4.0}
	lengths2 := []float64{5.9, 0.5}
	rect2, _ := NewBBox(q, lengths2[0], lengths2[1])

	if yes := rect1.containsBBox(rect2); yes {
		t.Errorf("Expected %v doesn't contain %v", rect1, rect2)
	}
}

func TestNoIntersection(t *testing.T) {
	p := Point{2, 3}
	lengths1 := []float64{1, 1}
	rect1, _ := NewBBox(p, lengths1[0], lengths1[1])

	q := Point{-2, -3}
	lengths2 := []float64{3, 6.5}
	rect2, _ := NewBBox(q, lengths2[0], lengths2[1])

	// rect1 and rect2 fail to overlap in just one dimension (second)

	if intersect := intersect(rect1, rect2); intersect != nil {
		t.Errorf("Expected intersect(%v, %v) == nil, got %v", rect1, rect2, intersect)
	}
}

func TestNoIntersectionJustTouches(t *testing.T) {
	p := Point{2, 3}
	lengths1 := []float64{1, 1}
	rect1, _ := NewBBox(p, lengths1[0], lengths1[1])

	q := Point{-2, -3}
	lengths2 := []float64{4, 6.5}
	rect2, _ := NewBBox(q, lengths2[0], lengths2[1])

	// rect1 and rect2 fail to overlap in just one dimension (second)

	if intersect := intersect(rect1, rect2); intersect != nil {
		t.Errorf("Expected intersect(%v, %v) == nil, got %v", rect1, rect2, intersect)
	}
}

func TestContainmentIntersection(t *testing.T) {
	p := Point{2, 3}
	lengths1 := []float64{1, 1}
	rect1, _ := NewBBox(p, lengths1[0], lengths1[1])

	q := Point{2.2, 3.3}
	lengths2 := []float64{0.5, 0.5}
	rect2, _ := NewBBox(q, lengths2[0], lengths2[1])

	r := Point{2.2, 3.3}
	s := Point{2.7, 3.8}

	actual := intersect(rect1, rect2)
	d1 := r.dist(actual.min)
	d2 := s.dist(actual.max)
	if d1 > EPS || d2 > EPS {
		t.Errorf("intersect(%v, %v) != %v, %v, got %v", rect1, rect2, r, s, actual)
	}
}

func TestOverlapIntersection(t *testing.T) {
	p := Point{2, 3}
	lengths1 := []float64{2.5, 1}
	rect1, _ := NewBBox(p, lengths1[0], lengths1[1])

	q := Point{4, -3}
	lengths2 := []float64{2, 6.5}
	rect2, _ := NewBBox(q, lengths2[0], lengths2[1])

	r := Point{4, 3}
	s := Point{4.5, 3.5}

	actual := intersect(rect1, rect2)
	d1 := r.dist(actual.min)
	d2 := s.dist(actual.max)
	if d1 > EPS || d2 > EPS {
		t.Errorf("intersect(%v, %v) != %v, %v, got %v", rect1, rect2, r, s, actual)
	}
}

func TestToBBox(t *testing.T) {
	x := Point{-2.4, 0.0}
	tol := 0.05
	rect := x.ToBBox(tol)

	p := Point{-2.45, -0.05}
	q := Point{-2.35, 0.05}
	d1 := p.dist(rect.min)
	d2 := q.dist(rect.max)
	if d1 > EPS || d2 > EPS {
		t.Errorf("Expected %v.ToBBox(%v) == %v, %v, got %v", x, tol, p, q, rect)
	}
}

func TestBoundingBox(t *testing.T) {
	p := Point{-2.4, 0.0}
	lengths1 := []float64{15, 3}
	rect1, _ := NewBBox(p, lengths1[0], lengths1[1])

	q := Point{4.7, 2.5}
	lengths2 := []float64{5, 6}
	rect2, _ := NewBBox(q, lengths2[0], lengths2[1])

	r := Point{-2.4, 0.0}
	s := Point{12.6, 8.5}

	bb := boundingBox(rect1, rect2)
	d1 := r.dist(bb.min)
	d2 := s.dist(bb.max)
	if d1 > EPS || d2 > EPS {
		t.Errorf("boundingBox(%v, %v) != %v, %v, got %v", rect1, rect2, r, s, bb)
	}
}

func TestBoundingBoxContains(t *testing.T) {
	p := Point{-2.4, 0.0}
	lengths1 := []float64{15, 3}
	rect1, _ := NewBBox(p, lengths1[0], lengths1[1])

	q := Point{0.0, 1.5}
	lengths2 := []float64{6.222222, 0.946}
	rect2, _ := NewBBox(q, lengths2[0], lengths2[1])

	bb := boundingBox(rect1, rect2)
	d1 := rect1.min.dist(bb.min)
	d2 := rect1.max.dist(bb.max)
	if d1 > EPS || d2 > EPS {
		t.Errorf("boundingBox(%v, %v) != %v, got %v", rect1, rect2, rect1, bb)
	}
}

func TestBoundingBoxN(t *testing.T) {
	rect1, _ := NewBBox(Point{0, 0}, 1, 1)
	rect2, _ := NewBBox(Point{0, 1}, 1, 1)
	rect3, _ := NewBBox(Point{1, 0}, 1, 1)

	exp, _ := NewBBox(Point{0, 0}, 2, 2)
	bb := boundingBoxN(rect1, rect2, rect3)
	d1 := bb.min.dist(exp.min)
	d2 := bb.max.dist(exp.max)
	if d1 > EPS || d2 > EPS {
		t.Errorf("boundingBoxN(%v, %v, %v) != %v, got %v", rect1, rect2, rect3, exp, bb)
	}
}

func TestMinDistZero(t *testing.T) {
	p := Point{2, 3}
	r := p.ToBBox(1)
	if d := p.minDist(r); d > EPS {
		t.Errorf("Expected %v.minDist(%v) == 0, got %v", p, r, d)
	}
}

func TestMinDistPositive(t *testing.T) {
	p := Point{2, 3}
	r := &BBox{Point{-4, 7}, Point{-2, 9}}
	expected := float64((-2-2)*(-2-2) + (7-3)*(7-3))
	if d := p.minDist(r); math.Abs(d-expected) > EPS {
		t.Errorf("Expected %v.minDist(%v) == %v, got %v", p, r, expected, d)
	}
}

func TestMinMaxdist(t *testing.T) {
	p := Point{-2, -1}
	r := &BBox{Point{0, 0}, Point{2, 3}}

	// furthest points from p on the faces closest to p in each dimension
	q1 := Point{2, 3}
	q2 := Point{0, 3}
	q3 := Point{2, 0}

	// find the closest distance from p to one of these furthest points
	d1 := p.dist(q1)
	d2 := p.dist(q2)
	d3 := p.dist(q3)
	expected := math.Min(d1*d1, math.Min(d2*d2, d3*d3))

	if d := p.minMaxDist(r); math.Abs(d-expected) > EPS {
		t.Errorf("Expected %v.minMaxDist(%v) == %v, got %v", p, r, expected, d)
	}
}
