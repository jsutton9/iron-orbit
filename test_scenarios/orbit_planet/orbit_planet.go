package main

import (
	"fmt"

	k "github.com/jsutton9/iron-orbit/kinematics"
	s "github.com/jsutton9/iron-orbit/ships"
	t "github.com/jsutton9/iron-orbit/terrain"
)

func main() {
	star := t.Star{k.Vector{0, 0}, 1000000}
	planet := t.Planet{k.Vector{0, 0}, 39438.4}
	ship := s.Ship{k.Vector{101000, 0}, k.Vector{0, 6.28+3.16}}
	space := k.Space{1.0, 0.01, 0, []k.GravitySource{}, []k.BodyState{}, []k.OrbitState{}}
	space.AddGravitySource(&star)
	space.AddGravitySource(&planet)
	space.AddOrbiter(&planet, star, 100000, 0, false)
	space.AddBody(&ship)
	for i := 0; i < 100000; i++ {
		space.StepMotion()
		if (i % 1000 == 0) {
			fmt.Printf("\n(%.2f, %.2f) (%.2f, %.2f)\n", ship.Position().X, ship.Position().Y, ship.Velocity().X, ship.Velocity().Y)
			fmt.Printf("(%.2f, %.2f)\n", planet.Position().X, planet.Position().Y)
			rel := ship.Position().Minus(planet.Position())
			fmt.Printf("(%.2f, %.2f)\n", rel.X, rel.Y)
		}
	}
}
