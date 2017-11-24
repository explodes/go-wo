package internal

import (
	"fmt"

	"github.com/faiface/pixel"
)

type physics struct {
	pos pixel.Vec
	vel pixel.Vec
	acc pixel.Vec
}

func (p *physics) forceXY(ddx, ddy float64) {
	p.acc.X += ddx
	p.acc.Y += ddy
}

func (p *physics) speedXY(dx, dy float64) {
	p.vel.X = dx
	p.vel.Y = dy
}

func (p *physics) update(dt float64) {
	p.pos.X += p.vel.X * dt
	p.pos.Y += p.vel.Y * dt
	p.vel.X += p.acc.X*dt - p.vel.X*dt
	p.vel.Y += p.acc.Y*dt - p.vel.Y*dt
}

func (p *physics) String() string {
	return fmt.Sprintf("xy(%.02f, %.02f) dxy(%.02f, %.02f) ddxy(%.02f, %.02f)", p.pos.X, p.pos.Y, p.vel.X, p.vel.Y, p.acc.X, p.acc.Y)
}
