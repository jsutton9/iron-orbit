package vector

import (
	"math"
)

type Vector struct {
	X float64
	Y float64
}
func (v1 Vector) Plus(v2 Vector) Vector {
	return Vector{v1.X + v2.X, v1.Y + v2.Y}
}
func (v1 Vector) Minus(v2 Vector) Vector {
	return Vector{v1.X - v2.X, v1.Y - v2.Y}
}
func (v Vector) Scale(scalar float64) Vector {
	return Vector{v.X*scalar, v.Y*scalar}
}
func (v Vector) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}
