package main

import (
	"fmt"
	//"math"

	k "github.com/jsutton9/iron-orbit/kinematics"
	s "github.com/jsutton9/iron-orbit/ships"
	t "github.com/jsutton9/iron-orbit/terrain"
)

func main() {
	star := t.Star{k.Vector{0, 0}, 39438.4}
	ship := s.Ship{k.Vector{1000, 0}, k.Vector{0, 6.28}}
	space := k.Space{1.0, 0.01, 0, []k.GravitySource{}, []k.BodyState{}, []k.OrbitState{}}
	space.AddGravitySource(star)
	space.AddBody(&ship)
	for i := 0; i < 100000; i++ {
		space.StepMotion()
		if (i % 1000 == 0) {
			fmt.Printf("\n(%.2f, %.2f) (%.2f, %.2f)\n", ship.Position().X, ship.Position().Y, ship.Velocity().X, ship.Velocity().Y)
			r := ship.Position().Magnitude()
			speed := ship.Velocity().Magnitude()
			fmt.Printf("%.2f %.2f\n", r, speed)
			kinetic := 0.5*speed*speed
			potential := - space.GravityConstant*star.Mass() / r
			fmt.Printf("%.2f + %.2f = %.2f\n", kinetic, potential, kinetic + potential)
		}
	}
}
