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

type Broker[T any] struct {
	Sub chan chan T
	Unsub chan chan T

	quit chan struct{}
	publish chan T
}

func NewBroker[T any]() *Broker[T] {
	return &Broker[T]{
		Sub: make(chan chan T),
		Unsub: make(chan chan T),

		quit: make(chan struct{}),
		publish: make(chan T),
	}
}

func (b *Broker[T]) start() {
	subs := map[chan T]struct{}{}
	for {
		select {
		case <-b.quit:
			for range b.publish {}
			return
		case ch := <-b.Sub:
			subs[ch] = struct{}{}
		case ch := <-b.Unsub:
			delete(subs, ch)
		case update := <-b.publish:
			for ch, _ := range subs {
				ch <- update
			}
		}
	}
}

type Match struct {
	Space k.Space
	Ships []ShipAndPilot
	Mode TimeMode

	QuitChannel chan struct{}
	TimeChannel chan TimeMode
	MovementBroker *Broker[MovementUpdate]
}

func (match Match) Step() {
	match.Space.StepMotion()
	for i, _ := range match.Ships {
		match.Ships[i].Pilot.Update()
		match.Ships[i].Ship.Thrust(match.Space.TimeStep)
		s := &(match.Ships[i].Ship)
		match.MovementBroker.publish <- MovementUpdate{s.Id, s.P, s.V}
	}
}

func (match Match) Run(quitChannel chan struct{}, timeChannel chan TimeMode) {
	timer := time.NewTimer(0)
	duration := time.Duration(match.Space.TimeStep*float64(time.Second))
	go match.MovementBroker.start()
	for _, s := range match.Ships {
		match.Space.AddBody(&(s.Ship))
	}
	for {
		if match.Mode == Pause {
			select {
			case <- quitChannel:
				match.MovementBroker.quit <- struct{}{}
				return
			case match.Mode = <-timeChannel:
			}
		} else if match.Mode == OneStep {
			match.Step()
			match.Mode = Pause
		} else if match.Mode == RealTime {
			select {
			case <- quitChannel:
				match.MovementBroker.quit <- struct{}{}
				return
			case match.Mode = <-timeChannel:
			case <- timer.C:
				timer = time.NewTimer(duration)
				match.Step()
			}
		} else {
			select {
			case <- quitChannel:
				match.MovementBroker.quit <- struct{}{}
				return
			case match.Mode = <-timeChannel:
			default:
				match.Step()
			}
		}
	}
}
