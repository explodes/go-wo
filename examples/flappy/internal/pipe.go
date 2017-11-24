package internal

import (
	"math/rand"

	"github.com/explodes/go-wo"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/sirupsen/logrus"
)

const (
	pipeWidth        = 52
	basePipeSpeed    = 120
	pipeMin, pipeMax = 0.4, 0.6
	pipeGapSize      = 100
)

type playPipe struct {
	scored bool
	pos    pixel.Vec
}

func (p *playPipe) update(dt float64, speed float64, r *rand.Rand, limit pixel.Rect) {
	p.pos.X -= speed * dt

	if p.pos.X < -pipeWidth {
		p.reset(r, limit)
	}
}

func (p *playPipe) reset(r *rand.Rand, limit pixel.Rect) {
	p.scored = false
	p.pos.X = limit.Max.X
	p.pos.Y = limit.Min.Y + limit.H()*(pipeMin+(pipeMax-pipeMin)*r.Float64())
	logrus.WithField("pipeY", p.pos.Y).Info("pipe reset")
}

func (p *playPipe) draw(canvas *pixelgl.Canvas, sprite *pixel.Sprite) {
	mat := pixel.IM.Moved(p.pos)
	sprite.Draw(canvas, mat)

	bounds := wo.HitBox(sprite, p.pos)
	offset := pixel.V(0, bounds.H()+pipeGapSize)

	mat = mat.Moved(offset).ScaledXY(bounds.Center().Add(offset), pixel.V(1, -1))
	sprite.Draw(canvas, mat)
}

func (p *playPipe) invertedHitBox(sprite *pixel.Sprite) pixel.Rect {
	bounds := wo.HitBox(sprite, p.pos)
	offset := pixel.V(0, bounds.H()+pipeGapSize)
	return pixel.R(
		bounds.Min.X+offset.X,
		bounds.Min.Y+offset.Y,
		bounds.Max.X+offset.X,
		bounds.Max.Y+offset.Y,
	)
}
