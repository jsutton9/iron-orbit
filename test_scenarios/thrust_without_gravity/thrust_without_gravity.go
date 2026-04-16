package main

import (
	k "github.com/jsutton9/iron-orbit/kinematics"
	mats "github.com/jsutton9/iron-orbit/materials"
	m "github.com/jsutton9/iron-orbit/match"
	p "github.com/jsutton9/iron-orbit/pilot"
	s "github.com/jsutton9/iron-orbit/ships"
	sv "github.com/jsutton9/iron-orbit/server"
	t "github.com/jsutton9/iron-orbit/thrusters"
	v "github.com/jsutton9/iron-orbit/vector"
)

func main() {
	match := m.Match{
		k.Space{1.0, 0.01, 0, []k.GravitySource{}, []k.BodyState{}},
		[]m.ShipAndPilot{},
		m.Pause,
		make(chan struct{}),
		make(chan m.TimeMode),
		m.NewBroker[m.TrackingUpdate](),
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

	go match.Run()
	sv.Serve(&match)
}
