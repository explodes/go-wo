package internal

import (
	"github.com/explodes/go-wo"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	playerWidth  = 52
	playerHeight = 52 * 1190 / 2400
	playerSpeed  = 210
)

type player struct {
	pos pixel.Vec
	rot float64

	sprite *pixel.Sprite
}

func newPlayer(sprite *pixel.Sprite) (*player, error) {
	p := &player{
		sprite: sprite,
	}
	return p, nil
}

func (p *player) update(dt float64, ball *ball, input wo.Input) {
	var dx, dy float64
	if input.Pressed(pixelgl.KeyA) {
		dx -= playerSpeed * dt
	}
	if input.Pressed(pixelgl.KeyD) {
		dx += playerSpeed * dt
	}
	if input.Pressed(pixelgl.KeyS) {
		dy -= playerSpeed * dt
	}
	if input.Pressed(pixelgl.KeyW) {
		dy += playerSpeed * dt
	}
	dv := p.pos.Sub(ball.pos).Unit()
	p.rot = dv.Angle() + wo.DegToRad(90)
	p.pos = p.pos.Add(pixel.V(dx, dy))
}

func (p *player) draw(canvas *pixelgl.Canvas) {
	mat := pixel.IM.Moved(p.pos).Rotated(p.pos, p.rot+wo.DegToRad(90))
	p.sprite.Draw(canvas, mat)
}

func (p *player) hitBox() pixel.Rect {
	return wo.HitBox(p.sprite, p.pos)
}
