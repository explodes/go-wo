package internal

import (
	"github.com/explodes/go-wo"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type Particle interface {
	Draw(im *imdraw.IMDraw, canvas *pixelgl.Canvas)
	Update(dt float64)
	HitBox() pixel.Rect
}

type ParticleSet struct {
	im        *imdraw.IMDraw
	particles map[Particle]struct{}
}

func NewParticleCollection(size int) *ParticleSet {
	return &ParticleSet{
		im:        imdraw.New(nil),
		particles: make(map[Particle]struct{}, size),
	}
}

func (o *ParticleSet) Add(p Particle) {
	o.particles[p] = struct{}{}
}

func (o *ParticleSet) Remove(p Particle) {
	delete(o.particles, p)
}

func (o *ParticleSet) Draw(canvas *pixelgl.Canvas) {
	o.im.Clear()
	for p := range o.particles {
		p.Draw(o.im, canvas)
	}
	o.im.Draw(canvas)
}

func (o *ParticleSet) Update(dt float64) {
	for p := range o.particles {
		p.Update(dt)
	}
}

var _ Particle = &animatedParticle{}

type animatedParticle struct {
	pos  pixel.Vec
	size pixel.Vec
	age  float64

	frameNum int

	fps float64

	sheet *wo.SpriteSheet

	set                *ParticleSet
	removeWhenComplete bool
}

func newAnimatedParticle(particles *ParticleSet, sheet *wo.SpriteSheet, fps float64, pos, size pixel.Vec, removeWhenComplete bool) *animatedParticle {
	return &animatedParticle{
		pos:                pos,
		size:               size,
		fps:                fps,
		sheet:              sheet,
		set:                particles,
		removeWhenComplete: removeWhenComplete,
	}
}

func (p *animatedParticle) Update(dt float64) {
	p.age += dt
	frameNum := int(p.age * p.fps)
	if frameNum >= p.sheet.NumFrames() && p.removeWhenComplete {
		p.set.Remove(p)
		return
	}
	p.frameNum = frameNum % p.sheet.NumFrames()
}

func (p *animatedParticle) Draw(im *imdraw.IMDraw, canvas *pixelgl.Canvas) {
	sprite := p.sheet.SetFrame(p.frameNum)
	sx := p.size.X / float64(p.sheet.Options().Width)
	sy := p.size.Y / float64(p.sheet.Options().Height)
	mat := pixel.IM.Moved(p.pos).ScaledXY(p.pos, pixel.V(sx, sy))
	sprite.Draw(canvas, mat)
}

func (p *animatedParticle) HitBox() pixel.Rect {
	return wo.SizedRect(p.pos, p.size)
}
