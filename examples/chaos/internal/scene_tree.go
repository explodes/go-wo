package internal

import (
	"math/rand"

	"fmt"

	"image/color"

	"github.com/explodes/go-wo"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

const (
	maxTreeDepth       = 16
	branchLength       = 200
	branchLengthFactor = 0.8
)

type treeScene struct {
	w        *World
	rng      *rand.Rand
	bounds   pixel.Rect
	infoText *text.Text

	dirty bool
	im    *imdraw.IMDraw

	keyDelay float64

	depth     int
	angle     float64
	colorized bool
}

func (w *World) newTreeScene(canvas *pixelgl.Canvas) (wo.Scene, error) {

	infoFont, err := w.loader.FontFace("fonts/SourceSansPro-Regular.ttf", 16)
	if err != nil {
		return nil, err
	}
	defer infoFont.Close()
	infoText := text.New(pixel.V(0, 0), text.NewAtlas(infoFont, asciiPlusArrows))

	scene := &treeScene{
		w:        w,
		rng:      w.rng,
		bounds:   canvas.Bounds(),
		infoText: infoText,
		im:       imdraw.New(nil),
		depth:    3,
		angle:    45,
		dirty:    true,
	}
	return scene, nil
}

func (s *treeScene) Update(dt float64, input wo.Input) wo.SceneResult {
	s.keyDelay -= dt

	if s.keyDelay <= 0 {
		s.increaseEffectInt(input, pixelgl.KeyUp, pixelgl.KeyDown, 1, maxTreeDepth, 1, &s.depth)
		s.increaseEffectWrapped(input, pixelgl.KeyRight, pixelgl.KeyLeft, -180, 180, 1, &s.angle)
		if input.JustPressed(pixelgl.KeySpace) {
			s.colorized = !s.colorized
			s.keyDelay = keyDelay
			s.dirty = true
		}
	}
	if s.dirty {
		s.render()
		s.dirty = false
	}
	return s.w.maybeSelectScene(input)
}

func (s *treeScene) increaseEffectWrapped(input wo.Input, keyIncrease, keyDecrease pixelgl.Button, min, max, delta float64, value *float64) {
	if input.Pressed(keyIncrease) {
		*value += delta
		if *value > max {
			*value = min
		}
		s.dirty = true
		s.keyDelay = keyDelay
	} else if input.Pressed(keyDecrease) {
		*value -= delta
		if *value < min {
			*value = max
		}
		s.dirty = true
		s.keyDelay = keyDelay
	}
}

func (s *treeScene) increaseEffectInt(input wo.Input, keyIncrease, keyDecrease pixelgl.Button, min, max, delta int, value *int) {
	if input.Pressed(keyIncrease) {
		*value += delta
		if *value > max {
			*value = max
		}
		s.dirty = true
		s.keyDelay = keyDelay
	} else if input.Pressed(keyDecrease) {
		*value -= delta
		if *value < min {
			*value = min
		}
		s.dirty = true
		s.keyDelay = keyDelay
	}
}

func (s *treeScene) render() {
	im := s.im
	im.Clear()
	im.Color = colornames.White
	im.Push(pixel.V(width/2, 0), pixel.V(width/2, branchLength))
	im.Line(1)

	mat := pixel.IM.Moved(pixel.V(width/2, 0))
	s.renderTree(im, mat, 0, branchLength)
}

func (s *treeScene) renderTree(im *imdraw.IMDraw, mat pixel.Matrix, depth int, len float64) {
	if depth > s.depth {
		return
	}

	angle := wo.DegToRad(s.angle)
	branch := pixel.V(0, len)

	// draw
	if s.colorized {
		var r, g uint8
		if depth%2 == 0 {
			r = 255
			g = 0
		} else {
			r = 0
			g = 255
		}
		im.Color = color.RGBA{A: 255, R: r, G: g, B: uint8((1 - float64(depth)/float64(s.depth)) * 255.0)}
	}
	im.Push(mat.Project(pixel.ZV), mat.Project(branch))
	im.Line(1)

	// render deeper
	s.renderTree(im, pixel.IM.Rotated(pixel.ZV, float64(depth)*angle).Moved(mat.Project(branch)), depth+1, len*branchLengthFactor)
	if depth%2 == 0 {
		s.renderTree(im, pixel.IM.Rotated(pixel.ZV, float64(depth)*-angle).Moved(mat.Project(branch)), depth+1, len*branchLengthFactor)
	}
}

func (s *treeScene) Draw(canvas *pixelgl.Canvas) {
	s.im.Draw(canvas)
	drawText(
		canvas,
		pixel.V(10, canvas.Bounds().H()-10),
		s.infoText,
		fmt.Sprintf("↑↓) depth: %d", s.depth),
		fmt.Sprintf("←→) angle: %0.f", s.angle),
		fmt.Sprintf("space) colorized: %v", s.colorized),
	)
}
