package internal

import (
	"github.com/explodes/go-wo"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	ballWidth  = 18
	ballHeight = 18

	ballHitStrength = 525
	ballSlowdown    = 0.75
)

type ball struct {
	pos pixel.Vec
	vel pixel.Vec

	sprite *pixel.Sprite
}

func newBall(sprite *pixel.Sprite) (*ball, error) {
	b := &ball{
		sprite: sprite,
	}
	return b, nil
}

func (b *ball) update(dt float64) {
	b.pos = b.pos.Add(b.vel.Scaled(dt))
	b.vel = b.vel.Scaled(1 - ballSlowdown*dt)
}

func (b *ball) hit(direction pixel.Vec) {
	b.vel = b.vel.Add(direction.Scaled(ballHitStrength))
}

func (b *ball) draw(canvas *pixelgl.Canvas) {
	mat := pixel.IM.Moved(b.pos)
	b.sprite.Draw(canvas, mat)
}

func (b *ball) hitBox() pixel.Rect {
	return wo.HitBox(b.sprite, b.pos)
}
