package internal

import (
	"math/rand"
	"time"

	"golang.org/x/image/colornames"

	"image/color"

	"github.com/explodes/go-wo"
	"github.com/explodes/go-wo/wobj"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

var (
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
)

var _ wo.Scene = &scene{}

const (
	bulletSpeed    = 360
	newRowDuration = 3 * time.Second
	bulletDelay    = 0.35
)

type scene struct {
	time   float64
	bounds pixel.Rect

	aliens      []*AlienRow
	alienSprite *wobj.SpriteSheetDrawable

	bullets      *BulletPool
	bulletDelay  float64
	shotSfx      *wo.Audible
	explosionSfx *wo.Audible

	newRows *time.Ticker

	ship *Ship

	bg *imdraw.IMDraw
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
		alienSprite:  alienSpriteDrawable,
		bullets:      NewBulletPool(),
		ship:         ship,
		shotSfx:      speaker.Audible(shotSound),
		explosionSfx: speaker.Audible(explosionSound),
		newRows:      time.NewTicker(newRowDuration),
		bg:           newBackground(canvas),
	}
	scene.addNewAlienRow()

	return scene, nil
}

func newBackground(canvas *pixelgl.Canvas) *imdraw.IMDraw {
	const numStars = 700
	im := imdraw.New(nil)
	for i := 0; i < numStars; i++ {
		x := rng.Float64()*canvas.Bounds().W() + canvas.Bounds().Min.X
		y := rng.Float64()*canvas.Bounds().H() + canvas.Bounds().Min.Y
		im.Push(pixel.V(x, y))
		im.Color = color.Gray{Y: uint8(rng.Int31n(128))}
		im.Circle(1, 0)
	}
	return im
}

func (s *scene) addNewAlienRow() {
	s.aliens = append(s.aliens, s.newAlienRow())
}

func (s *scene) newAlienRow() *AlienRow {
	const numAliens = 12
	aliens := &AlienRow{
		aliens:    make([]*Alien, 0, numAliens),
		direction: DirectionLeft,
		speed:     baseAlienSpeed,
	}
	kind := randomAlienKind()
	for i := 0; i < numAliens; i++ {
		alien := newAlien(s.alienSprite, pixel.V(float64(60*i)+10, s.bounds.H()-alienSize), kind)
		aliens.aliens = append(aliens.aliens, alien)
	}
	return aliens
}

func (s *scene) Update(dt float64, input wo.Input) wo.SceneResult {
	s.time += dt
	s.bulletDelay -= dt

	// process input
	if input.Pressed(pixelgl.KeyEscape) {
		return wo.SceneResultError
	}
	if input.JustPressed(pixelgl.KeySpace) || input.JustPressed(pixelgl.KeyW) || (s.bulletDelay <= 0 && (input.Pressed(pixelgl.KeySpace) || input.Pressed(pixelgl.KeyW))) {
		bullet := s.bullets.Spawn()
		bullet.Pos = s.ship.Pos.Add(pixel.V(s.ship.Size.X*0.25, s.ship.Size.Y))
		s.shotSfx.Play()
		s.bulletDelay = bulletDelay
	}

	// create elements
	if len(s.aliens) == 0 {
		s.addNewAlienRow()
	} else {
		select {
		case <-s.newRows.C:
			s.addNewAlienRow()
		default:
		}
	}

	// update elements
	s.ship.Update(dt, input)
	for _, row := range s.aliens {
		row.Update(dt, s.bounds)
	}
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
		for _, row := range s.aliens {
			for alienIndex, alien := range row.aliens {
				if alien != nil && wo.Collision(bulletBounds, alien.Bounds()) {
					row.aliens[alienIndex] = nil
					s.explosionSfx.Play()
				}
			}
		}
	}
}

func (s *scene) Draw(canvas *pixelgl.Canvas) {
	s.bg.Draw(canvas)
	for bullet := range s.bullets.Active {
		bullet.Draw(canvas)
	}
	s.ship.Draw(canvas)
	for _, row := range s.aliens {
		row.Draw(canvas)
	}
}
