// Copyright 2012 Daniel Connelly.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"math"
)

// DistError is an improper distance measurement.  It implements the error
// and is generated when a distance-related assertion fails.
type DistError float64

func (err DistError) Error() string {
	return "rtree: improper distance"
}

// Point represents a point in n-dimensional Euclidean space.
type Point struct {
	X, Y float64
}

// Dist computes the Euclidean distance between two points p and q.
func (p Point) dist(q Point) float64 {
	dx := p.X - q.X
	dy := p.Y - q.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// minDist computes the square of the distance from a point to a bounding box.
// If the point is contained in the bounding box then the distance is zero.
//
// Implemented per Definition 2 of "Nearest Neighbor Queries" by
// N. Roussopoulos, S. Kelley and F. Vincent, ACM SIGMOD, pages 71-79, 1995.
func (p Point) minDist(bb *BBox) float64 {
	sum := 0.0
	if p.X < bb.min.X {
		d := p.X - bb.min.X
		sum += d * d
	} else if p.X > bb.max.X {
		d := p.X - bb.max.X
		sum += d * d
	}

	if p.Y < bb.min.Y {
		d := p.Y - bb.min.Y
		sum += d * d
	} else if p.Y > bb.max.Y {
		d := p.Y - bb.max.Y
		sum += d * d
	}

	return sum
}

// minMaxDist computes the minimum of the maximum distances from p to points
// on r.  If r is the bounding box of some geometric objects, then there is
// at least one object contained in r within minMaxDist(p, r) of p.
//
// Implemented per Definition 4 of "Nearest Neighbor Queries" by
// N. Roussopoulos, S. Kelley and F. Vincent, ACM SIGMOD, pages 71-79, 1995.
func (p Point) minMaxDist(bb *BBox) float64 {
	// by definition, MinMaxDist(p, r) =
	// min{1<=k<=n}(|pk - rmk|^2 + sum{1<=i<=n, i != k}(|pi - rMi|^2))
	// where rmk and rMk are defined as follows:

	var rmx, rmy, rMx, rMy float64
	if p.X <= (bb.min.X+bb.max.X)/2 {
		rmx = bb.min.X
	} else {
		rmx = bb.max.X
	}

	if p.Y <= (bb.min.Y+bb.max.Y)/2 {
		rmy = bb.min.Y
	} else {
		rmy = bb.max.Y
	}

	if p.X >= (bb.min.X+bb.max.X)/2 {
		rMx = bb.min.X
	} else {
		rMx = bb.max.X
	}

	if p.Y >= (bb.min.Y+bb.max.Y)/2 {
		rMy = bb.min.Y
	} else {
		rMy = bb.max.Y
	}

	// This formula can be computed in linear time by precomputing
	// S = sum{1<=i<=n}(|pi - rMi|^2).

	s := 0.0
	d := p.X - rMx
	s += d * d
	d = p.Y - rMy
	s += d * d

	// Compute MinMaxDist using the precomputed s for X.
	d1 := p.X - rMx
	d2 := p.X - rmx
	d = s - d1*d1 + d2*d2
	min := d

	// and for Y
	d1 = p.Y - rMy
	d2 = p.Y - rmy
	d = s - d1*d1 + d2*d2
	if d < min {
		min = d
	}

	return min
}

// BBox represents a subset of 2-dimensional Euclidean space of the form
// min:[a1, b1] x max:[a2, b2], where a1 < a2 and b1 < b2
type BBox struct {
	min, max Point
}

func (bb *BBox) String() string {
	return fmt.Sprintf("%sx%s", bb.min, bb.max)
}

func (p *Point) String() string {
	return fmt.Sprintf("[%.2f, %.2f]", p.X, p.Y)
}

// NewRect constructs and returns a pointer to a Bbox given two corner points.
// The point p should be the most-negative point on the rectangle and x, y
// should be positive lengths.
func NewBBox(p Point, x, y float64) (*BBox, error) {
	if x < 0 {
		return nil, DistError(x)
	}
	if y < 0 {
		return nil, DistError(y)
	}

	return &BBox{
		min: p,
		max: Point{X: p.X + x, Y: p.Y + y},
	}, nil
}

// size computes the measure of a bounding box
func (bb *BBox) size() float64 {
	return (bb.max.X - bb.min.X) * (bb.max.Y - bb.min.Y)
}

// margin computes the sum of the edge lengths of a bounding box.
func (bb *BBox) margin() float64 {
	return 2 * ((bb.max.X - bb.min.X) + (bb.max.Y - bb.min.Y))
}

// containsPoint tests whether p is located inside or on the boundary of bb.
func (bb *BBox) containsPoint(p Point) bool {
	return bb.min.X < p.X && bb.max.X > p.X && bb.min.Y < p.Y && bb.max.Y > p.Y
}

// containsBBox tests whether bb2 is is located inside bb.
func (bb *BBox) containsBBox(bb2 *BBox) bool {
	return bb.min.X <= bb2.min.X && bb.max.X >= bb2.max.X && bb.min.Y <= bb2.min.Y && bb.max.Y >= bb2.max.Y
}

// intersect computes the intersection of two bounding boxes.  If no
// intersection exists, the intersection is nil.
func intersect(bb1, bb2 *BBox) *BBox {
	// There are four cases of overlap:
	//
	//     1.  a1------------b1
	//              a2------------b2
	//              p--------q
	//
	//     2.       a1------------b1
	//         a2------------b2
	//              p--------q
	//
	//     3.  a1-----------------b1
	//              a2-------b2
	//              p--------q
	//
	//     4.       a1-------b1
	//         a2-----------------b2
	//              p--------q
	//
	// There are only two cases of non-overlap:
	//
	//     1. a1------b1
	//                    a2------b2
	//
	//     2.             a1------b1
	//        a2------b2

	if bb1.max.X < bb2.min.X || bb2.max.X < bb1.min.X || bb1.max.Y < bb2.min.Y || bb2.max.Y < bb2.min.Y {
		return nil
	}
	return &BBox{
		min: Point{X: math.Max(bb1.min.X, bb2.min.X), Y: math.Max(bb1.min.Y, bb2.min.Y)},
		max: Point{X: math.Min(bb1.max.X, bb2.max.X), Y: math.Min(bb1.max.Y, bb2.max.Y)},
	}
}

// ToBBox constructs a bounding box containing p with side lengths 2*tol.
func (p Point) ToBBox(tol float64) *BBox {
	return &BBox{
		min: Point{X: p.X - tol, Y: p.Y - tol},
		max: Point{X: p.X + tol, Y: p.Y + tol},
	}
}

// boundingBox constructs the smallest bounding box containing both bb1 and bb2.
func boundingBox(bb1, bb2 *BBox) *BBox {
	return &BBox{
		min: Point{X: math.Min(bb1.min.X, bb2.min.X), Y: math.Min(bb1.min.Y, bb2.min.Y)},
		max: Point{X: math.Max(bb1.max.X, bb2.max.X), Y: math.Max(bb1.max.Y, bb2.max.Y)},
	}
}

// boundingBoxN constructs the smallest rectangle containing all of bbs...
func boundingBoxN(bbs ...*BBox) *BBox {
	if len(bbs) == 1 {
		return bbs[0]
	}
	bb := boundingBox(bbs[0], bbs[1])
	for _, other := range bbs[2:] {
		bb = boundingBox(bb, other)
	}
	return bb
}
