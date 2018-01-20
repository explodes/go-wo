package wobj

import (
	"github.com/explodes/go-wo"
	"github.com/faiface/pixel"
)

type Object struct {
	Drawable Drawable
	Pos      pixel.Vec
	Size     pixel.Vec
	Rot      float64
}

func LoadSpriteObject(loader wo.Loader, spriteName string, x, y, w, h float64, transforms ...wo.ImageTransformer) (*Object, error) {
	var sprite *pixel.Sprite
	var err error
	if loader != nil && spriteName != "" {
		sprite, err = loader.Sprite(spriteName, transforms...)
		if err != nil {
			return nil, err
		}
	}
	return NewSpriteObject(sprite, x, y, w, h), nil
}

func NewSpriteObject(sprite *pixel.Sprite, x, y, w, h float64) *Object {
	drawable := &SpriteDrawable{
		Sprite: sprite,
	}
	return NewDrawableObject(drawable, x, y, w, h)
}

func NewDrawableObject(drawable Drawable, x, y, w, h float64) *Object {
	return &Object{
		Drawable: drawable,
		Pos:      pixel.V(x, y),
		Size:     pixel.V(w, h),
		Rot:      0,
	}
}

func (o *Object) Bounds() pixel.Rect {
	return pixel.R(o.Pos.X, o.Pos.Y, o.Pos.X+o.Size.X, o.Pos.Y+o.Size.Y)
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
