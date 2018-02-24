package internal

import (
	"math"

	"image/color"

	"github.com/explodes/go-wo"
	"github.com/explodes/go-wo/wobj"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/usedbytes/hsv"
	"golang.org/x/image/colornames"
)

const (
	layerBackground = iota
	layerPlatforms
	layerPlayer
	numLayers
)

const (
	tagPlayer   = "player"
	tagPlatform = "platform"
	tagFinish   = "finish"
)

type mainScene struct {
	w *World

	bounds pixel.Rect
	time   float64

	level int

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
	s.generateNextLevel()
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
		PostSteps: wobj.MakeBehaviors(
			func(source *wobj.Object, dt float64) {
				iter := s.layers.TagIterator(tagFinish)
				for finish, ok := iter(); ok; finish, ok = iter() {
					if wo.Collision(source.Bounds(), finish.Bounds()) {
						s.level++
						s.generateNextLevel()
					}
				}
			},
		),
	}
}

func (s *mainScene) createPlatform(x, y, size float64, color color.Color, tag string) *wobj.Object {
	im := imdraw.New(nil)
	im.Color = color
	im.Push(pixel.V(0, 0), pixel.V(size, size))
	im.Rectangle(1)
	drawable := &IMDrawable{
		bounds: pixel.R(0, 0, size, size),
		im:     im,
	}
	return &wobj.Object{
		Tag:      tag,
		Pos:      pixel.V(x, y),
		Size:     pixel.V(size, size),
		Drawable: drawable,
	}
}

func (s *mainScene) randomColor() color.Color {
	return hsv.HSVColor{
		S: 255,
		V: 255,
		H: uint16(s.w.rng.Intn(360)),
	}
}

func (s *mainScene) randomColors() [7]color.Color {
	var c [7]color.Color
	for i := 0; i < len(c); i++ {
		c[i] = s.randomColor()
	}
	return c
}

func (s *mainScene) generateNextLevel() {

	platformColors := s.randomColors()

	layers := wobj.NewLayers(numLayers)
	s.layers = layers

	levelNum := s.level % len(levels)
	level := levels[levelNum]

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
			case 9:
				layers[layerPlatforms].Add(s.createPlatform(x, y, size, colornames.White, tagFinish))
			default:
				layers[layerBackground].Add(s.createPlatform(x, y, size, platformColors[piece-2], tagPlatform))
			}

		}
	}
}

var (
	levels = [][]int{
		//{
		//	2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
		//	2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
		//	2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
		//	2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
		//	2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
		//	2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
		//	2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
		//	2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
		//	2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
		//	2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
		//	2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
		//	2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
		//	2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
		//	2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
		//	2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
		//},
		{
			2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 9, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 3, 4, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 3, 4, 0, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 3, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 1, 3, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
		},
		{
			2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 6, 0, 4, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 0, 5, 0, 5, 0, 3, 0, 0, 2,
			2, 0, 0, 0, 0, 6, 0, 4, 0, 4, 0, 0, 0, 0, 2,
			2, 0, 5, 0, 0, 5, 0, 3, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 4, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 1, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
		},
		{
			2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 2,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 2,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 3, 3, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 3, 4, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
			2, 1, 3, 0, 0, 0, 4, 0, 4, 0, 0, 0, 0, 0, 2,
			2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
		},
		{
			2, 2, 2, 2, 2, 2, 2,
			2, 9, 0, 0, 0, 1, 2,
			2, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 0, 0, 2,
			2, 0, 0, 0, 3, 0, 2,
			2, 0, 0, 0, 3, 0, 2,
			2, 2, 2, 2, 2, 2, 2,
		},
	}
)
