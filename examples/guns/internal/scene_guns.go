package internal

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/explodes/go-wo"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type particleSetKey uint8

const (
	particlesBullets particleSetKey = iota
	particlesExplosions
)

const (
	explosionFps = 60
)

var (
	explosionSheetOptions = wo.SpriteSheetOptions{
		Rows:       9,
		Columns:    9,
		Width:      900 / 9,
		Height:     900 / 9,
		ExactCount: 74,
	}
	enemySheetOptions = wo.SpriteSheetOptions{
		Height:  512 / 4,
		Width:   512 / 4,
		Rows:    4,
		Columns: 4,
	}
)

type spaceScene struct {
	time   float64
	bounds pixel.Rect

	space          *imdraw.IMDraw
	explosionSheet *wo.SpriteSheet
	enemySheet     *wo.SpriteSheet

	speaker        *wo.Speaker
	shotSound      *wo.Sound
	explosionSound *wo.Sound

	ship *ship

	particleSets map[particleSetKey]*ParticleSet
	enemies      map[*enemy]struct{}

	rng *rand.Rand

	debug bool
}

var _ wo.Scene = &spaceScene{}

func (w *World) newSpaceScene(canvas *pixelgl.Canvas) (wo.Scene, error) {

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

	particles := map[particleSetKey]*ParticleSet{
		particlesBullets:    NewParticleCollection(256),
		particlesExplosions: NewParticleCollection(64),
	}

	ship, err := newShip(w.loader, particles[particlesBullets], speaker.Audible(shotSound), w.debug)
	if err != nil {
		return nil, err
	}
	ship.pos = canvas.Bounds().Center()

	explosionSheet, err := w.loader.SpriteSheet("img/explosion.png", explosionSheetOptions)
	if err != nil {
		return nil, err
	}

	enemySheet, err := w.loader.SpriteSheet("img/enemies.png", enemySheetOptions, wo.AlphaKeyTransformer(color.White))
	if err != nil {
		return nil, err
	}

	scene := &spaceScene{
		bounds:         canvas.Bounds(),
		space:          createSpace(canvas),
		explosionSheet: explosionSheet,
		enemySheet:     enemySheet,
		speaker:        speaker,
		shotSound:      shotSound,
		explosionSound: explosionSound,
		ship:           ship,
		particleSets:   particles,
		enemies:        make(map[*enemy]struct{}, 16),
		rng:            rand.New(rand.NewSource(time.Now().UnixNano())),
		debug:          w.debug,
	}
	return scene, nil
}

func (s *spaceScene) Update(dt float64, input wo.Input) wo.SceneResult {
	s.time += dt
	s.ship.update(dt, input)
	s.constrainShip()
	s.particleSets[particlesBullets].Update(dt)
	s.particleSets[particlesExplosions].Update(dt)
	s.cleanBullets()
	s.createEnemies()
	s.collisionDetection()
	s.updateEnemies(dt)

	if input.JustPressed(pixelgl.KeyF) {
		s.explode(s.bounds.Center(), pixel.V(500, 500))
	}

	return wo.SceneResultNone
}
func (s *spaceScene) collisionDetection() {
	// bullets -> enemies
	for bullet := range s.particleSets[particlesBullets].particles {
		box := bullet.HitBox()
		for enemy := range s.enemies {
			enemyBox := enemy.hitBox()
			if !wo.Collision(box, enemyBox) {
				continue
			}
			s.particleSets[particlesBullets].Remove(bullet)
			enemy.health -= 0.25
			if enemy.health < 0 {
				delete(s.enemies, enemy)
				s.explode(enemy.pos, enemy.artBox().Size())
			}
		}
	}
}

func (s *spaceScene) explode(pos, size pixel.Vec) {
	set := s.particleSets[particlesExplosions]
	ex := newAnimatedParticle(set, s.explosionSheet, explosionFps, pos, size, true)
	set.Add(ex)
	s.speaker.Play(s.explosionSound)
}

func (s *spaceScene) createEnemies() {
	if len(s.enemies) == 0 {
		class := s.rng.Intn(numEnemyTypes)
		pos := pixel.V(s.rng.Float64()*s.bounds.W(), s.rng.Float64()*s.bounds.H())
		enemy := newEnemy(enemyType(class), s.enemySheet, pos, s.debug)
		s.enemies[enemy] = struct{}{}
	}
}

func (s *spaceScene) updateEnemies(dt float64) {
	for enemy := range s.enemies {
		enemy.update(dt, s.ship)
	}
}

func (s *spaceScene) constrainShip() {
	shipBox := s.ship.hitBox()
	if shipBox.Min.X < s.bounds.Min.X {
		dx := s.bounds.Min.X - shipBox.Min.X
		s.ship.pos.X += dx
	}
	if shipBox.Max.X > s.bounds.Max.X {
		dx := s.bounds.Max.X - shipBox.Max.X
		s.ship.pos.X += dx
	}
	if shipBox.Min.Y < s.bounds.Min.Y {
		dy := s.bounds.Min.Y - shipBox.Min.Y
		s.ship.pos.Y += dy
	}
	if shipBox.Max.Y > s.bounds.Max.Y {
		dy := s.bounds.Max.Y - shipBox.Max.Y
		s.ship.pos.Y += dy
	}
}

func (s *spaceScene) cleanBullets() {
	bullets := s.particleSets[particlesBullets]
	for bullet := range bullets.particles {
		if !wo.Collision(bullet.HitBox(), s.bounds) {
			bullets.Remove(bullet)
		}
	}
}

func (s *spaceScene) Draw(canvas *pixelgl.Canvas) {
	s.space.Draw(canvas)
	s.ship.draw(canvas)
	s.particleSets[particlesBullets].Draw(canvas)
	s.particleSets[particlesExplosions].Draw(canvas)
	for enemy := range s.enemies {
		enemy.draw(canvas)
	}
}

func createSpace(canvas *pixelgl.Canvas) *imdraw.IMDraw {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	im := imdraw.New(nil)
	for stars := 1000; stars > 0; stars-- {
		pos := pixel.V(
			rng.Float64()*canvas.Bounds().W(),
			rng.Float64()*canvas.Bounds().H(),
		)
		im.Color = color.White
		im.Push(pos, pos.Add(pixel.V(1, 1)))
		im.Rectangle(0)
	}
	return im
}
