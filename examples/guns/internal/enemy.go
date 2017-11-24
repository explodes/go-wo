package internal

import (
	"github.com/explodes/go-wo"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type enemyType int

const (
	enemyTypeHelo enemyType = iota
	enemyTypeFighter
	enemyTypeRamStation
	enemyTypeMagnum
)

const (
	numEnemyTypes  = 4
	numEnemyFrames = 4

	enemySpeed            = 25
	enemyHitBoxDifficulty = 0.6 // lower is harder to hit
)

type enemy struct {
	frames [4]pixel.Rect
	pos    pixel.Vec
	rot    float64
	sheet  *wo.SpriteSheet
	health float64
	frame  pixel.Rect

	debug bool
}

func newEnemy(class enemyType, sheet *wo.SpriteSheet, pos pixel.Vec, debug bool) *enemy {
	e := &enemy{
		pos:    pos,
		sheet:  sheet,
		health: 1,
		debug:  debug,
	}

	allFrames := sheet.Frames()

	for row := 0; row < enemySheetOptions.Rows; row++ {
		sheetIndex := int(class) + row*enemySheetOptions.Columns
		e.frames[numEnemyFrames-1-row] = allFrames[sheetIndex]
	}
	return e
}

func (e *enemy) update(dt float64, ship *ship) {
	if e.health < 0 {
		e.health = 0
	} else if e.health > 1 {
		e.health = 1
	}
	frameNum := int((numEnemyFrames - 1) * e.health)
	e.frame = e.frames[frameNum]

	dv := e.pos.Sub(ship.pos).Unit()

	e.rot = dv.Angle() + wo.DegToRad(90)
	e.pos = e.pos.Sub(dv.Scaled(enemySpeed * dt))
}

func (e *enemy) draw(canvas *pixelgl.Canvas) {
	sprite := e.sheet.SetFrame(0)
	sprite.Set(sprite.Picture(), e.frame)
	mat := pixel.IM.Moved(e.pos).Rotated(e.pos, e.rot)
	sprite.Draw(canvas, mat)

	if e.debug {
		im := imdraw.New(nil)

		im.Color = colornames.Orange
		im.Push(e.artBox().Min, e.artBox().Max)
		im.Rectangle(1)

		im.Color = colornames.Red
		im.Push(e.hitBox().Min, e.hitBox().Max)
		im.Rectangle(1)

		im.Draw(canvas)
	}
}

func (e *enemy) artBox() pixel.Rect {
	sprite := e.sheet.SetFrame(0)
	sprite.Set(sprite.Picture(), e.frame)
	bounds := pixel.R(
		e.pos.X,
		e.pos.Y,
		e.pos.X+sprite.Frame().W(),
		e.pos.Y+sprite.Frame().H(),
	)
	offset := pixel.V(-sprite.Frame().W()/2, -sprite.Frame().H()/2)
	return wo.TranslateRect(bounds, offset)
}

func (e *enemy) hitBox() pixel.Rect {
	return wo.ScaleRect(e.artBox(), enemyHitBoxDifficulty)
}
