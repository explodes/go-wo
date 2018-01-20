package internal

import (
	"github.com/explodes/go-wo"
	"github.com/explodes/go-wo/wobj"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	objSpeed = 120
)

type Obj struct {
	*wobj.Object
}

func newObj(loader wo.Loader) (*Obj, error) {
	o, err := wobj.LoadSpriteObject(loader, "img/ship_512.png", 200, 200, 24, 24)
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
		o.Move(pixel.V(-objSpeed*dt, 0))
	}
	if input.Pressed(pixelgl.KeyRight) {
		o.Move(pixel.V(objSpeed*dt, 0))
	}
	if input.Pressed(pixelgl.KeyUp) {
		o.Move(pixel.V(0, objSpeed*dt))
	}
	if input.Pressed(pixelgl.KeyDown) {
		o.Move(pixel.V(0, -objSpeed*dt))
	}

}
