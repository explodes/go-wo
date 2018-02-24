package internal

import (
	"math"

	"image/color"

	"github.com/explodes/go-wo"
	"github.com/explodes/go-wo/wobj"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	layerBackground = iota
	layerPlatforms
	layerPlayer
	layerForeground
	numLayers
)

const (
	tagPlayer   = "player"
	tagPlatform = "platform"
)

type mainScene struct {
	w *World

	bounds pixel.Rect
	time   float64

	layers wobj.Layers
}

func (w *World) createMainScene(canvas *pixelgl.Canvas) (wo.Scene, error) {
	if width != height {
		panic("window is not square!")
	}
	s := &mainScene{
		w:      w,
		bounds: canvas.Bounds(),
		layers: wobj.NewLayers(numLayers),
	}
	s.generateLevel()
	return s, nil
}

func (s *mainScene) Update(dt float64, input wo.Input) wo.SceneResult {
	s.time += dt

	if input.JustPressed(pixelgl.KeyEscape) {
		return wo.SceneResultWindowClosed
	}

	s.layers.Update(dt)

	return wo.SceneResultNone
}

func (s *mainScene) Draw(canvas *pixelgl.Canvas) {
	s.layers.Draw(canvas)
}

func (s *mainScene) createPlayer(x, y, size float64) *wobj.Object {
	const (
		gravity   = -512
		jumpSpeed = 340
		runSpeed  = 128
	)
	im := imdraw.New(nil)
	im.Color = colornames.Red
	im.Push(pixel.V(size/2, size/2))
	im.Circle(size/2, 1)
	drawable := &IMDrawable{
		bounds: pixel.R(0, 0, size, size),
		im:     im,
	}

	input := s.w.input

	return &wobj.Object{
		Tag:      tagPlayer,
		Pos:      pixel.V(x, y),
		Size:     pixel.V(size, size),
		Drawable: drawable,
		Steps: wobj.MakeBehaviors(
			func(player *wobj.Object, dt float64) {

				ground := false

				// apply gravity and velocity
				player.Velocity.Y += gravity * dt
				player.Pos = player.Pos.Add(player.Velocity.Scaled(dt))

				// check collisions against each platform
				bounds := player.Bounds()

				iter := s.layers.TagIterator(tagPlatform)
				for platform, ok := iter(); ok; platform, ok = iter() {
					platformBounds := platform.Bounds()
					if bounds.Max.X <= platformBounds.Min.X || bounds.Min.X >= platformBounds.Max.X {
						continue
					}
					if bounds.Min.Y > platformBounds.Max.Y || bounds.Min.Y < platformBounds.Max.Y+player.Velocity.Y*dt {
						continue
					}
					player.Velocity.Y = 0
					player.Pos = player.Pos.Add(pixel.V(0, platformBounds.Max.Y-bounds.Min.Y))
					ground = true
				}

				// can jump off of the ground
				if ground && input.Pressed(pixelgl.KeyW, pixelgl.KeyUp) {
					player.Velocity.Y = jumpSpeed
				}

				// can move but not into walls
				switch {
				case input.Pressed(pixelgl.KeyA, pixelgl.KeyLeft):
					player.Velocity.X = -runSpeed
				case input.Pressed(pixelgl.KeyD, pixelgl.KeyRight):
					player.Velocity.X = +runSpeed
				default:
					player.Velocity.X = 0
				}

			},
		),
	}
}

func (s *mainScene) createPlatform(x, y, size float64, color color.Color) *wobj.Object {
	im := imdraw.New(nil)
	im.Color = color
	im.Push(pixel.V(0, 0), pixel.V(size, size))
	im.Rectangle(1)
	drawable := &IMDrawable{
		bounds: pixel.R(0, 0, size, size),
		im:     im,
	}
	return &wobj.Object{
		Tag:      tagPlatform,
		Pos:      pixel.V(x, y),
		Size:     pixel.V(size, size),
		Drawable: drawable,
	}
}

func (s *mainScene) generateLevel() {
	layers := s.layers

	square := math.Sqrt(float64(len(level)))
	isquare := int(square)
	if float64(isquare) != square {
		panic("level is not square!")
	}
	size := width / square
	//isize := int(size)

	for i := 0; i < isquare; i++ {
		for j := 0; j < isquare; j++ {
			x := float64(i) * size
			y := height - float64(j)*size - size

			index := i + j*isquare
			piece := level[index]

			switch piece {
			case 0:
			case 1:
				layers[layerPlayer].Add(s.createPlayer(x, y, size))
			case 2:
				layers[layerBackground].Add(s.createPlatform(x, y, size, colornames.Blue))
			case 3:
				layers[layerPlatforms].Add(s.createPlatform(x, y, size, colornames.Green))
			case 4:
				layers[layerPlatforms].Add(s.createPlatform(x, y, size, colornames.Orange))
			}

		}
	}
}

var (
	level = []int{
		2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
		2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
		2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
		2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
		2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
		2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
		2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
		2, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 2,
		2, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 2,
		2, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 2,
		2, 0, 0, 0, 0, 3, 3, 0, 0, 0, 0, 0, 0, 0, 2,
		2, 0, 0, 0, 0, 3, 4, 0, 0, 0, 0, 0, 0, 0, 2,
		2, 0, 0, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
		2, 1, 3, 0, 0, 0, 4, 0, 4, 0, 0, 0, 0, 0, 2,
		2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
	}
)
