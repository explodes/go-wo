package internal

import (
	"fmt"
	_ "image/png"
	"math/rand"
	"time"

	"github.com/explodes/go-wo"
	"github.com/explodes/go-wo/examples/flappy/res"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

const (
	debugPlayScene = true
	numPipes       = 3
)

var (
	scoreColor = colornames.White
)

type playScene struct {
	game       *FlappyWorld
	log        *logrus.Entry
	pipeLimits pixel.Rect
	rng        *rand.Rand

	// backgroundSprite is the backgroundSprite image of the scene
	backgroundSprite *pixel.Sprite

	// birdSprites are the frames for the bird animation
	birdSprites []*pixel.Sprite

	// pipeSprite is the graphic used to draw
	// the pipeSprite obstacles
	pipeSprite *pixel.Sprite

	scoreText *text.Text

	// bird is the player character
	bird *playBird

	// pipes are the pipeSprite obstacles
	pipes []*playPipe

	score int
	dead  bool
}

func decodeFontFace(ttf []byte, size float64) (font.Face, error) {
	f, err := truetype.Parse(ttf)
	if err != nil {
		return nil, err
	}
	return truetype.NewFace(f, &truetype.Options{Size: size}), nil
}

func (g *FlappyWorld) createPlayScene(canvas *pixelgl.Canvas) (wo.Scene, error) {
	fontSrc, err := res.Load("fonts/Flappy.ttf")
	if err != nil {
		return nil, errors.Errorf("could not find font: %v", err)
	}
	small, err := decodeFontFace(fontSrc, 32)
	defer small.Close()
	scoreText := text.New(pixel.V(0, 0), text.NewAtlas(small, text.ASCII))
	scoreText.Color = scoreColor

	pic, err := res.NewPic("imgs/background.png")
	if err != nil {
		return nil, errors.Errorf("could not load background image: %v", err)
	}
	backgroundSprite := pixel.NewSprite(pic, pic.Bounds())

	birdSprites := make([]*pixel.Sprite, 4, 4)
	for i := 1; i <= 4; i++ {
		pic, err := res.NewPic(fmt.Sprintf("imgs/bird_frame_%d.png", i))
		if err != nil {
			return nil, errors.Errorf("could not load bird image %d: %v", i, err)
		}
		sprite := pixel.NewSprite(pic, pic.Bounds())
		birdSprites[i-1] = sprite
	}

	pic, err = res.NewPic("imgs/pipe.png")
	if err != nil {
		return nil, errors.Errorf("could not load pipe image: %v", err)
	}
	pipeSprite := pixel.NewSprite(pic, pic.Bounds())
	pipeBounds := pipeSprite.Picture().Bounds()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	w := canvas.Bounds().W()
	pipeLimits := pixel.R(0, -pipeBounds.H(), w, canvas.Bounds().H()-pipeBounds.H())
	pipes := make([]*playPipe, numPipes, numPipes)
	for index := range pipes {
		pipe := &playPipe{}
		pipe.reset(rng, pipeLimits)
		pipe.pos.X = w + (float64(index)/float64(len(pipes)))*w + pipeWidth
		pipes[index] = pipe
	}

	scene := &playScene{
		game:             g,
		log:              g.log.WithField("scene", "FlappyWorld"),
		rng:              rng,
		pipeLimits:       pipeLimits,
		backgroundSprite: backgroundSprite,
		birdSprites:      birdSprites,
		pipeSprite:       pipeSprite,
		scoreText:        scoreText,
		bird: &playBird{
			physics: physics{
				pos: pixel.V(10, canvas.Bounds().H()/2),
			},
		},
		pipes: pipes,
	}
	return scene, nil
}

func (s *playScene) Update(dt float64, input wo.Input) wo.SceneResult {
	s.bird.update(dt, input, s.dead)

	if s.bird.pos.Y <= 0 {
		s.dead = true
	}

	if s.dead {
		if s.bird.pos.Y < 0 {
			s.game.lastScore = s.score
			return SceneResultGoToTitle
		}
		return wo.SceneResultNone
	}

	for _, pipe := range s.pipes {
		pipe.update(dt, basePipeSpeed, s.rng, s.pipeLimits)
	}

	s.detectPipes()

	s.scoreText.Clear()
	s.scoreText.WriteString(fmt.Sprintf("\r%dpts", s.score))

	return wo.SceneResultNone
}

func (s *playScene) detectPipes() {
	birdBox := wo.HitBox(s.birdSprites[0], s.bird.pos)
	for _, pipe := range s.pipes {
		pipeBox := wo.HitBox(s.pipeSprite, pipe.pos)
		if wo.Collision(birdBox, pipeBox) {
			s.dead = true
			return
		}
		inverted := pipe.invertedHitBox(s.pipeSprite)
		if wo.Collision(birdBox, inverted) {
			s.dead = true
			return
		}
		if !pipe.scored && birdBox.Min.X >= pipeBox.Max.X {
			s.score++
			pipe.scored = true
		}
	}
}

func (s *playScene) Draw(canvas *pixelgl.Canvas) {
	s.drawBackground(canvas)
	s.drawPipes(canvas)
	s.drawScore(canvas)
	s.drawBird(canvas)
}

func (s *playScene) drawBackground(canvas *pixelgl.Canvas) {
	mat := wo.FitAtZero(s.backgroundSprite.Frame(), canvas.Bounds())
	s.backgroundSprite.Draw(canvas, mat)
}

func (s *playScene) drawPipes(canvas *pixelgl.Canvas) {
	for _, pipe := range s.pipes {
		pipe.draw(canvas, s.pipeSprite)
		s.drawDebug(canvas, s.pipeSprite, pipe.pos)
		s.drawDebugRect(canvas, pipe.invertedHitBox(s.pipeSprite))
	}
}

func (s *playScene) drawScore(canvas *pixelgl.Canvas) {
	s.scoreText.Draw(canvas, pixel.IM.Moved(pixel.V(10, canvas.Bounds().Max.Y-s.scoreText.LineHeight)))
}

func (s *playScene) drawBird(canvas *pixelgl.Canvas) {
	s.bird.draw(canvas, s.birdSprites)
	s.drawDebug(canvas, s.birdSprites[0], s.bird.pos)
}

func (s *playScene) drawDebug(canvas *pixelgl.Canvas, sprite *pixel.Sprite, pos pixel.Vec) {
	if !debugPlayScene {
		return
	}
	hitbox := wo.HitBox(sprite, pos)
	s.drawDebugRect(canvas, hitbox)
}

func (s *playScene) drawDebugRect(canvas *pixelgl.Canvas, r pixel.Rect) {
	if !debugPlayScene {
		return
	}
	im := imdraw.New(nil)
	im.Color = colornames.Red
	im.Push(r.Min, r.Max)
	im.Rectangle(2)
	im.Draw(canvas)
}
