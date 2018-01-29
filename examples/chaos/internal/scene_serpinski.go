package internal

import (
	"math/rand"

	"image/color"

	"math"

	"github.com/explodes/go-wo"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	triangleVertexRadius = 2
	pointRadius          = 1

	pointsPerStep = 40
)

type sierpinskiScene struct {
	w   *World
	rng *rand.Rand

	bounds pixel.Rect
	im     *imdraw.IMDraw

	triangle [3]pixel.Vec
	colors   [3]color.Color
	seed     pixel.Vec
}

func (w *World) newSierpinskiScene(canvas *pixelgl.Canvas) (wo.Scene, error) {
	scene := &sierpinskiScene{
		w:   w,
		rng: w.rng,

		bounds: canvas.Bounds(),
		im:     imdraw.New(nil),

		colors: [3]color.Color{colornames.Red, colornames.Green, colornames.Blue},
	}
	scene.triangle = scene.randomTriangle(w.rng, canvas.Bounds())
	scene.seed = scene.randomPoint(w.rng, canvas.Bounds())
	scene.drawTriangle()
	return scene, nil
}

func (s *sierpinskiScene) Update(dt float64, input wo.Input) wo.SceneResult {
	for i := 0; i < pointsPerStep; i++ {
		s.addNewPoint()
	}
	return s.w.maybeSelectScene(input)
}

func (s *sierpinskiScene) drawTriangle() {
	for i := 0; i < 3; i++ {
		s.im.Color = s.colors[i]
		s.im.Push(s.triangle[i])
		s.im.Circle(triangleVertexRadius, 0)
	}

}

func (s *sierpinskiScene) addNewPoint() {
	n := s.rng.Intn(3)
	mid := s.midpoint(s.triangle[n], s.seed)

	s.im.Color = s.colorFor(n, mid)
	s.im.Push(mid)
	s.im.Circle(pointRadius, 0)

	s.seed = mid
}

func (s *sierpinskiScene) Draw(canvas *pixelgl.Canvas) {
	s.im.Draw(canvas)
}

func (s *sierpinskiScene) colorFor(n int, pt pixel.Vec) color.Color {
	return color.RGBA{
		R: s.colorComponent(s.distanceOf(0, pt)),
		G: s.colorComponent(s.distanceOf(1, pt)),
		B: s.colorComponent(s.distanceOf(2, pt)),
		A: 1,
	}
}

func (s *sierpinskiScene) distanceOf(n int, pt pixel.Vec) float64 {
	host := s.triangle[n]

	other1 := s.triangle[(n+1)%3]
	other2 := s.triangle[(n+1)%3]
	maxDistance := math.Max(
		s.positive(host.Sub(other1).Len()),
		s.positive(host.Sub(other2).Len()),
	)
	if maxDistance == 0 {
		return 0
	}

	distance := s.positive(host.Sub(pt).Len()) / maxDistance

	if distance > 1 {
		return 1
	}
	return distance
}

func (s *sierpinskiScene) colorComponent(p float64) uint8 {
	return uint8(255 * p)
}

func (s *sierpinskiScene) positive(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

func (s *sierpinskiScene) randomTriangle(rng *rand.Rand, r pixel.Rect) [3]pixel.Vec {
	return [3]pixel.Vec{
		s.randomPoint(rng, r),
		s.randomPoint(rng, r),
		s.randomPoint(rng, r),
	}
}

func (s *sierpinskiScene) randomPoint(rng *rand.Rand, r pixel.Rect) pixel.Vec {
	return pixel.V(
		rng.Float64()*r.W()+r.Min.X,
		rng.Float64()*r.H()+r.Min.Y,
	)
}

func (s *sierpinskiScene) midpoint(a, b pixel.Vec) pixel.Vec {
	return pixel.V(
		(a.X+b.X)/2,
		(a.Y+b.Y)/2,
	)
}
