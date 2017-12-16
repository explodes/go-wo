package internal

import (
	"github.com/explodes/go-wo"

	"github.com/faiface/pixel/pixelgl"
)

const (
	objSpeed = 120
)

type Obj struct {
	*wo.Object
}

func newObj(loader wo.Loader) (*Obj, error) {
	o, err := wo.NewObject(loader, "img/ship_512.png", 200, 200, 24, 24)
	if err != nil {
		return nil, err
	}
	obj := &Obj{
		Object: o,
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
