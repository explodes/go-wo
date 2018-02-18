package internal

import (
	"math/rand"
	"time"

	"github.com/explodes/go-wo"
	"github.com/explodes/go-wo/wobj"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var _ wo.Scene = &scene{}

const (
	tankRotatesPerSecond = 0.5
	tankSpeed            = 175

	autoShotPerSecond = 0.5

	bulletSpeed = 560
)

var (
	tankRotateOffset = wo.DegToRad(90)
)

type scene struct {
	time   float64
	bounds pixel.Rect
	rng    *rand.Rand
	input  wo.Input

	player1 *wobj.Object
	player2 *wobj.Object

	shot wobj.Drawable

	blueShotDelay float64
	redShotDelay  float64

	objects *wobj.Objects
}

func (w *World) newScene(canvas *pixelgl.Canvas) (wo.Scene, error) {

	shotSprite, err := w.loader.Sprite("img/shot.png")
	if err != nil {
		return nil, err
	}

	tankSheet, err := w.loader.SpriteSheet("img/tanks3.png", wo.SpriteSheetOptions{
		Width:   149,
		Height:  166,
		Columns: 1,
		Rows:    2,
	})
	if err != nil {
		return nil, err
	}
	tank1Drawable := wobj.NewSpriteSheetDrawable(tankSheet)

	tankSheet, err = w.loader.SpriteSheet("img/tanks3.png", wo.SpriteSheetOptions{
		Width:   149,
		Height:  166,
		Columns: 1,
		Rows:    2,
	})
	if err != nil {
		return nil, err
	}
	tank2Drawable := wobj.NewSpriteSheetDrawable(tankSheet)

	tank1Drawable.Sheet.SetFrame(0)
	tank2Drawable.Sheet.SetFrame(1)

	scene := &scene{
		bounds:  canvas.Bounds(),
		objects: wobj.NewObjects(),
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
		shot:    wobj.NewSpriteDrawable(shotSprite),
	}

	player1 := &wobj.Object{
		Tag:      "player1",
		Pos:      pixel.V(100, 200),
		Size:     pixel.V(170*3/10, 200*3/10),
		Drawable: tank1Drawable,
		Rot:      tankRotateOffset + wo.DegToRad(135),

		Steps: wobj.MakeBehaviors(
			scene.behaviorBlueRotateOnButton,
		),
		PostSteps: wobj.MakeBehaviors(
			scene.behaviorReflectInBounds,
		),
	}
	scene.player1 = player1

	player2 := &wobj.Object{
		Tag:      "player2",
		Pos:      pixel.V(700, 200),
		Size:     pixel.V(170*3/10, 200*3/10),
		Drawable: tank2Drawable,
		Rot:      tankRotateOffset + wo.DegToRad(-45),

		Steps: wobj.MakeBehaviors(
			scene.behaviorRedRotateOnButton,
		),
		PostSteps: wobj.MakeBehaviors(
			scene.behaviorReflectInBounds,
		),
	}
	scene.player2 = player2

	scene.objects.Add(player1)
	scene.objects.Add(player2)

	return scene, nil
}

func (s *scene) Update(dt float64, input wo.Input) wo.SceneResult {
	s.time += dt
	s.blueShotDelay += dt
	s.redShotDelay += dt
	s.input = input

	s.objects.Update(dt)

	return wo.SceneResultNone
}

func (s *scene) Draw(canvas *pixelgl.Canvas) {
	s.objects.Draw(canvas)
}

func (s *scene) behaviorBlueRotateOnButton(source *wobj.Object, dt float64) {
	if s.input.Pressed(pixelgl.KeyA) {
		// rotate
		source.Rot += wo.DegToRad(-tankRotatesPerSecond*360) * dt
		s.blueShotDelay = 0
	} else {
		source.Velocity = pixel.V(tankSpeed, 0).Rotated(source.Rot - tankRotateOffset)
		wobj.Movement(source, dt)
		if s.blueShotDelay > 1.0/autoShotPerSecond {
			s.spawnBlueShots()
			s.blueShotDelay = 0
		}
	}
}

func (s *scene) behaviorRedRotateOnButton(source *wobj.Object, dt float64) {
	if s.input.Pressed(pixelgl.KeyL) {
		// rotate
		source.Rot += wo.DegToRad(-tankRotatesPerSecond*360) * dt
		s.redShotDelay = 0
	} else {
		source.Velocity = pixel.V(tankSpeed, 0).Rotated(source.Rot - tankRotateOffset)
		wobj.Movement(source, dt)
		if s.redShotDelay > 1.0/autoShotPerSecond {
			s.spawnRedShots()
			s.redShotDelay = 0
		}
	}
}

func (s *scene) spawnBlueShots() {

	bounds := s.player1.Bounds()
	pos1 := bounds.Center().Add(pixel.V(bounds.W()/2, 2).Rotated(s.player1.Rot - tankRotateOffset))
	pos2 := bounds.Center().Add(pixel.V(bounds.W()/2, -8).Rotated(s.player1.Rot - tankRotateOffset))

	blueBullet1 := &wobj.Object{
		Tag:      "blueBullet",
		Pos:      pos1,
		Size:     pixel.V(8, 8),
		Drawable: s.shot,
		Velocity: pixel.V(bulletSpeed, 0).Rotated(s.player1.Rot - tankRotateOffset),
		Steps: wobj.MakeBehaviors(
			wobj.Movement,
		),
		PostSteps: wobj.MakeBehaviors(
			s.behaviorRemoveOutOfBounds,
		),
	}
	blueBullet2 := &wobj.Object{
		Tag:      "blueBullet",
		Pos:      pos2,
		Size:     pixel.V(8, 8),
		Drawable: s.shot,
		Velocity: pixel.V(bulletSpeed, 0).Rotated(s.player1.Rot - tankRotateOffset),
		Steps: wobj.MakeBehaviors(
			wobj.Movement,
		),
		PostSteps: wobj.MakeBehaviors(
			s.behaviorRemoveOutOfBounds,
		),
	}
	s.objects.Add(blueBullet1)
	s.objects.Add(blueBullet2)
}

func (s *scene) spawnRedShots() {

	bounds := s.player2.Bounds()
	offset := pixel.V(bounds.H()/2, -8).Rotated(s.player2.Rot - tankRotateOffset)
	pos := bounds.Center().Add(offset)

	redBullet := &wobj.Object{
		Tag:      "redBullet",
		Pos:      pos,
		Size:     pixel.V(14, 14),
		Drawable: s.shot,
		Velocity: pixel.V(bulletSpeed, 0).Rotated(s.player2.Rot - tankRotateOffset),
		Steps: wobj.MakeBehaviors(
			wobj.Movement,
		),
		PostSteps: wobj.MakeBehaviors(
			s.behaviorRemoveOutOfBounds,
		),
	}
	s.objects.Add(redBullet)
}

func (s *scene) behaviorRemoveOutOfBounds(source *wobj.Object, dt float64) {
	if !source.Collides(s.bounds) {
		s.objects.Remove(source)
	}
}

func (s *scene) behaviorReflectInBounds(source *wobj.Object, dt float64) {
	objBounds := source.Bounds()
	switch {
	case objBounds.Min.X <= s.bounds.Min.X:
		source.Velocity = pixel.V(-source.Velocity.X, source.Velocity.Y)
		source.Rot = source.Velocity.Angle() + tankRotateOffset
	case objBounds.Max.X >= s.bounds.Max.X:
		source.Velocity = pixel.V(-source.Velocity.X, source.Velocity.Y)
		source.Rot = source.Velocity.Angle() + tankRotateOffset
	}
	switch {
	case objBounds.Min.Y <= s.bounds.Min.Y:
		source.Velocity = pixel.V(source.Velocity.X, -source.Velocity.Y)
		source.Rot = source.Velocity.Angle() + tankRotateOffset
	case objBounds.Max.Y >= s.bounds.Max.Y:
		source.Velocity = pixel.V(source.Velocity.X, -source.Velocity.Y)
		source.Rot = source.Velocity.Angle() + tankRotateOffset
	}
}
