package wo

import (
	"github.com/faiface/pixel"
)

// SpriteSheet is a Sprite with a set of Rectangles that
// defines sections of the Sprite to render.

// The use of the SetFrame and Sprite methods
// should not be used concurrently.
type SpriteSheet struct {
	sprite  *pixel.Sprite
	pic     *pixel.PictureData
	frames  []pixel.Rect
	options SpriteSheetOptions
}

// SpriteSheetOptions specifies options for loading a SpriteSheet.
type SpriteSheetOptions struct {
	Width      int
	Height     int
	Columns    int
	Rows       int
	ExactCount int
}

// NewSpriteSheet creates a SpriteSheet from
// PictureData and the given options.
func NewSpriteSheet(pic *pixel.PictureData, opts SpriteSheetOptions) *SpriteSheet {
	frames := makeFrames(int(pic.Bounds().H()), opts)
	return &SpriteSheet{
		pic:     pic,
		sprite:  pixel.NewSprite(pic, frames[0]),
		frames:  frames,
		options: opts,
	}
}

// makeFrames is the function to create the
// rectangles for portions of the Sprite.
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

// Options returns the SpriteSheetOptions used to load this SpriteSheet.
func (ss *SpriteSheet) Options() SpriteSheetOptions {
	return ss.options
}

// Sprite returns the Sprite with the current frame.
func (ss *SpriteSheet) Sprite() *pixel.Sprite {
	return ss.sprite
}

// Frames returns the Frames defined in this SpriteSheet.
func (ss *SpriteSheet) Frames() []pixel.Rect {
	return ss.frames
}

// SetFrame prepares the Sprite to be rendered with a
// particular frame number.
func (ss *SpriteSheet) SetFrame(frameNum int) *pixel.Sprite {
	ss.sprite.Set(ss.pic, ss.frames[frameNum])
	return ss.sprite
}

// NumFrames returns the total number of frames
// available in this SpriteSheet.
func (ss *SpriteSheet) NumFrames() int {
	return len(ss.frames)
}

// Bounds returns the bounds of the current frame
func (ss *SpriteSheet) Bounds() pixel.Rect {
	return ss.sprite.Frame()
}
