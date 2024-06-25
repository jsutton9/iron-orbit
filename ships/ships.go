package ships

import (
	k "github.com/jsutton9/iron-orbit/kinematics"
)

type Ship struct {
	P k.Vector
	V k.Vector
}

func (ship Ship) Position() k.Vector {
	return ship.P
}
func (ship *Ship) SetPosition(position k.Vector) {
	ship.P = position
}
func (ship Ship) Velocity() k.Vector {
	return ship.V
}
func (ship *Ship) AddVelocity(deltaV k.Vector) {
	ship.V = ship.V.Plus(deltaV)
}
