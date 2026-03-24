package pilot

import (
	"math"

	k "github.com/jsutton9/iron-orbit/kinematics"
	s "github.com/jsutton9/iron-orbit/ships"
	v "github.com/jsutton9/iron-orbit/vector"
)

type NavigationIntent int
const (
	Idle NavigationIntent = 1
	MaxThrustOnVector NavigationIntent = 2
	Strafe NavigationIntent = 3
	Rendezvous NavigationIntent = 4
)

type Pilot struct {
	Ship *s.Ship
	Navigation NavigationIntent
	Heading float64
	Target k.DynamicBody
	TargetProximity float64
	//GravitySource k.GravitySource
}

func (pilot *Pilot) Update() {
	if pilot.Navigation == MaxThrustOnVector {
		for i, _ := range pilot.Ship.Thrusters {
			pilot.Ship.Thrusters[i].ThrustFactor = 1.0
			pilot.Ship.Thrusters[i].ThrustAngle = pilot.Heading
		}
	}
	if pilot.Navigation == Rendezvous {
		pRelative := pilot.Target.Position().Minus(pilot.Ship.P)
		vRelative := pilot.Target.Velocity().Minus(pilot.Ship.V)
		distance := pRelative.Magnitude()
		radialUnit := pRelative.Scale(1./distance)
		lateralUnit := v.Vector{-radialUnit.Y, radialUnit.X}
		radialSpeed := vRelative.Dot(radialUnit)
		lateralSpeed := vRelative.Dot(lateralUnit)
		if (lateralSpeed < 0) {
			lateralSpeed = -lateralSpeed
			lateralUnit = lateralUnit.Scale(-1)
		}

		totalThrust := 0.0
		for _, thruster := range pilot.Ship.Thrusters {
			totalThrust += thruster.MaxThrust
		}
		aMax := 0.95 * totalThrust / pilot.Ship.M

		lateralMatchTime := lateralSpeed/aMax

		timeForSpeedMatch := radialSpeed/aMax
		distanceForSpeedMatch := 0.5*radialSpeed*timeForSpeedMatch
		if (radialSpeed < 0 && distanceForSpeedMatch > distance) {
			distance = -distance
			radialUnit = radialUnit.Scale(-1)
			timeForSpeedMatch = -timeForSpeedMatch
		}
		radialMatchTime := timeForSpeedMatch + 2*math.Pow((distance + distanceForSpeedMatch)/aMax, 0.5)

		if (radialMatchTime + lateralMatchTime < 0.1) {
			pilot.Navigation = Idle
			for i, _ := range pilot.Ship.Thrusters {
				pilot.Ship.Thrusters[i].ThrustFactor = 0.0
			}
			return
		}

		if (distanceForSpeedMatch > 0.95*distance && distanceForSpeedMatch < 1.05*distance && lateralMatchTime < radialMatchTime) {
			for i, _ := range pilot.Ship.Thrusters {
				pilot.Ship.Thrusters[i].ThrustFactor = 0.0
			}
			return
		}

		pilot.Heading = radialUnit.Scale(radialMatchTime).Plus(lateralUnit.Scale(2*lateralMatchTime)).Angle()
		for i, _ := range pilot.Ship.Thrusters {
			pilot.Ship.Thrusters[i].ThrustFactor = 1.0
			pilot.Ship.Thrusters[i].ThrustAngle = pilot.Heading
		}
	}
}
