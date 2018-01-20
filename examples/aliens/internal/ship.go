package internal

import (
	"github.com/explodes/go-wo"
	"github.com/explodes/go-wo/wobj"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	shipSpeed = 120
)

type Ship struct {
	*wobj.Object
}

func newShip(loader wo.Loader, bounds pixel.Rect) (*Ship, error) {
	o, err := wobj.LoadSpriteObject(loader, "img/ship_512.png", bounds.W()/2-12, 0, 24, 24)
	if err != nil {
		return nil, err
	}
	o.Rot = wo.DegToRad(45)
	ship := &Ship{
		Object: o,
	}
	return ship, nil
}

func (o *Ship) Update(dt float64, input wo.Input) {
	if input.Pressed(pixelgl.KeyLeft) {
		o.Move(pixel.V(-shipSpeed*dt, 0))
	}
	if input.Pressed(pixelgl.KeyRight) {
		o.Move(pixel.V(shipSpeed*dt, 0))
	}
}
