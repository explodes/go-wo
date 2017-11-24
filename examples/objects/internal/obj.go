package internal

import (
	"github.com/explodes/go-wo"
	"github.com/faiface/pixel"

	"github.com/faiface/pixel/pixelgl"
)

const (
	objSpeed = 120
)

type Obj struct {
	*wo.Object
}

func newObj(loader wo.Loader) (*Obj, error) {
	sprite, err := loader.Sprite("img/ship_512.png")
	if err != nil {
		return nil, err
	}
	obj := &Obj{
		Object: &wo.Object{
			Sprite: sprite,
			Pos:    pixel.V(400, 50),
			Rot:    wo.DegToRad(45),
			Size:   pixel.V(24, 24),
		},
	}
	return obj, nil
}

func (o *Obj) update(dt float64, input wo.Input) {
	o.Rot += wo.DegToRad(90) * dt

	if input.Pressed(pixelgl.KeyLeft) {
		o.MoveXY(-objSpeed*dt, 0)
	}
	if input.Pressed(pixelgl.KeyRight) {
		o.MoveXY(objSpeed*dt, 0)
	}
	if input.Pressed(pixelgl.KeyUp) {
		o.MoveXY(0, objSpeed*dt)
	}
	if input.Pressed(pixelgl.KeyDown) {
		o.MoveXY(0, -objSpeed*dt)
	}

}

func (o *Obj) draw(canvas *pixelgl.Canvas) {
	o.Object.Draw(canvas)
}
