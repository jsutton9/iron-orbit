package terrain

import (
	k "github.com/jsutton9/iron-orbit/kinematics"
)

type Star struct {
	P k.Vector
	M float64
}

func (star Star) Position() k.Vector {
	return star.P
}
func (star Star) Mass() float64 {
	return star.M
}

type Planet struct {
	P k.Vector
	M float64
}

func (planet Planet) Position() k.Vector {
	return planet.P
}
func (planet *Planet) SetPosition(position k.Vector) {
	planet.P = position
}
func (planet Planet) Mass() float64 {
	return planet.M
}
