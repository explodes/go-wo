package internal

import (
	"math"

	"fmt"

	"github.com/explodes/go-wo"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/usedbytes/hsv"
	"golang.org/x/image/colornames"
)

const (
	maxCircles = 750

	defaultCosEffect = 45
	defaultPowEffect = 0.85
	cosEffectDelta   = 0.1
	powEffectDelta   = 0.01

	keyDelay = 0
)

var (
	asciiPlusArrows = append(text.ASCII, []rune{'↑', '↓', '←', '→'}...)
)

type circlesScene struct {
	w *World

	time float64

	bounds pixel.Rect
	im     *imdraw.IMDraw

	infoText *text.Text

	keyDelay  float64
	cosEffect float64
	powEffect float64
	colorized bool
}

func (w *World) newCirclesScene(canvas *pixelgl.Canvas) (wo.Scene, error) {

	infoFont, err := w.loader.FontFace("fonts/SourceSansPro-Regular.ttf", 16)
	if err != nil {
		return nil, err
	}
	defer infoFont.Close()
	infoText := text.New(pixel.V(0, 0), text.NewAtlas(infoFont, asciiPlusArrows))

	scene := &circlesScene{
		w:         w,
		bounds:    canvas.Bounds(),
		im:        imdraw.New(nil),
		infoText:  infoText,
		cosEffect: defaultCosEffect,
		powEffect: defaultPowEffect,
	}
	return scene, nil
}

func (s *circlesScene) Update(dt float64, input wo.Input) wo.SceneResult {
	s.time += dt
	s.keyDelay -= dt

	if s.keyDelay < 0 {
		s.increaseEffect(input, pixelgl.KeyUp, pixelgl.KeyDown, cosEffectDelta, &s.cosEffect)
		s.increaseEffect(input, pixelgl.KeyRight, pixelgl.KeyLeft, powEffectDelta, &s.powEffect)
		if input.JustPressed(pixelgl.KeySpace) {
			s.colorized = !s.colorized
			s.keyDelay = keyDelay
		}
	}

	s.renderCircles()

	return s.w.maybeSelectScene(input)
}

func (s *circlesScene) increaseEffect(input wo.Input, keyIncrease, keyDecrease pixelgl.Button, delta float64, value *float64) {
	if input.Pressed(keyIncrease) {
		*value += delta
		s.keyDelay = keyDelay
	} else if input.Pressed(keyDecrease) {
		*value -= delta
		if *value < delta {
			*value = delta
		}
		s.keyDelay = keyDelay
	}
}

func (s *circlesScene) renderCircles() {
	dt := math.Cos(s.time)*s.cosEffect + s.cosEffect*1.12

	s.im.Clear()
	s.im.Color = colornames.White
	center := s.bounds.Center()
	maxRadius := maxRadius(s.bounds)

	for i := 0; i < maxCircles; i++ {
		t := float64(i)
		t = math.Pow(t, s.powEffect)
		radius := t * dt
		if radius > maxRadius {
			break
		}
		if s.colorized {
			s.im.Color = hsv.HSVColor{S: 255, V: 255, H: uint16(float64(i) / maxCircles * 360)}
		}
		s.im.Push(center)
		s.im.Circle(radius, 1)
	}
}

func maxRadius(r pixel.Rect) float64 {
	w := r.W()
	h := r.H()
	return math.Sqrt(w*w + h*h)
}

func (s *circlesScene) Draw(canvas *pixelgl.Canvas) {
	s.im.Draw(canvas)
	drawText(
		canvas,
		pixel.V(10, canvas.Bounds().H()-10),
		s.infoText,
		fmt.Sprintf("cos (↑↓) : %f", s.cosEffect),
		fmt.Sprintf("pow (←→): %f", s.powEffect),
		fmt.Sprintf("colorized (space): %v", s.colorized),
	)
}
