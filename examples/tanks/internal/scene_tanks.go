package internal

import (
	"math/rand"
	"time"

	"strconv"

	"fmt"

	"image/color"

	_ "image/jpeg"

	"github.com/explodes/go-wo"
	"github.com/explodes/go-wo/wobj"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

var _ wo.Scene = &gameScene{}

type Phase uint8

const (
	phaseCountdown Phase = iota
	phaseBattle
	phaseBlueVictory
	phaseRedVictory
)

const (
	tankRotatesPerSecond  = 0.5
	tankSpeed             = 185
	tankWidth, tankHeight = 170 * 3 / 10, 200 * 3 / 10

	victoryMessageDuration = 3

	autoShotPerSecond = 0.5

	bulletSpeed = 580

	tagBackground = "background"
	tagBluePlayer = "bluePlayer"
	tagBlueBullet = "blueBullet"
	tagRedPlayer  = "redPlayer"
	tagRedBullet  = "redBullet"
)

const (
	layerBackground = iota
	layerTanks
	layerBullets
	numLayers
)

var (
	tankRotateOffset = wo.DegToRad(90)

	winningMessages = []string{
		"%s has become the champion",
		"%s is victorious",
		"%s was better",
	}

	countdownColors = []color.Color{
		colornames.Red,
		colornames.Blue,
		colornames.White,
	}
)

type gameScene struct {
	w      *World
	time   float64
	bounds pixel.Rect
	rng    *rand.Rand
	input  wo.Input

	phase Phase

	message *text.Text

	cannon *wo.Sound

	bluePlayer *wobj.Object
	redPlayer  *wobj.Object

	victoryTime float64

	shot wobj.Drawable

	blueShotDelay float64
	redShotDelay  float64

	layers wobj.Layers
}

func (w *World) newGameScene(canvas *pixelgl.Canvas) (wo.Scene, error) {

	countdownFont, err := w.loader.FontFace("fonts/DampfPlatzs.ttf", 42)
	if err != nil {
		return nil, err
	}
	defer countdownFont.Close()
	countdownText := text.New(pixel.V(canvas.Bounds().W()/2, 10), text.NewAtlas(countdownFont, text.ASCII))

	cannon, err := w.loader.Sound("wav", "sound/tank.wav")
	if err != nil {
		return nil, err
	}

	shotSprite, err := w.loader.Sprite("img/shot.png")
	if err != nil {
		return nil, err
	}

	dirtSprite, err := w.loader.Sprite("img/dirt.jpg")
	if err != nil {
		return nil, err
	}

	tankSheet, err := w.loader.SpriteSheet("img/tanks.png", wo.SpriteSheetOptions{
		Width:   149,
		Height:  166,
		Columns: 1,
		Rows:    2,
	})
	if err != nil {
		return nil, err
	}
	blueTankDrawable := wobj.NewSpriteSheetDrawable(tankSheet)

	tankSheet, err = w.loader.SpriteSheet("img/tanks.png", wo.SpriteSheetOptions{
		Width:   149,
		Height:  166,
		Columns: 1,
		Rows:    2,
	})
	if err != nil {
		return nil, err
	}
	tank2Drawable := wobj.NewSpriteSheetDrawable(tankSheet)

	blueTankDrawable.Sheet.SetFrame(0)
	tank2Drawable.Sheet.SetFrame(1)

	scene := &gameScene{
		w:       w,
		phase:   phaseCountdown,
		bounds:  canvas.Bounds(),
		layers:  wobj.NewLayers(numLayers),
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
		shot:    wobj.NewSpriteDrawable(shotSprite),
		cannon:  cannon,
		message: countdownText,
	}

	rot1 := wo.DegToRad(135)
	rot2 := wo.DegToRad(-45)
	if scene.rng.Float64() < 0.5 {
		rot1, rot2 = rot2, rot1
	}

	player1 := &wobj.Object{
		Tag:       tagBluePlayer,
		Pos:       pixel.V(100, height/2-tankHeight/2),
		Size:      pixel.V(tankWidth, tankHeight),
		Drawable:  blueTankDrawable,
		Rot:       rot1,
		RotNormal: tankRotateOffset,

		Steps: wobj.MakeBehaviors(
			scene.behaviorBlueRotateOnButton,
		),
		PostSteps: wobj.MakeBehaviors(
			wobj.ReflectWithin(scene),
			scene.behaviorBlueHitsRedBullet,
		),
	}
	scene.bluePlayer = player1
	scene.layers[layerTanks].Add(player1)

	player2 := &wobj.Object{
		Tag:       tagRedPlayer,
		Pos:       pixel.V(width-100-tankWidth, height/2-tankHeight/2),
		Size:      pixel.V(tankWidth, tankHeight),
		Drawable:  tank2Drawable,
		Rot:       rot2,
		RotNormal: tankRotateOffset,

		Steps: wobj.MakeBehaviors(
			scene.behaviorRedRotateOnButton,
		),
		PostSteps: wobj.MakeBehaviors(
			wobj.ReflectWithin(scene),
			scene.behaviorRedHitsBlueBullet,
		),
	}
	scene.redPlayer = player2
	scene.layers[layerTanks].Add(player2)

	dirt := &wobj.Object{
		Tag:      tagBackground,
		Size:     canvas.Bounds().Max,
		Drawable: wobj.NewSpriteDrawable(dirtSprite),
	}
	scene.layers[layerBackground].Add(dirt)

	return scene, nil
}

func (s *gameScene) Update(dt float64, input wo.Input) wo.SceneResult {
	if input.JustPressed(pixelgl.KeyEscape) {
		return wo.SceneResultWindowClosed
	}
	s.time += dt
	s.input = input

	switch s.phase {
	case phaseCountdown:

		countdownTime := s.time * 2
		if countdownTime >= 3 {
			s.phase = phaseBattle
			break
		}
		seconds := 3 - int(countdownTime)

		countdownColorIndex := 3 - seconds
		if countdownColorIndex < 0 {
			countdownColorIndex = 0
		}

		s.message.Clear()
		s.message.Color = countdownColors[countdownColorIndex]
		s.message.WriteString(strconv.Itoa(seconds))
	case phaseBattle:
		s.blueShotDelay += dt
		s.redShotDelay += dt
		s.layers.Update(dt)
	case phaseBlueVictory:
		fallthrough
	case phaseRedVictory:
		s.victoryTime -= dt
		if s.victoryTime <= 0 {
			return gotoTitle
		}
	}

	return wo.SceneResultNone
}

func (s *gameScene) Draw(canvas *pixelgl.Canvas) {
	s.layers.Draw(canvas)

	switch s.phase {
	case phaseBattle:
	case phaseBlueVictory:
		fallthrough
	case phaseRedVictory:
		fallthrough
	case phaseCountdown:
		s.message.Draw(canvas, pixel.IM.Moved(canvas.Bounds().Center().Sub(s.message.Bounds().Center())))
	}
}

func (s *gameScene) behaviorBlueRotateOnButton(source *wobj.Object, dt float64) {
	if s.input.Pressed(pixelgl.KeyA) {
		// rotate
		source.Rot += wo.DegToRad(-tankRotatesPerSecond*360) * dt
		s.blueShotDelay = 0
	} else {
		source.Velocity = pixel.V(tankSpeed, 0).Rotated(source.Rot)
		wobj.Movement(source, dt)
		if s.blueShotDelay > 1.0/autoShotPerSecond {
			s.spawnBlueShots()
			s.blueShotDelay = 0
		}
	}
}

func (s *gameScene) behaviorRedRotateOnButton(source *wobj.Object, dt float64) {
	if s.input.Pressed(pixelgl.KeyL) {
		// rotate
		source.Rot += wo.DegToRad(-tankRotatesPerSecond*360) * dt
		s.redShotDelay = 0
	} else {
		source.Velocity = pixel.V(tankSpeed, 0).Rotated(source.Rot)
		wobj.Movement(source, dt)
		if s.redShotDelay > 1.0/autoShotPerSecond {
			s.spawnRedShots()
			s.redShotDelay = 0
		}
	}
}

func (s *gameScene) spawnBlueShots() {

	bounds := s.bluePlayer.Bounds()
	pos1 := bounds.Center().Add(pixel.V(bounds.W()/2, 2).Rotated(s.bluePlayer.Rot))
	pos2 := bounds.Center().Add(pixel.V(bounds.W()/2, -8).Rotated(s.bluePlayer.Rot))

	blueBullet1 := &wobj.Object{
		Tag:      tagBlueBullet,
		Pos:      pos1,
		Size:     pixel.V(8, 8),
		Drawable: s.shot,
		Velocity: pixel.V(bulletSpeed, 0).Rotated(s.bluePlayer.Rot),
		Steps: wobj.MakeBehaviors(
			wobj.Movement,
		),
		PostSteps: wobj.MakeBehaviors(
			s.behaviorRemoveOutOfBounds,
		),
	}
	blueBullet2 := &wobj.Object{
		Tag:      tagBlueBullet,
		Pos:      pos2,
		Size:     pixel.V(8, 8),
		Drawable: s.shot,
		Velocity: pixel.V(bulletSpeed, 0).Rotated(s.bluePlayer.Rot),
		Steps: wobj.MakeBehaviors(
			wobj.Movement,
		),
		PostSteps: wobj.MakeBehaviors(
			s.behaviorRemoveOutOfBounds,
		),
	}
	s.layers[layerBullets].Add(blueBullet1)
	s.layers[layerBullets].Add(blueBullet2)

	s.w.speaker.Play(s.cannon)
}

func (s *gameScene) spawnRedShots() {

	bounds := s.redPlayer.Bounds()
	offset := pixel.V(bounds.H()/2, -8).Rotated(s.redPlayer.Rot)
	pos := bounds.Center().Add(offset)

	redBullet := &wobj.Object{
		Tag:      tagRedBullet,
		Pos:      pos,
		Size:     pixel.V(14, 14),
		Drawable: s.shot,
		Velocity: pixel.V(bulletSpeed, 0).Rotated(s.redPlayer.Rot),
		Steps: wobj.MakeBehaviors(
			wobj.Movement,
		),
		PostSteps: wobj.MakeBehaviors(
			s.behaviorRemoveOutOfBounds,
		),
	}
	s.layers[layerBullets].Add(redBullet)

	s.w.speaker.Play(s.cannon)
}

func (s *gameScene) behaviorRemoveOutOfBounds(source *wobj.Object, dt float64) {
	if !wo.Collision(source.Bounds(), s.bounds) {
		s.layers[layerBullets].Remove(source)
	}
}

func (s *gameScene) Bounds() pixel.Rect {
	return s.bounds
}

func (s *gameScene) behaviorRedHitsBlueBullet(source *wobj.Object, dt float64) {
	if s.phase != phaseBattle {
		return
	}
	sourceBounds := source.Bounds()
	for bullet := range s.layers[layerBullets].Tagged(tagBlueBullet) {
		if wo.Collision(sourceBounds, bullet.Bounds()) {
			s.w.blueScore++
			s.phase = phaseBlueVictory
			s.onVictory("Blue", colornames.Cadetblue)
		}
	}
}

func (s *gameScene) behaviorBlueHitsRedBullet(source *wobj.Object, dt float64) {
	if s.phase != phaseBattle {
		return
	}
	sourceBounds := source.Bounds()
	for bullet := range s.layers[layerBullets].Tagged(tagRedBullet) {
		if wo.Collision(sourceBounds, bullet.Bounds()) {
			s.w.redScore++
			s.phase = phaseRedVictory
			s.onVictory("Red", colornames.Indianred)
		}
	}
}

func (s *gameScene) onVictory(winner string, textColor color.Color) {
	s.victoryTime = victoryMessageDuration
	s.message.Clear()
	s.message.Color = textColor

	saying := winningMessages[s.rng.Intn(len(winningMessages))]
	victoryMessage := fmt.Sprintf(saying, winner)

	s.message.WriteString(victoryMessage)
}
