package internal

import (
	"image/color"

	"github.com/explodes/go-wo"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

var _ Particle = &bullet{}

type bullet struct {
	pos pixel.Vec
	vel pixel.Vec

	size  float64
	color color.Color
}

func newBullet(pos pixel.Vec, vel pixel.Vec, size float64, color color.Color) *bullet {
	return &bullet{
		pos:   pos,
		vel:   vel,
		size:  size,
		color: color,
	}
}

func (b *bullet) Update(dt float64) {
	b.pos = b.pos.Add(b.vel.Scaled(dt))
}

func (b *bullet) Draw(im *imdraw.IMDraw, canvas *pixelgl.Canvas) {
	im.Color = b.color
	im.Push(b.pos)
	im.Circle(b.size, 0)
}

func (b *bullet) HitBox() pixel.Rect {
	return wo.SizedRect(b.pos, pixel.V(b.size, b.size))
}
