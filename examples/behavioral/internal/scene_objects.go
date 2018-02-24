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

const (
	shipTag            = "ship"
	numShips           = 100
	sizeMin, sizeMax   = 5, 12
	speedMin, speedMax = 10, 150
)

type scene struct {
	time   float64
	bounds pixel.Rect
	rng    *rand.Rand
	input  wo.Input

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
	s.input = input

	if s.objects.Len() < numShips {
		s.addObject()
	}

	s.objects.Update(dt)

	return wo.SceneResultNone
}

func (s *scene) addObject() {

	size := s.norm(sizeMin, sizeMax)
	speed := s.norm(speedMin, speedMax)
	angle := 2 * math.Pi * s.rng.Float64()

	var onBorderCollision wobj.Behavior
	if s.rng.Float64() < .5 {
		onBorderCollision = s.behaviorReflectInBounds
	} else {
		onBorderCollision = s.behaviorRemoveOutOfBounds
	}

	o := &wobj.Object{
		Tag:       shipTag,
		Pos:       s.bounds.Center(),
		Size:      pixel.V(size, size),
		Drawable:  s.drawable,
		Velocity:  pixel.V(0, speed).Rotated(angle),
		RotNormal: wo.DegToRad(-45),
		Steps: wobj.MakeBehaviors(
			s.behaviorShipInput,
		),
		PostSteps: wobj.MakeBehaviors(
			wobj.FaceDirection,
			onBorderCollision,
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
	stddev99 := (max - min) / 3.2165 / 2
	mean := (max + min) / 2
	r := s.rng.NormFloat64()*stddev99 + mean
	return math.Max(min, math.Min(max, r))
}

func (s *scene) behaviorShipInput(source *wobj.Object, dt float64) {
	switch {
	case s.input.Pressed(pixelgl.KeyA):
		mag := source.Velocity.Len()
		angle := wo.DegToRad(45)
		source.Velocity = pixel.V(mag, 0).Rotated(angle)
	case s.input.Pressed(pixelgl.KeySpace):
		mag := source.Velocity.Len()
		angle := source.Velocity.Angle() + wo.DegToRad(360)*dt
		source.Velocity = pixel.V(mag, 0).Rotated(angle)
	default:
		wobj.Movement(source, dt)
	}
}

func (s *scene) behaviorReflectInBounds(source *wobj.Object, dt float64) {
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
}

func (s *scene) behaviorRemoveOutOfBounds(source *wobj.Object, dt float64) {
	if !wo.Collision(source.Bounds(), s.bounds) {
		s.objects.Remove(source)
	}
}
