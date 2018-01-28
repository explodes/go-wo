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

const (
	bulletSpeed = 240
)

type scene struct {
	time   float64
	bounds pixel.Rect

	aliens      *AlienRow
	alienSprite *wobj.SpriteSheetDrawable

	bullets      *BulletPool
	shotSfx      *wo.Audible
	explosionSfx *wo.Audible

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
	alienSpriteDrawable := &wobj.SpriteSheetDrawable{Sheet: alienSprite}

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

	speaker, err := wo.NewSpeaker()
	if err != nil {
		return nil, err
	}
	shotSound, err := w.loader.Sound("wav", "sound/shot.wav")
	if err != nil {
		return nil, err
	}
	explosionSound, err := w.loader.Sound("wav", "sound/boom.wav")
	if err != nil {
		return nil, err
	}

	scene := &scene{
		bounds:       canvas.Bounds(),
		aliens:       aliens,
		alienSprite:  alienSpriteDrawable,
		bullets:      NewBulletPool(),
		ship:         ship,
		shotSfx:      speaker.Audible(shotSound),
		explosionSfx: speaker.Audible(explosionSound),
	}
	return scene, nil
}

func (s *scene) Update(dt float64, input wo.Input) wo.SceneResult {
	s.time += dt

	// process input
	if input.Pressed(pixelgl.KeyEscape) {
		return wo.SceneResultError
	}
	if input.JustPressed(pixelgl.KeySpace) || input.JustPressed(pixelgl.KeyW) {
		bullet := s.bullets.Spawn()
		bullet.Pos = s.ship.Pos.Add(pixel.V(s.ship.Size.X*0.25, s.ship.Size.Y))
		s.shotSfx.Play()
	}

	// update elements
	s.ship.Update(dt, input)
	s.aliens.Update(dt, s.bounds)
	for bullet := range s.bullets.Active {
		bullet.Pos.Y += bulletSpeed * dt
	}

	// test collisions
	s.collideBullets()

	return wo.SceneResultNone
}

func (s *scene) collideBullets() {
	for bullet := range s.bullets.Active {
		bulletBounds := bullet.Bounds()
		for alienIndex, alien := range s.aliens.aliens {
			if alien != nil && wo.Collision(bulletBounds, alien.Bounds()) {
				s.aliens.aliens[alienIndex] = nil
				s.explosionSfx.Play()
			}
		}
	}
}

func (s *scene) Draw(canvas *pixelgl.Canvas) {
	for bullet := range s.bullets.Active {
		bullet.Draw(canvas)
	}
	s.ship.Draw(canvas)
	s.aliens.Draw(canvas)

}
