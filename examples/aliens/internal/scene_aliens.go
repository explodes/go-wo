package internal

import (
	"math/rand"
	"time"

	"golang.org/x/image/colornames"

	"github.com/explodes/go-wo"
	"github.com/explodes/go-wo/wobj"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
)

var _ wo.Scene = &scene{}

type scene struct {
	time   float64
	bounds pixel.Rect

	aliens      *AlienRow
	alienSprite *wobj.SpriteSheetDrawable

	ship *Ship
}

func (w *World) newSpaceScene(canvas *pixelgl.Canvas) (wo.Scene, error) {
	ship, err := newShip(w.loader, canvas.Bounds())
	if err != nil {
		return nil, err
	}

	whiteKey := wo.AlphaKeyTransformer(colornames.White)
	alienSprite, err := w.loader.SpriteSheet("img/ships.png", alienSheetOptions, whiteKey)
	if err != nil {
		return nil, err
	}
	alienSpriteDrawable := &wobj.SpriteSheetDrawable{alienSprite}

	const numAliens = 12
	aliens := &AlienRow{
		aliens:    make([]*Alien, 0, numAliens),
		direction: DirectionLeft,
		speed:     baseAlienSpeed,
	}
	kind := randomAlienKind()
	for i := 0; i < numAliens; i++ {
		alien := newAlien(alienSpriteDrawable, pixel.V(float64(60*i)+10, canvas.Bounds().H()-alienSize), kind)
		aliens.aliens = append(aliens.aliens, alien)
	}

	scene := &scene{
		bounds:      canvas.Bounds(),
		aliens:      aliens,
		alienSprite: alienSpriteDrawable,
		ship:        ship,
	}
	return scene, nil
}

func (s *scene) Update(dt float64, input wo.Input) wo.SceneResult {
	if input.Pressed(pixelgl.KeyEscape) {
		return wo.SceneResultError
	}

	s.time += dt
	s.ship.Update(dt, input)
	s.aliens.Update(dt, s.bounds)
	return wo.SceneResultNone
}

func (s *scene) Draw(canvas *pixelgl.Canvas) {
	s.ship.Draw(canvas)
	s.aliens.Draw(canvas)
}
