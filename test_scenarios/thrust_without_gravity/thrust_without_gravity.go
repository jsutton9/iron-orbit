package main

import (
	"fmt"

	k "github.com/jsutton9/iron-orbit/kinematics"
	mats "github.com/jsutton9/iron-orbit/materials"
	m "github.com/jsutton9/iron-orbit/match"
	p "github.com/jsutton9/iron-orbit/pilot"
	s "github.com/jsutton9/iron-orbit/ships"
	t "github.com/jsutton9/iron-orbit/thrusters"
	v "github.com/jsutton9/iron-orbit/vector"
)

func main() {
	match := m.Match{
		k.Space{1.0, 0.01, 0, []k.GravitySource{}, []k.BodyState{}},
		[]m.ShipAndPilot{},
		m.Fast,
		make(chan struct{}),
		make(chan m.TimeMode),
		m.NewBroker[m.MovementUpdate](),
	}

	match.Ships = append(match.Ships, m.ShipAndPilot{
		s.Ship{
			/*Id=*/1,
			/*M=*/1000,
			/*P=*/v.Vector{0, 0},
			/*V=*/v.Vector{0, 10},
			[]t.Thruster{t.Thruster{1, 500, mats.HOFuel, 100.0, 0.0, 0.5}},
			[]mats.Material{mats.Material{mats.HOFuel, 500}}},
		p.Pilot{}})
	match.Ships[0].Pilot = p.Pilot{
		&(match.Ships[0].Ship),
		p.MaxThrustOnVector,
		0.0, nil, 0.0}

	match.Space.AddBody(&(match.Ships[0].Ship))

	quitChannel := make(chan struct{}, 2)
	timeChannel := make(chan m.TimeMode, 2)
	updates := make(chan m.MovementUpdate)

	go func() {
		match.MovementBroker.Sub <- updates
		for i := 0; i < 1000; i++ {
			movementUpdate := <- updates
			if (i % 100 == 0) {
				t := float64(i) * match.Space.TimeStep
				p := movementUpdate.P
				v := movementUpdate.V
				fmt.Printf("t=%4.1f P=(%5.2f, %5.2f) V=(%5.2f, %5.2f)\n", t, p.X, p.Y, v.X, v.Y)
			}
		}
		quitChannel <- struct{}{}
		for range updates {}
	}()

	match.Run(quitChannel, timeChannel)
}
