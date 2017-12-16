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

func NewObject(loader Loader, spriteName string, x, y, w, h float64) (*Object, error) {
	var sprite *pixel.Sprite
	var err error
	if loader != nil && spriteName != "" {
		sprite, err = loader.Sprite(spriteName)
		if err != nil {
			return nil, err
		}
	}
	o := &Object{
		Sprite: sprite,
		Pos:    pixel.V(x, y),
		Size:   pixel.V(w, h),
		Rot:    0,
	}
	return o, nil
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
	if o.Sprite == nil {
		return
	}
	mat := Fit(o.Sprite.Picture().Bounds(), o.Bounds()).Rotated(o.Pos, o.Rot)
	o.Sprite.Draw(canvas, mat)
}
