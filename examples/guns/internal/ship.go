package internal

import (
	"github.com/explodes/go-wo"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"math"
)

const (
	shipRotationSpeed   = 270 * math.Pi / 180 // rotations per second
	shipDefaultRotation = 45 * math.Pi / 180  // rotate the image to point up initially
	shipImageScale      = 50. / 512.          // 50x50 render size
	shipHitBoxLeniency  = 0.75                // gotta get real close
	shipMoveSpeed       = 200

	shipBulletFireDelay = 1. / 6. // 3 bullets per second
	shipBulletSize      = 2.5
	shipBulletSpeed     = 1000
)

type ship struct {
	pos pixel.Vec

	rot    float64
	sprite *pixel.Sprite

	bulletTime float64
	bullets    *ParticleSet

	gunAudible *wo.Audible

	debug bool
}

func newShip(loader wo.Loader, bullets *ParticleSet, gunAudible *wo.Audible, debug bool) (*ship, error) {
	sprite, err := loader.Sprite("img/ship_512.png")
	if err != nil {
		return nil, err
	}
	ship := &ship{
		sprite:     sprite,
		bullets:    bullets,
		gunAudible: gunAudible,
		debug:      debug,
	}
	return ship, nil
}

func (s *ship) update(dt float64, input wo.Input) {

	s.bulletTime -= dt

	if input.Pressed(pixelgl.KeyA) || input.Pressed(pixelgl.KeyLeft) {
		s.rot += dt * shipRotationSpeed
	}
	if input.Pressed(pixelgl.KeyD) || input.Pressed(pixelgl.KeyRight) {
		s.rot -= dt * shipRotationSpeed
	}
	if input.Pressed(pixelgl.KeyW) || input.Pressed(pixelgl.KeyUp) {
		mov := pixel.V(0, shipMoveSpeed*dt).Rotated(s.rot)
		s.pos = s.pos.Add(mov)
	}
	if input.Pressed(pixelgl.KeySpace) && s.bulletTime <= 0 {
		offset := pixel.V(0, s.hitBox().H()/2).Rotated(s.rot)
		pos := s.pos.Add(offset)
		vel := pixel.V(0, shipBulletSpeed).Rotated(s.rot)
		s.bullets.Add(newBullet(pos, vel, shipBulletSize, colornames.Orange))
		s.bulletTime = shipBulletFireDelay
		s.gunAudible.Play()
	}
}

func (s *ship) draw(canvas *pixelgl.Canvas) {
	mat := pixel.IM.Moved(s.pos).Rotated(s.pos, s.rot+shipDefaultRotation).Scaled(s.pos, shipImageScale)
	s.sprite.Draw(canvas, mat)

	if s.debug {
		im := imdraw.New(nil)

		im.Color = colornames.Orange
		im.Push(s.artBox().Min, s.artBox().Max)
		im.Rectangle(1)

		im.Color = colornames.Red
		im.Push(s.hitBox().Min, s.hitBox().Max)
		im.Rectangle(1)

		im.Draw(canvas)
	}
}

func (s *ship) artBox() pixel.Rect {
	bounds := pixel.R(
		s.pos.X,
		s.pos.Y,
		s.pos.X+s.sprite.Frame().W()*shipImageScale,
		s.pos.Y+s.sprite.Frame().H()*shipImageScale,
	)
	offset := pixel.V(-s.sprite.Frame().W()/2*shipImageScale, -s.sprite.Frame().H()/2*shipImageScale)
	return wo.TranslateRect(bounds, offset)
}

func (s *ship) hitBox() pixel.Rect {
	return wo.ScaleRect(s.artBox(), shipHitBoxLeniency)
}
