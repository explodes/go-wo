package internal

import (
	"math/rand"

	"image/color"

	"github.com/explodes/go-wo"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

const (
	minTriSpeed = 0.5
	maxTriSpeed = 3
	triSize     = 35

	textX, textY = 724, 171
)

var (
	triLineColor = color.Alpha{A: 64}
)

type titleScene struct {
	w   *World
	rng *rand.Rand

	bounds   pixel.Rect
	infoText *text.Text

	background *imdraw.IMDraw
	triangles  []*triangle
}

type triangle struct {
	vertices []pixel.Vec
	im       *imdraw.IMDraw

	velocity []pixel.Vec
}

func newTriangle(rng *rand.Rand, r pixel.Rect) *triangle {

	center := pixel.V(r.Min.X+r.W()*rng.Float64(), r.Min.Y+r.H()*rng.Float64())

	n := 3 + rng.Intn(5-3)

	vex := make([]pixel.Vec, n)
	for i := 0; i < n; i++ {
		vex[i] = center.Add(pixel.V(rng.Float64()*triSize, rng.Float64()*triSize))
	}
	vel := make([]pixel.Vec, n)
	for i := 0; i < n; i++ {
		vel[i] = pixel.V(1, 1).Scaled(minTriSpeed + rng.Float64()*maxTriSpeed).Rotated(wo.DegToRad(360 * rng.Float64()))
	}

	return &triangle{
		vertices: vex,
		velocity: vel,
		im:       imdraw.New(nil),
	}
}

func (t *triangle) update(dt float64) {
	for i := 0; i < len(t.vertices); i++ {
		t.vertices[i] = t.vertices[i].Add(t.velocity[i].Scaled(dt))
	}
}

func (t triangle) draw(canvas *pixelgl.Canvas) {
	t.im.Clear()
	for i := 1; i < len(t.vertices); i++ {
		for j := 0; j < i; j++ {
			t.im.Color = triLineColor
			t.im.Push(t.vertices[i], t.vertices[j])
			t.im.Line(1)
		}
	}
	for i := 0; i < len(t.vertices); i++ {
		t.im.Color = colornames.Lightblue
		t.im.Push(t.vertices[i])
		t.im.Circle(2, 0)
	}
	t.im.Draw(canvas)
}

func (w *World) newTitleScene(canvas *pixelgl.Canvas) (wo.Scene, error) {

	infoFont, err := w.loader.FontFace("fonts/Lekton-Regular.ttf", 24)
	if err != nil {
		return nil, err
	}
	defer infoFont.Close()
	infoText := text.New(pixel.V(0, 0), text.NewAtlas(infoFont, asciiPlusArrows))
	infoText.Color = colornames.Green

	triangles := make([]*triangle, 40)
	for i := 0; i < len(triangles); i++ {
		triangles[i] = newTriangle(w.rng, canvas.Bounds())
	}

	background := imdraw.New(nil)
	background.Color = colornames.Black
	background.Push(
		pixel.V(textX, textY),
		pixel.V(canvas.Bounds().Max.X, 0),
	)
	background.Rectangle(0)

	scene := &titleScene{
		w:          w,
		rng:        w.rng,
		bounds:     canvas.Bounds(),
		infoText:   infoText,
		triangles:  triangles,
		background: background,
	}
	return scene, nil
}

func (s *titleScene) Update(dt float64, input wo.Input) wo.SceneResult {
	for _, triangle := range s.triangles {
		triangle.update(dt)
	}
	return s.w.maybeSelectScene(input)
}

func (s *titleScene) Draw(canvas *pixelgl.Canvas) {
	for _, triangle := range s.triangles {
		triangle.draw(canvas)
	}

	s.background.Draw(canvas)

	drawText(
		canvas,
		pixel.V(724, 193),
		s.infoText,
		"     CHAOS",
		"~ ~ ~ ~ ~ ~ ~ ~",
		"1) Sierpinski",
		"2) Circles",
		"3) Tree",
		"0) Title Screen",
		"esc) Quit",
	)
}
