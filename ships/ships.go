package ships

import (

	m "github.com/jsutton9/iron-orbit/materials"
	t "github.com/jsutton9/iron-orbit/thrusters"
	v "github.com/jsutton9/iron-orbit/vector"
)

type Ship struct {
	Id int
	M float64
	P v.Vector
	V v.Vector
	Thrusters []t.Thruster
	CargoMaterials []m.Material
}

func (ship Ship) Position() v.Vector {
	return ship.P
}
func (ship *Ship) SetPosition(position v.Vector) {
	ship.P = position
}
func (ship Ship) Velocity() v.Vector {
	return ship.V
}
func (ship *Ship) AddVelocity(deltaV v.Vector) {
	ship.V = ship.V.Plus(deltaV)
}
func (ship *Ship) Thrust(deltaT float64) {
	for _, thruster := range ship.Thrusters {
		for i, _ := range ship.CargoMaterials {
			mat := &ship.CargoMaterials[i]
			if (mat.Type != thruster.FuelType) {
				continue
			}
			impulse, cost := thruster.Impulse(deltaT, mat.Mass)
			mat.Mass -= cost
			ship.M -= cost
			ship.AddVelocity(impulse.Scale(1/ship.M))
			break
		}
	}
}
