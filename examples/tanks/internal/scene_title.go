package internal

import (
	"math/rand"
	"time"

	"image/color"

	"fmt"

	"github.com/explodes/go-wo"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

var _ wo.Scene = &gameScene{}

var (
	instructions = []string{"Blue will rotate the tank with (A).",
		"Red will rotate the tank with (L).",
		"Your tank will automatically fire.",
		"Last one standing wins.",
	}
)

type titleScene struct {
	rng *rand.Rand

	title      *text.Text
	titlePos   []pixel.Vec
	titleColor []color.Color

	help         *text.Text
	score        *text.Text
	instructions *text.Text
}

func (w *World) newTitleScene(canvas *pixelgl.Canvas) (wo.Scene, error) {
	instructionsFont, err := w.loader.FontFace("fonts/Lekton-Regular.ttf", 12)
	if err != nil {
		return nil, err
	}
	defer instructionsFont.Close()
	instructionsText := text.New(pixel.V(800, 40), text.NewAtlas(instructionsFont, text.ASCII))
	instructionsText.Color = colornames.White
	for _, line := range instructions {
		instructionsText.Dot.X -= instructionsText.BoundsOf(line).W() + 10
		fmt.Fprintln(instructionsText, line)
	}

	helpFont, err := w.loader.FontFace("fonts/DampfPlatzs.ttf", 24)
	if err != nil {
		return nil, err
	}
	defer helpFont.Close()
	helpText := text.New(pixel.V(canvas.Bounds().W()/2, 10), text.NewAtlas(helpFont, text.ASCII))
	helpText.Color = colornames.White
	helpText.WriteString("press space to battle")

	titleFont, err := w.loader.FontFace("fonts/DampfPlatz.ttf", 240)
	if err != nil {
		return nil, err
	}
	defer titleFont.Close()
	titleText := text.New(pixel.V(0, 0), text.NewAtlas(titleFont, text.ASCII))
	titleText.Color = colornames.White
	titleText.WriteString("Tanks")

	scoreFont, err := w.loader.FontFace("fonts/BlackKnightFLF.ttf", 36)
	if err != nil {
		return nil, err
	}
	defer scoreFont.Close()
	scoreText := text.New(pixel.V(10, canvas.Bounds().H()-36), text.NewAtlas(scoreFont, text.ASCII))
	scoreText.Color = colornames.Blue
	scoreText.WriteString(fmt.Sprintf("Blue: %d", w.blueScore))
	scoreText.Color = colornames.White
	scoreText.WriteString(" - ")
	scoreText.Color = colornames.Red
	scoreText.WriteString(fmt.Sprintf("Red: %d", w.redScore))

	scene := &titleScene{
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
		title:    titleText,
		titlePos: make([]pixel.Vec, 5),
		titleColor: []color.Color{
			colornames.Red,
			colornames.Lightblue,
			colornames.Coral,
			colornames.Cornflowerblue,
			colornames.White,
		},
		help:         helpText,
		score:        scoreText,
		instructions: instructionsText,
	}

	return scene, nil
}

func (s *titleScene) Update(dt float64, input wo.Input) wo.SceneResult {
	if input.Pressed(pixelgl.KeySpace) {
		return gotoBattle
	}
	for i := 0; i < len(s.titlePos); i++ {
		s.titlePos[i] = pixel.V(-3+s.rng.Float64()*6, -3+s.rng.Float64()*6)
	}
	return wo.SceneResultNone
}

func (s *titleScene) Draw(canvas *pixelgl.Canvas) {
	for i := 0; i < len(s.titlePos); i++ {
		textColor := s.titleColor[i]
		offset := s.titlePos[i]
		mat := pixel.IM.Moved(canvas.Bounds().Center().Sub(s.title.Bounds().Center()).Add(offset))
		s.title.DrawColorMask(canvas, mat, textColor)
	}
	s.score.Draw(canvas, pixel.IM)
	s.help.Draw(canvas, pixel.IM.Moved(pixel.V(-s.help.Bounds().W()/2, 0)))
	s.instructions.Draw(canvas, pixel.IM)
}
