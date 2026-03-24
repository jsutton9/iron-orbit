package main

import (
	"fmt"

	k "github.com/jsutton9/iron-orbit/kinematics"
	m "github.com/jsutton9/iron-orbit/materials"
	p "github.com/jsutton9/iron-orbit/pilot"
	s "github.com/jsutton9/iron-orbit/ships"
	t "github.com/jsutton9/iron-orbit/thrusters"
	v "github.com/jsutton9/iron-orbit/vector"
)

func trial(space *k.Space, ship *s.Ship, pilot *p.Pilot, target k.DynamicBody) {
	timeStep := 0.01
	for i := 0; i < 5500; i++ {
		if (i % 1 == 0) {
			t := float64(i) * timeStep
			r := ship.Position().Minus(target.Position()).Magnitude()
			s := ship.Velocity().Minus(target.Velocity()).Magnitude()
			fmt.Printf("t=%4.1f P=(%5.2f, %5.2f) V=(%5.2f, %5.2f) M=%6.2f r=%5.2f s=%5.2f\n", t, ship.Position().X, ship.Position().Y, ship.Velocity().X, ship.Velocity().Y, ship.M, r, s)
		}
		space.StepMotion()
		pilot.Update()
		ship.Thrust(timeStep)
	}
}

func main() {
	space := k.Space{1.0, 0.01, 0, []k.GravitySource{}, []k.BodyState{}}

	thruster := t.Thruster{1, 1000, m.HOFuel, 1.0, 0.0, 0.5}
	fuel := m.Material{m.HOFuel, 500}
	ship := s.Ship{1000, v.Vector{0, 0}, v.Vector{0, 0}, []t.Thruster{thruster}, []m.Material{fuel}}
	space.AddBody(&ship)
	target := s.Ship{1000, v.Vector{0, 0}, v.Vector{0, 0}, []t.Thruster{}, []m.Material{}}
	space.AddBody(&target)

	pilot := p.Pilot{&ship, p.Rendezvous, 0.0, &target, 10}

	ship.P = v.Vector{0, 0}
	ship.V = v.Vector{0, 0}
	target.P = v.Vector{100, 0}
	target.V = v.Vector{0, 10}
	trial(&space, &ship, &pilot, &target)
}
