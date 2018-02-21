package wobj

import "github.com/explodes/go-wo"

type Behavior func(source *Object, dt float64)

type Reaction func(source, with *Object, dt float64)

type Behaviors []Behavior

func MakeBehaviors(behaviors ...Behavior) Behaviors {
	return Behaviors(behaviors)
}

func (b Behaviors) Execute(source *Object, dt float64) {
	for _, behavior := range b {
		behavior(source, dt)
	}
}

func Collision(with *Object, reaction Reaction) Behavior {
	return func(source *Object, dt float64) {
		if wo.Collision(source.Bounds(), with.Bounds()) {
			reaction(source, with, dt)
		}
	}
}

var Movement = Behavior(func(source *Object, dt float64) {
	v := source.Velocity.Scaled(dt)
	source.Pos = source.Pos.Add(v)
})

var FaceDirection = FaceDirectionOffset(0)

func FaceDirectionOffset(offset float64) Behavior {
	return func(source *Object, dt float64) {
		source.Rot = source.Velocity.Angle() + offset
	}
}
