package wo

import (
	"github.com/faiface/pixel"
)

type SpriteSheet struct {
	sprite  *pixel.Sprite
	pic     *pixel.PictureData
	frames  []pixel.Rect
	options SpriteSheetOptions
}

type SpriteSheetOptions struct {
	Width      int
	Height     int
	Columns    int
	Rows       int
	ExactCount int
}

func NewSpriteSheet(pic *pixel.PictureData, opts SpriteSheetOptions) *SpriteSheet {
	frames := makeFrames(int(pic.Bounds().H()), opts)
	return &SpriteSheet{
		pic:     pic,
		sprite:  pixel.NewSprite(pic, frames[0]),
		frames:  frames,
		options: opts,
	}
}

func makeFrames(imageHeight int, opts SpriteSheetOptions) []pixel.Rect {
	var N int
	if opts.ExactCount != 0 {
		N = opts.ExactCount
	} else {
		N = opts.Rows * opts.Columns
	}
	frames := make([]pixel.Rect, 0, N)
	for row := 0; row < opts.Rows; row++ {
		y := imageHeight - row*opts.Height
		for col := 0; col < opts.Columns; col++ {
			x := col * opts.Width
			frame := pixel.R(float64(x), float64(y-opts.Height), float64(x+opts.Width), float64(y))
			frames = append(frames, frame)
			if len(frames) == N {
				return frames
			}
		}
	}
	return frames
}

func (ss *SpriteSheet) Options() SpriteSheetOptions {
	return ss.options
}

func (ss *SpriteSheet) Frame() *pixel.Sprite {
	return ss.sprite
}

func (ss *SpriteSheet) Frames() []pixel.Rect {
	return ss.frames
}

func (ss *SpriteSheet) SetFrame(frame int) *pixel.Sprite {
	ss.sprite.Set(ss.pic, ss.frames[frame])
	return ss.sprite
}

func (ss *SpriteSheet) NumFrames() int {
	return len(ss.frames)
}
