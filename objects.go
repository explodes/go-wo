package wo

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Object struct {
	Sprite *pixel.Sprite
	Pos    pixel.Vec
	Size   pixel.Vec
	Rot    float64
}

func (o *Object) Bounds() pixel.Rect {
	return pixel.R(o.Pos.X, o.Pos.Y, o.Pos.X+o.Size.X, o.Pos.Y+o.Size.Y)
}

func (o *Object) MoveXY(x, y float64) {
	o.Move(pixel.V(x, y))
}

func (o *Object) Move(v pixel.Vec) {
	o.Pos = o.Pos.Add(v)
}

func (o *Object) Draw(canvas *pixelgl.Canvas) {
	mat := Fit(o.Sprite.Picture().Bounds(), o.Bounds()).Rotated(o.Pos, o.Rot)
	o.Sprite.Draw(canvas, mat)
}
