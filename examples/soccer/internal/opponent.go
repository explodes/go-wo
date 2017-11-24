package internal

import (
	"github.com/explodes/go-wo"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	opponentWidth  = 52
	opponentHeight = 52 * 397 / 800
	opponentSpeed  = 140
)

type opponent struct {
	pos pixel.Vec
	rot float64

	sprite *pixel.Sprite
}

func newOpponent(sprite *pixel.Sprite) (*opponent, error) {
	p := &opponent{
		sprite: sprite,
	}
	return p, nil
}

func (o *opponent) update(dt float64, ball *ball) {
	dv := o.pos.Sub(ball.pos).Unit()
	o.rot = dv.Angle() + wo.DegToRad(90)
	o.pos = o.pos.Sub(dv.Scaled(opponentSpeed * dt))
}

func (o *opponent) draw(canvas *pixelgl.Canvas) {
	mat := pixel.IM.Moved(o.pos).Rotated(o.pos, o.rot+wo.DegToRad(90))
	o.sprite.Draw(canvas, mat)
}

func (o *opponent) hitBox() pixel.Rect {
	return wo.HitBox(o.sprite, o.pos)
}
