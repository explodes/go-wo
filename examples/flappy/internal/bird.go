package internal

import (
	"github.com/explodes/go-wo"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	// flapRate is how many times per second the frame
	// advances in the flying animation of the bird
	flapRate = 6

	horiSpeed = 2750
	vertSpeed = 200
	gravity   = -175
	flapWait  = 0.85

	deadWeight = 3
)

// playBird is a player-controlled bird
type playBird struct {
	physics

	// age is how old the bird is in seconds
	age float64

	// flapWait is how long left until the
	// bird can flag again
	flapWait float64
}

func (b *playBird) update(dt float64, input wo.Input, dead bool) {

	if dead {
		b.speedXY(0, b.vel.Y+gravity*dt*deadWeight)
	} else {
		b.updateAlive(dt, input)
	}

	b.physics.update(dt)
	b.age += dt
	b.flapWait -= dt
}

func (b *playBird) updateAlive(dt float64, input wo.Input) {
	dx := b.vel.X
	if input.Pressed(pixelgl.KeyLeft) {
		dx = -horiSpeed * dt
	}
	if input.Pressed(pixelgl.KeyRight) {
		dx = horiSpeed * dt
	}
	dy := b.vel.Y
	if b.flapWait <= 0 && input.Pressed(pixelgl.KeyUp) {
		dy += vertSpeed
		b.flapWait = flapWait
	}
	dy += gravity * dt
	b.speedXY(dx, dy)
}

func (b *playBird) draw(canvas *pixelgl.Canvas, sprites []*pixel.Sprite) {
	frameNum := int(b.age*flapRate) % len(sprites)
	frame := sprites[frameNum]
	mat := pixel.IM.Moved(b.pos)
	frame.Draw(canvas, mat)
}
