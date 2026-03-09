package kinematics

import (
	//"fmt"
	"math"

	v "github.com/jsutton9/iron-orbit/vector"
)

type GravitySource interface {
	Position() v.Vector
	Mass() float64
}

/*
type Orbiter interface {
	SetPosition(position v.Vector)
}

type OrbitState struct {
	Handle Orbiter
	Center GravitySource
	Radius float64
	PhaseOffset float64
	Rate float64
	NextPosition v.Vector
	Clockwise bool
}
*/

type DynamicBody interface {
	Position() v.Vector
	SetPosition(position v.Vector)
	Velocity() v.Vector
	AddVelocity(deltaV v.Vector)
}

type BodyState struct {
	Handle DynamicBody
	NextPosition v.Vector
}

type Space struct {
	GravityConstant float64
	TimeStep float64
	Time float64
	GravitySources []GravitySource
	Bodies []BodyState
	//Orbiters []OrbitState
}

/*
func (state OrbitState) Position(t float64) v.Vector {
	angle := math.Mod((t*state.Rate) + state.PhaseOffset, (2*math.Pi))
	if state.Clockwise {
		angle = -angle
	}
	unitVector := v.Vector{math.Cos(angle), math.Sin(angle)}
	return unitVector.Scale(state.Radius).Plus(state.Center.Position())
}
*/

func (space *Space) AddGravitySource(source GravitySource) {
	space.GravitySources = append(space.GravitySources, source)
}

/*
func (space *Space) AddOrbiter(handle Orbiter, center GravitySource, radius float64, phaseOffset float64, clockwise bool) {
	rate := math.Sqrt(space.GravityConstant*center.Mass()/(radius*radius*radius))
	state := OrbitState{handle, center, radius, phaseOffset, rate, v.Vector{0, 0}, clockwise}
	state.Handle.SetPosition(state.Position(space.Time))
	space.Orbiters = append(space.Orbiters, state)
}
*/

func (space *Space) AddBody(handle DynamicBody) {
	space.Bodies = append(space.Bodies, BodyState{handle, v.Vector{0, 0}})
}

func (space *Space) Gravitate(source GravitySource, body *BodyState) {
	gravityStep := space.TimeStep * space.GravityConstant * source.Mass()
	deltaPosition := source.Position().Minus(body.Handle.Position())
	distanceFactor := math.Pow(deltaPosition.Magnitude(), -3)
	body.Handle.AddVelocity(deltaPosition.Scale(gravityStep*distanceFactor))
}

func (space *Space) StepMotion() {
	for i, _ := range space.Bodies {
		body := &space.Bodies[i]
		body.NextPosition = body.Handle.Position().Plus(body.Handle.Velocity().Scale(space.TimeStep))
	}
	/*
	for i, _ := range space.Orbiters {
		o := &space.Orbiters[i]
		o.NextPosition = o.Position(space.Time + space.TimeStep)
	}
	*/
	for i, _ := range space.Bodies {
		for _, source := range space.GravitySources{
			space.Gravitate(source, &space.Bodies[i])
		}
	}
	for i, _ := range space.Bodies {
		body := &space.Bodies[i]
		body.Handle.SetPosition(body.NextPosition)
	}
	/*
	for i, _ := range space.Orbiters {
		o := &space.Orbiters[i]
		o.Handle.SetPosition(o.NextPosition)
	}
	*/
	space.Time += space.TimeStep
}
