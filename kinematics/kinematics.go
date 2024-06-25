package kinematics

import (
	//"fmt"
	"math"
)

type Vector struct {
	X float64
	Y float64
}
func (v1 Vector) Plus(v2 Vector) Vector {
	return Vector{v1.X + v2.X, v1.Y + v2.Y}
}
func (v1 Vector) Minus(v2 Vector) Vector {
	return Vector{v1.X - v2.X, v1.Y - v2.Y}
}
func (v Vector) Scale(scalar float64) Vector {
	return Vector{v.X*scalar, v.Y*scalar}
}
func (v Vector) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

type GravitySource interface {
	Position() Vector
	Mass() float64
}

type Orbiter interface {
	SetPosition(position Vector)
}

type OrbitState struct {
	Handle Orbiter
	Center GravitySource
	Radius float64
	PhaseOffset float64
	Rate float64
	NextPosition Vector
	Clockwise bool
}

type DynamicBody interface {
	Position() Vector
	SetPosition(position Vector)
	Velocity() Vector
	AddVelocity(deltaV Vector)
}

type BodyState struct {
	Handle DynamicBody
	NextPosition Vector
}

type Space struct {
	GravityConstant float64
	TimeStep float64
	Time float64
	GravitySources []GravitySource
	Bodies []BodyState
	Orbiters []OrbitState
}

func (state OrbitState) Position(t float64) Vector {
	angle := math.Mod((t*state.Rate) + state.PhaseOffset, (2*math.Pi))
	if state.Clockwise {
		angle = -angle
	}
	unitVector := Vector{math.Cos(angle), math.Sin(angle)}
	return unitVector.Scale(state.Radius).Plus(state.Center.Position())
}

func (space *Space) AddGravitySource(source GravitySource) {
	space.GravitySources = append(space.GravitySources, source)
}

func (space *Space) AddOrbiter(handle Orbiter, center GravitySource, radius float64, phaseOffset float64, clockwise bool) {
	rate := math.Sqrt(space.GravityConstant*center.Mass()/(radius*radius*radius))
	state := OrbitState{handle, center, radius, phaseOffset, rate, Vector{0, 0}, clockwise}
	state.Handle.SetPosition(state.Position(space.Time))
	space.Orbiters = append(space.Orbiters, state)
}

func (space *Space) AddBody(handle DynamicBody) {
	space.Bodies = append(space.Bodies, BodyState{handle, Vector{0, 0}})
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
	for i, _ := range space.Orbiters {
		o := &space.Orbiters[i]
		o.NextPosition = o.Position(space.Time + space.TimeStep)
	}
	for i, _ := range space.Bodies {
		for _, source := range space.GravitySources{
			space.Gravitate(source, &space.Bodies[i])
		}
	}
	for i, _ := range space.Bodies {
		body := &space.Bodies[i]
		body.Handle.SetPosition(body.NextPosition)
	}
	for i, _ := range space.Orbiters {
		o := &space.Orbiters[i]
		o.Handle.SetPosition(o.NextPosition)
	}
	space.Time += space.TimeStep
}
