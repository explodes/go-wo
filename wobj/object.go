package wobj

import (
	"github.com/explodes/go-wo"
	"github.com/faiface/pixel"
)

// Object is a game object that has basic physics, optional
// graphics, and associated Behaviors. It can be used standalone
// or managed (Updated and Drawn) by Objects.
type Object struct {
	// Tag is an optional identifier for this type of object.
	// It can be retrieved as an ObjectSet from an Objects by
	// this tag along with other Objects with the same tag.
	Tag string

	// Pos is the position of the Object. The Drawable, if any,
	// will be drawn with this as the origin.
	Pos pixel.Vec
	// Size is the size of the Object. The Drawable, if any,
	// will be scaled to fit.
	Size pixel.Vec
	// Velocity is the Vec describing the movement speed
	// and direction of this Object.
	Velocity pixel.Vec

	// Drawable is an optional Drawable to use to draw this
	// Object on a Target.
	Drawable Drawable
	// Rot is an amount in radians used to rotate the Drawable.
	Rot float64

	// PreSteps is Behaviors to execute before Steps and
	// PostSteps during an Update performed by Objects.
	PreSteps Behaviors
	// Steps is Behaviors to execute before PostSteps and
	// after PreSteps during an Update performed by Objects.
	Steps Behaviors
	// PostSteps is Behaviors to execute after Steps during
	// an Update performed by Objects.
	PostSteps Behaviors
}

// Bounds gets the hitbox for this Object. Any Drawable will
// scaled and translated to fit this box. Collision detection
// can be performed using this Rect.
func (o *Object) Bounds() pixel.Rect {
	return pixel.R(o.Pos.X, o.Pos.Y, o.Pos.X+o.Size.X, o.Pos.Y+o.Size.Y)
}

// Collides tests to see if this Object collides with another.
func (o *Object) Collides(other *Object) bool {
	return wo.Collision(o.Bounds(), other.Bounds())
}

// Move is a helper function used to add a Vec to this Object's current Pos.
func (o *Object) Move(v pixel.Vec) {
	o.Pos = o.Pos.Add(v)
}

// Draw will render this Object on a target if a Drawable is associated with
// this Object. The Object's Drawable will be scaled and translated to fit
// this Object's Bounds. It will also be rotated by Rot radians to
// This function does nothing if this Object has no Drawable.
func (o *Object) Draw(target pixel.Target) {
	if o.Drawable == nil {
		return
	}
	bounds := o.Bounds()
	center := pixel.V(bounds.W()/2, bounds.H()/2)
	mat := wo.Fit(o.Drawable.Bounds(), o.Bounds()).Moved(center).Rotated(center.Add(o.Pos), o.Rot)
	o.Drawable.Draw(target, mat)
}
