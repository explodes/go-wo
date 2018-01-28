package internal

import (
	"github.com/explodes/go-wo"
	"github.com/explodes/go-wo/wobj"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Direction uint8

const (
	DirectionLeft Direction = iota
	DirectionRight
)

const (
	alienSize           = 28
	baseAlienSpeed      = 30
	alienSpeedPerSecond = 0.5
)

var (
	alienSheetOptions = wo.SpriteSheetOptions{
		Width:      16,
		Height:     16,
		Columns:    100,
		Rows:       75,
		ExactCount: 7500,
	}
)

type AlienRow struct {
	aliens    []*Alien
	direction Direction
	speed     float64
}

func (a *AlienRow) Update(dt float64, bounds pixel.Rect) {
	a.speed += dt * alienSpeedPerSecond
	if a.direction == DirectionLeft {
		a.UpdateDirection(dt, -1)
		for _, alien := range a.aliens {
			if alien == nil {
				continue
			}
			if alien.Bounds().Min.X <= bounds.Min.X {
				a.direction = DirectionRight
				a.DownShift()
			}
			break
		}
	} else if a.direction == DirectionRight {
		a.UpdateDirection(dt, 1)
		for i := len(a.aliens) - 1; i >= 0; i++ {
			alien := a.aliens[i]
			if alien == nil {
				continue
			}
			if alien.Bounds().Max.X >= bounds.Max.X {
				a.direction = DirectionLeft
				a.DownShift()
			}
			break
		}
	}
}

func (a *AlienRow) DownShift() {
	for _, alien := range a.aliens {
		if alien == nil {
			continue
		}
		alien.Move(pixel.V(0, -alien.Bounds().H()))
	}
}

func (a *AlienRow) Draw(canvas *pixelgl.Canvas) {
	for _, alien := range a.aliens {
		if alien == nil {
			continue
		}
		alien.Draw(canvas)
	}
}

func (a *AlienRow) UpdateDirection(dt, moveFactor float64) {
	for _, alien := range a.aliens {
		if alien == nil {
			continue
		}
		alien.Move(pixel.V(a.speed*dt*moveFactor, 0))
	}
}

func randomAlienKind() int {
	return rng.Intn(alienSheetOptions.ExactCount)
}

type Alien struct {
	*wobj.Object
	drawable *wobj.SpriteSheetDrawable
	kind     int
}

func newAlien(drawable *wobj.SpriteSheetDrawable, pos pixel.Vec, kind int) *Alien {
	o := wobj.NewDrawableObject(drawable, pos.X, pos.Y, alienSize, alienSize)
	o.Rot = wo.DegToRad(180)
	return &Alien{
		Object:   o,
		drawable: drawable,
		kind:     kind,
	}
}

func (a *Alien) Move(v pixel.Vec) {
	if a != nil {
		a.Object.Move(v)
	}
}

func (a *Alien) Draw(canvas *pixelgl.Canvas) {
	if a != nil {
		a.drawable.Sheet.SetFrame(a.kind)
		a.Object.Draw(canvas)
	}
}
