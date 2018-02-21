package wobj

import (
	"github.com/explodes/go-wo"
	"github.com/faiface/pixel"
)

// Behavior is what happens when an object meets a condition for a given time delta
type Behavior func(source *Object, dt float64)

// Reaction is what happens when two objects meet a condition for a given time delta.
// Source is the Object performing a behavior "with" is the object the source is
// reacting with.
type Reaction func(source, with *Object, dt float64)

// Behaviors is a slice of Behavior that should happen in succession
type Behaviors []Behavior

// MakeBehaviors is a convenience function for turning a sequence
// of Behaviors or Behavior functions into Behaviors
func MakeBehaviors(behaviors ...Behavior) Behaviors {
	return Behaviors(behaviors)
}

// Execute executes all behaviors for an object with a time delta
func (b Behaviors) Execute(source *Object, dt float64) {
	for _, behavior := range b {
		behavior(source, dt)
	}
}

// Collision is a behavior that executes when an Object's
// Bounds intersects with a Bounder's Bounds
func Collision(with Bounder, then Behavior) Behavior {
	return func(source *Object, dt float64) {
		if wo.Collision(source.Bounds(), with.Bounds()) {
			then(source, dt)
		}
	}
}

// ObjectCollision is a behavior that executes when an
// Object's Bounds intersects
// with another Object's Bounds
func ObjectCollision(with *Object, reaction Reaction) Behavior {
	return func(source *Object, dt float64) {
		if wo.Collision(source.Bounds(), with.Bounds()) {
			reaction(source, with, dt)
		}
	}
}

// Movement is a Behavior that will move a source an object
// by its velocity scaled by time delta
var Movement = Behavior(func(source *Object, dt float64) {
	v := source.Velocity.Scaled(dt)
	source.Pos = source.Pos.Add(v)
})

// FaceDirection is a behavior that adjusts an Object's
// Rot (rotation) to face the same angle as its Velocity.
func FaceDirection(source *Object, dt float64) {
	source.Rot = source.Velocity.Angle()
}

// ReflectWithin creates a Behavior that will reflect the object
// with confined bounds. The source object will always be placed
// within the bounder's Bounds.
func ReflectWithin(bounder Bounder) Behavior {
	return func(source *Object, dt float64) {
		objBounds := source.Bounds()
		boundary := bounder.Bounds()
		switch {
		case objBounds.Min.X <= boundary.Min.X:
			source.Velocity = pixel.V(-source.Velocity.X, source.Velocity.Y)
			source.Rot = source.Velocity.Angle()
			source.Pos = pixel.V(boundary.Min.X, source.Pos.Y)
		case objBounds.Max.X >= boundary.Max.X:
			source.Velocity = pixel.V(-source.Velocity.X, source.Velocity.Y)
			source.Rot = source.Velocity.Angle()
			source.Pos = pixel.V(boundary.Max.X-source.Size.X, source.Pos.Y)
		}
		switch {
		case objBounds.Min.Y <= boundary.Min.Y:
			source.Velocity = pixel.V(source.Velocity.X, -source.Velocity.Y)
			source.Rot = source.Velocity.Angle()
			source.Pos = pixel.V(source.Pos.X, boundary.Min.Y)
		case objBounds.Max.Y >= boundary.Max.Y:
			source.Velocity = pixel.V(source.Velocity.X, -source.Velocity.Y)
			source.Rot = source.Velocity.Angle()
			source.Pos = pixel.V(source.Pos.X, boundary.Max.Y-source.Size.Y)
		}
	}
}
