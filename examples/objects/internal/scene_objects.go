package internal

import (
	"math/rand"
	"time"

	"math"

	"github.com/explodes/go-wo"
	"github.com/explodes/go-wo/wobj"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var _ wo.Scene = &scene{}

type scene struct {
	time   float64
	bounds pixel.Rect
	rng    *rand.Rand

	drawable wobj.Drawable

	objects *wobj.Objects
}

func (w *World) newObjectsScene(canvas *pixelgl.Canvas) (wo.Scene, error) {

	sprite, err := w.loader.Sprite("img/ship_512.png")
	if err != nil {
		return nil, err
	}
	spriteDrawable := wobj.NewSpriteDrawable(sprite)

	scene := &scene{
		bounds:   canvas.Bounds(),
		objects:  wobj.NewObjects(),
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
		drawable: spriteDrawable,
	}
	return scene, nil
}

func (s *scene) Update(dt float64, input wo.Input) wo.SceneResult {
	s.time += dt

	if s.objects.Size() < 100 {
		s.addObject()
	}

	s.objects.Update(dt)

	return wo.SceneResultNone
}

func (s *scene) addObject() {

	size := s.norm(5, 8)
	speed := s.norm(10, 500)
	angle := -math.Pi + 2*math.Pi*s.rng.Float64()

	o := &wobj.Object{
		Tag:  "ship",
		Pos:  s.bounds.Center(),
		Size: pixel.V(size, size),

		Drawable: s.drawable,

		Velocity: pixel.V(0, speed).Rotated(angle),

		Steps: wobj.MakeBehaviors(
			wobj.Movement,
		),
		PostSteps: wobj.MakeBehaviors(
			wobj.FaceDirectionOffset(wo.DegToRad(-45)),
			func(source *wobj.Object, dt float64) {
				objBounds := source.Bounds()
				switch {
				case objBounds.Min.X <= s.bounds.Min.X:
					source.Velocity = pixel.V(-source.Velocity.X, source.Velocity.Y)
				case objBounds.Max.X >= s.bounds.Max.X:
					source.Velocity = pixel.V(-source.Velocity.X, source.Velocity.Y)
				}
				switch {
				case objBounds.Min.Y <= s.bounds.Min.Y:
					source.Velocity = pixel.V(source.Velocity.X, -source.Velocity.Y)
				case objBounds.Max.Y >= s.bounds.Max.Y:
					source.Velocity = pixel.V(source.Velocity.X, -source.Velocity.Y)
				}
			},
		),
	}
	s.objects.Add(o)
}

func (s *scene) Draw(canvas *pixelgl.Canvas) {
	s.objects.Draw(canvas)
}

func (s *scene) norm(min, max float64) float64 {
	if min == max {
		return min
	}

	stddev := min + (max-min)/3

	mean := (max + min) / 2
	r := s.rng.NormFloat64()*stddev + mean
	return math.Max(min, math.Min(max, r))
}
