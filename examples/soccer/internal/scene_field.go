package internal

import (
	"github.com/explodes/go-wo"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	goalWidth  = 50
	goalHeight = 200
	goalOffset = height/2 - goalHeight/2
)

type hitBoxer interface {
	hitBox() pixel.Rect
}

type fieldScene struct {
	ball     *ball
	player   *player
	opponent *opponent

	playerGoal   pixel.Rect
	opponentGoal pixel.Rect

	playerScore, opponentScore int

	backgroundSprite *pixel.Sprite

	bounds pixel.Rect

	debug bool
}

func (w *World) createFieldScene(canvas *pixelgl.Canvas) (wo.Scene, error) {
	backgroundSprite, err := w.loader.Sprite("img/field.png", wo.ResizeTransformer(width, height))
	if err != nil {
		return nil, err
	}

	ballSprite, err := w.loader.Sprite("img/ball.png", wo.ResizeTransformer(ballWidth, ballHeight))
	if err != nil {
		return nil, err
	}
	ball, err := newBall(ballSprite)
	if err != nil {
		return nil, err
	}

	playerSprite, err := w.loader.Sprite("img/player.png", wo.ResizeTransformer(playerWidth, playerHeight))
	if err != nil {
		return nil, err
	}
	player, err := newPlayer(playerSprite)
	if err != nil {
		return nil, err
	}

	opponentSprite, err := w.loader.Sprite("img/opponent.png", wo.ResizeTransformer(opponentWidth, opponentHeight))
	if err != nil {
		return nil, err
	}
	opponent, err := newOpponent(opponentSprite)
	if err != nil {
		return nil, err
	}

	scene := &fieldScene{
		ball:             ball,
		player:           player,
		opponent:         opponent,
		playerGoal:       pixel.R(0, goalOffset, goalWidth, goalOffset+goalHeight),
		opponentGoal:     pixel.R(width-goalWidth, goalOffset, width, goalOffset+goalHeight),
		backgroundSprite: backgroundSprite,
		bounds:           canvas.Bounds(),
		debug:            w.debug,
	}

	scene.resetPositions()

	return scene, nil
}

func (s *fieldScene) resetPositions() {
	s.ball.pos = s.bounds.Center()
	s.player.pos = pixel.V(0, s.bounds.H()/2)
	s.opponent.pos = pixel.V(s.bounds.W(), s.bounds.H()/2)
}

func (s *fieldScene) Update(dt float64, input wo.Input) wo.SceneResult {
	s.ball.update(dt)
	s.player.update(dt, s.ball, input)
	s.opponent.update(dt, s.ball)

	s.hitBall(s.player)
	s.hitBall(s.opponent)

	s.constrainBall()

	if s.collideGoal(s.playerGoal) {
		s.opponentScore++
		s.resetPositions()
	}
	if s.collideGoal(s.opponentGoal) {
		s.playerScore++
		s.resetPositions()
	}

	return wo.SceneResultNone
}

func (s *fieldScene) hitBall(hitBoxer hitBoxer) {
	ballBox := s.ball.hitBox()
	box := hitBoxer.hitBox()
	if !wo.Collision(ballBox, box) {
		return
	}
	dir := ballBox.Center().Sub(box.Center()).Unit()
	s.ball.hit(dir)
}

func (s *fieldScene) constrainBall() {
	ballBox := s.ball.hitBox()
	if ballBox.Min.X < 0 {
		s.ball.vel = reflectVert(s.ball.vel)
		s.ball.pos = pixel.V(0, s.ball.pos.Y)
	}
	if ballBox.Max.X > s.bounds.Max.X {
		s.ball.vel = reflectVert(s.ball.vel)
		s.ball.pos = pixel.V(s.bounds.Max.X-ballBox.W(), s.ball.pos.Y)
	}
	if ballBox.Min.Y < 0 {
		s.ball.vel = reflectHori(s.ball.vel)
		s.ball.pos = pixel.V(s.ball.pos.X, 0)
	}
	if ballBox.Max.Y > s.bounds.Max.Y {
		s.ball.vel = reflectHori(s.ball.vel)
		s.ball.pos = pixel.V(s.ball.pos.X, s.bounds.Max.Y-ballBox.H())
	}
}

func (s *fieldScene) collideGoal(goal pixel.Rect) bool {
	ballBox := s.ball.hitBox()
	if !wo.Collision(ballBox, goal) {
		return false
	}
	return true
}

func reflectVert(vel pixel.Vec) pixel.Vec {
	return pixel.V(-vel.X, vel.Y)
}

func reflectHori(vel pixel.Vec) pixel.Vec {
	return pixel.V(vel.X, -vel.Y)
}

func (s *fieldScene) Draw(canvas *pixelgl.Canvas) {
	s.backgroundSprite.Draw(canvas, pixel.IM.Moved(canvas.Bounds().Center()))
	s.ball.draw(canvas)
	s.player.draw(canvas)
	s.opponent.draw(canvas)
	if s.debug {
		drawHitBoxes(canvas, s.ball, s.player, s.opponent, rectHitBox(s.playerGoal), rectHitBox(s.opponentGoal))
	}
}

type rectHitBox pixel.Rect

func (r rectHitBox) hitBox() pixel.Rect {
	return pixel.Rect(r)
}

func drawHitBoxes(canvas *pixelgl.Canvas, hitBoxers ...hitBoxer) {
	im := imdraw.New(nil)
	im.Color = colornames.Red
	for _, h := range hitBoxers {
		box := h.hitBox()
		im.Push(box.Min, box.Max)
		im.Rectangle(1)
	}
	im.Draw(canvas)
}
