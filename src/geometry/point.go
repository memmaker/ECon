package geometry

import (
	"fmt"
	"math"
)

type Point struct {
	X int
	Y int
}

// String returns a string representation of the form "(x,y)".
func (p Point) String() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}

// Shift returns a new point with coordinates shifted by (x,y). It's a
// shorthand for p.Add(Point{x,y}).
func (p Point) Shift(x, y int) Point {
	return Point{X: p.X + x, Y: p.Y + y}
}

// Add returns vector p+q.
func (p Point) Add(q Point) Point {
	return Point{X: p.X + q.X, Y: p.Y + q.Y}
}

// Sub returns vector p-q.
func (p Point) Sub(q Point) Point {
	return Point{X: p.X - q.X, Y: p.Y - q.Y}
}

// In reports whether the position is within the given range.
func (p Point) In(rg Rect) bool {
	return p.X >= rg.Min.X && p.X < rg.Max.X && p.Y >= rg.Min.Y && p.Y < rg.Max.Y
}

// Mul returns the vector p*k.
func (p Point) Mul(k int) Point {
	return Point{X: p.X * k, Y: p.Y * k}
}

// Div returns the vector p/k.
func (p Point) Div(k int) Point {
	return Point{X: p.X / k, Y: p.Y / k}
}

func Distance(p, q Point) float64 {
	// euclidean distance
	return math.Sqrt(float64((p.X-q.X)*(p.X-q.X) + (p.Y-q.Y)*(p.Y-q.Y)))
}

func ManhattanDistance(p, q Point) int {
	// manhattan distance
	return int(math.Abs(float64(p.X-q.X)) + math.Abs(float64(p.Y-q.Y)))
}
func ChebyshevDistance(p, q Point) int {
	// chebyshev distance
	return int(math.Max(math.Abs(float64(p.X-q.X)), math.Abs(float64(p.Y-q.Y))))
}
