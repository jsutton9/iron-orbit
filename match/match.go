package match

import (
	"time"

	k "github.com/jsutton9/iron-orbit/kinematics"
	p "github.com/jsutton9/iron-orbit/pilot"
	s "github.com/jsutton9/iron-orbit/ships"
	v "github.com/jsutton9/iron-orbit/vector"
)

type TimeMode int
const (
	Pause TimeMode = 0
	RealTime TimeMode = 1
	OneStep TimeMode = 2
	Fast TimeMode = 3
)

type ShipAndPilot struct {
	Ship s.Ship
	Pilot p.Pilot
}

type MovementUpdate struct {
	Id int
	P v.Vector
	V v.Vector
}

type Match struct {
	Space k.Space
	Ships []ShipAndPilot
	Mode TimeMode
	MovementChannel chan MovementUpdate
}

func (match Match) Step() {
	match.Space.StepMotion()
	for i, _ := range match.Ships {
		match.Ships[i].Pilot.Update()
		match.Ships[i].Ship.Thrust(match.Space.TimeStep)
		s := match.Ships[i].Ship
		match.MovementChannel <- MovementUpdate{s.Id, s.P, s.V}
	}
}

func (match Match) Run(quitChannel chan int, timeChannel chan TimeMode) {
	timer := time.NewTimer(0)
	duration := time.Duration(match.Space.TimeStep*float64(time.Second))
	for _, s := range match.Ships {
		match.Space.AddBody(&(s.Ship))
		match.MovementChannel <- MovementUpdate{s.Ship.Id, s.Ship.P, s.Ship.V}
	}
	for {
		if match.Mode == Pause {
			select {
			case <- quitChannel:
				return
			case match.Mode = <-timeChannel:
			}
		} else if match.Mode == OneStep {
			match.Step()
			match.Mode = Pause
		} else if match.Mode == RealTime {
			select {
			case <- quitChannel:
				return
			case match.Mode = <-timeChannel:
			case <- timer.C:
				timer = time.NewTimer(duration)
				match.Step()
			}
		} else {
			select {
			case <- quitChannel:
				return
			case match.Mode = <-timeChannel:
			default:
				match.Step()
			}
		}
	}
}
