package main

import (
	"fmt"

	k "github.com/jsutton9/iron-orbit/kinematics"
	m "github.com/jsutton9/iron-orbit/materials"
	s "github.com/jsutton9/iron-orbit/ships"
	t "github.com/jsutton9/iron-orbit/thrusters"
	v "github.com/jsutton9/iron-orbit/vector"
)

func main() {
	space := k.Space{1.0, 0.01, 0, []k.GravitySource{}, []k.BodyState{}}

	thruster := t.Thruster{1, 500, m.HOFuel, 0.01, 1.0, 0.0}
	fuel := m.Material{m.HOFuel, 5}
	ship := s.Ship{1000, v.Vector{0, 0}, v.Vector{0, 10}, []t.Thruster{thruster}, []m.Material{fuel}}
	space.AddBody(&ship)

	timeStep := 0.01
	for i := 0; i < 1000; i++ {
		t := float64(i) * timeStep
		if (i % 100 == 0) {
			fmt.Printf("t=%4.1f P=(%5.2f, %5.2f) V=(%5.2f, %5.2f)\n", t, ship.Position().X, ship.Position().Y, ship.Velocity().X, ship.Velocity().Y)
		}
		space.StepMotion()
		ship.Thrust(timeStep)
	}
}
