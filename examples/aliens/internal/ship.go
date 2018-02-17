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
	sprite, err := loader.Sprite("img/ship_512.png")
	if err != nil {
		return nil, err
	}
	ship := &Ship{
		Object: &wobj.Object{
			Pos:      pixel.V(bounds.W()/2-12, 0),
			Size:     pixel.V(24, 24),
			Drawable: wobj.NewSpriteDrawable(sprite),
			Rot:      wo.DegToRad(45),
		},
	}
	return ship, nil
}

func (o *Ship) Update(dt float64, input wo.Input) {
	if input.Pressed(pixelgl.KeyLeft) || input.Pressed(pixelgl.KeyA) {
		o.Move(pixel.V(-shipSpeed*dt, 0))
	}
	if input.Pressed(pixelgl.KeyRight) || input.Pressed(pixelgl.KeyD) {
		o.Move(pixel.V(shipSpeed*dt, 0))
	}
}
