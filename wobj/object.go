package wobj

import (
	"github.com/explodes/go-wo"
	"github.com/faiface/pixel"
)

type Object struct {
	Tag string

	Pos  pixel.Vec
	Size pixel.Vec
	Rot  float64

	Velocity pixel.Vec

	Drawable Drawable

	PreSteps  Behaviors
	Steps     Behaviors
	PostSteps Behaviors
}

func (o *Object) Bounds() pixel.Rect {
	return pixel.R(o.Pos.X, o.Pos.Y, o.Pos.X+o.Size.X, o.Pos.Y+o.Size.Y)
}

func (o *Object) Collides(other *Object) bool {
	return wo.Collision(o.Bounds(), other.Bounds())
}

func (o *Object) Move(v pixel.Vec) {
	o.Pos = o.Pos.Add(v)
}

func (o *Object) Draw(target pixel.Target) {
	if o.Drawable == nil {
		return
	}
	bounds := o.Bounds()
	center := pixel.V(bounds.W()/2, bounds.H()/2)
	mat := wo.Fit(o.Drawable.Bounds(), o.Bounds()).Moved(center).Rotated(center.Add(o.Pos), o.Rot)
	o.Drawable.Draw(target, mat)
}
