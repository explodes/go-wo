package wo

import (
	"testing"
	"github.com/cocoonlife/testify/assert"
	"github.com/faiface/pixel"
)

func TestMakeFrames_specificCount(t *testing.T) {
	as := assert.New(t)

	opts := SpriteSheetOptions{
		Width:      10,
		Height:     10,
		Rows:       3,
		Columns:    3,
		ExactCount: 8,
	}

	frames := makeFrames(30, opts)

	as.Equal(8, len(frames))
	as.Equal(pixel.R(0, 20, 10, 30), frames[0])
	as.Equal(pixel.R(10, 20, 20, 30), frames[1])
	as.Equal(pixel.R(20, 20, 30, 30), frames[2])
	as.Equal(pixel.R(0, 10, 10, 20), frames[3])
	as.Equal(pixel.R(10, 10, 20, 20), frames[4])
	as.Equal(pixel.R(20, 10, 30, 20), frames[5])
	as.Equal(pixel.R(0, 0, 10, 10), frames[6])
	as.Equal(pixel.R(10, 0, 20, 10), frames[7])
}

func TestMakeFrames_noSpecificCount(t *testing.T) {
	as := assert.New(t)

	opts := SpriteSheetOptions{
		Width:      10,
		Height:     10,
		Rows:       3,
		Columns:    3,
		ExactCount: 0,
	}

	frames := makeFrames(30, opts)

	as.Equal(9, len(frames))
	as.Equal(pixel.R(0, 20, 10, 30), frames[0])
	as.Equal(pixel.R(10, 20, 20, 30), frames[1])
	as.Equal(pixel.R(20, 20, 30, 30), frames[2])
	as.Equal(pixel.R(0, 10, 10, 20), frames[3])
	as.Equal(pixel.R(10, 10, 20, 20), frames[4])
	as.Equal(pixel.R(20, 10, 30, 20), frames[5])
	as.Equal(pixel.R(0, 0, 10, 10), frames[6])
	as.Equal(pixel.R(10, 0, 20, 10), frames[7])
	as.Equal(pixel.R(20, 0, 30, 10), frames[8])
}
