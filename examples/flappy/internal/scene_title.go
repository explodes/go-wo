package internal

import (
	"fmt"
	"image/color"

	"github.com/explodes/go-wo"
	"github.com/explodes/go-wo/examples/flappy/res"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	beginText = "(press any key to begin)"
)

type titleScene struct {
	game      *FlappyWorld
	largeText *text.Text
	smallText *text.Text
	bg        *imdraw.IMDraw
	log       *logrus.Entry
}

func (g *FlappyWorld) createTitleScene(canvas *pixelgl.Canvas) (wo.Scene, error) {
	fontSrc, err := res.Load("fonts/Flappy.ttf")
	if err != nil {
		return nil, errors.Errorf("could not find font: %v", err)
	}
	font, err := truetype.Parse(fontSrc)
	if err != nil {
		return nil, errors.Errorf("could not parse font: %v", err)
	}
	large := truetype.NewFace(font, &truetype.Options{Size: 48})
	defer large.Close()
	small := truetype.NewFace(font, &truetype.Options{Size: 16})
	defer small.Close()
	largeText := text.New(pixel.V(0, 0), text.NewAtlas(large, text.ASCII))
	largeText.WriteString(title)
	smallText := text.New(pixel.V(0, 0), text.NewAtlas(small, text.ASCII))

	scene := &titleScene{
		game:      g,
		largeText: largeText,
		smallText: smallText,
		bg:        createDebugIMDraw(canvas),
		log:       g.log.WithField("scene", "title"),
	}
	return scene, nil
}

func (s *titleScene) Update(dt float64, input wo.Input) wo.SceneResult {
	if input.Typed() != "" {
		return SceneResultGoToGame
	}
	return wo.SceneResultNone
}

func (s *titleScene) Draw(canvas *pixelgl.Canvas) {

	const pad = 10

	s.bg.Draw(canvas)

	s.smallText.Clear()
	s.smallText.WriteString(fmt.Sprintf("\r%s", beginText))

	titleBounds := s.largeText.Bounds()
	helpBounds := s.smallText.Bounds()

	s.largeText.Draw(canvas, pixel.IM.Moved(pixel.V(width-pad-titleBounds.Max.X, pad+titleBounds.Min.Y+helpBounds.H()+pad)))

	s.smallText.Draw(canvas, pixel.IM.Moved(pixel.V(width-pad-helpBounds.Max.X, pad+helpBounds.Min.Y)))

	if s.game.lastScore != noPreviousScore {
		s.smallText.Clear()
		s.smallText.WriteString(fmt.Sprintf("\rPrevious score: %dpts", s.game.lastScore))
		scoreBounds := s.smallText.Bounds()
		s.smallText.Draw(canvas, pixel.IM.Moved(pixel.V(10, canvas.Bounds().H()-scoreBounds.H()-10)))
	}
}

func createDebugIMDraw(canvas *pixelgl.Canvas) *imdraw.IMDraw {

	const div = float64(19)

	im := imdraw.New(nil)

	x, y, w, h := wo.Shape(canvas.Bounds())

	for i := float64(0); i <= div; i += 1 {
		for j := float64(0); j <= div; j += 1 {

			dx := x + i*w/div
			dy := y + j*h/div

			im.Color = pixel.RGBA{R: i / div, G: j / div, B: 0, A: 0.1}
			im.SetColorMask(color.Alpha{A: 64})
			im.Push(pixel.V(dx, dy))
			im.Circle(5, 0)
		}
	}
	return im
}
