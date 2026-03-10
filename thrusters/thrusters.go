package thrusters

import (
	"math"

	m "github.com/jsutton9/iron-orbit/materials"
	v "github.com/jsutton9/iron-orbit/vector"
)

type Thruster struct {
	Mass float64
	MaxThrust float64
	FuelType m.MaterialType
	FuelCost float64

	ThrustFactor float64
	ThrustAngle float64
}

func (thruster Thruster) Impulse(deltaT float64, fuel float64) (v.Vector, float64) {
	cost := deltaT * thruster.ThrustFactor * thruster.FuelCost
	factor := thruster.ThrustFactor
	if (cost > fuel) {
		cost = fuel
		factor = cost / (thruster.FuelCost * deltaT)
	}

	magnitude := factor * thruster.MaxThrust * deltaT
	iX := magnitude*math.Cos(thruster.ThrustAngle)
	iY := magnitude*math.Sin(thruster.ThrustAngle)
	return v.Vector{iX, iY}, cost
}
