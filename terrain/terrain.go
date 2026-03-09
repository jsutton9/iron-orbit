package terrain

import (
	v "github.com/jsutton9/iron-orbit/vector"
)

type Star struct {
	P v.Vector
	M float64
}

func (star Star) Position() v.Vector {
	return star.P
}
func (star Star) Mass() float64 {
	return star.M
}

type Planet struct {
	P v.Vector
	M float64
}

func (planet Planet) Position() v.Vector {
	return planet.P
}
func (planet *Planet) SetPosition(position v.Vector) {
	planet.P = position
}
func (planet Planet) Mass() float64 {
	return planet.M
}
