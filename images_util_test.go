package wo

import (
	"image"
	"image/color"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testImageStr = "test0123456789abcdef"
)

var (
	testImageBytes  = []byte(testImageStr)
	testImageColors = [][]int{
		{0xffff0000, 0xffff0000, 0xffff0000, 0xffff0000},
		{0xff00ff00, 0xff00ff00, 0xff00ff00, 0xff00ff00},
		{0xff0000ff, 0xff0000ff, 0xff0000ff, 0xff0000ff},
		{0xff000000, 0x00000000, 0xffffffff, 0x00ffffff},
	}
)

type testImage struct {
	pixels [][]int
}

func NewTestImage() *testImage {
	return NewTestImagePixels(testImageColors)
}

func NewTestImagePixels(pixels [][]int) *testImage {
	if len(pixels) != 4 {
		panic("test image must have 4 rows")
	}
	for _, row := range pixels {
		if len(row) != 4 {
			panic("test image row must have 4 columns")
		}
	}
	return &testImage{
		pixels: pixels,
	}
}

func (i *testImage) ColorModel() color.Model { return color.RGBAModel }
func (i *testImage) Bounds() image.Rectangle { return image.Rect(0, 0, 4, 4) }
func (i *testImage) At(x, y int) color.Color {
	px := i.pixels[x][y]
	return intAsColor(px)
}

func decodeTestImage(r io.Reader) (image.Image, error) {
	return NewTestImage(), nil
}

func decodeTestImageConfig(r io.Reader) (image.Config, error) {
	config := image.Config{
		Width:      4,
		Height:     4,
		ColorModel: color.RGBAModel,
	}
	return config, nil
}

func init() {
	image.RegisterFormat("test", "test", decodeTestImage, decodeTestImageConfig)
}

func intAsColor(i int) color.Color {
	return color.RGBA{
		A: uint8((i & 0xff000000) >> 24),
		R: uint8((i & 0x00ff0000) >> 16),
		G: uint8((i & 0x0000ff00) >> 8),
		B: uint8((i & 0x000000ff) >> 0),
	}
}

func colorAsInt(c color.Color) int {
	r, g, b, a := c.RGBA()
	return int((a&0xff)<<24) | int((r&0xff)<<16) | int((g&0xff)<<8) | int(b&0xff)
}

func assertImageBoundsEqual(t *testing.T, img1, img2 image.Image) {
	t.Helper()

	b1 := img1.Bounds()
	b2 := img2.Bounds()

	assert.Equal(t, b1, b2, "image bounds are unequal: %v != %v", b1, b2)
}

func assertImageEquals(t *testing.T, img1, img2 image.Image) {
	t.Helper()

	assertImageBoundsEqual(t, img1, img2)

	b1 := img1.Bounds()

	for x := b1.Min.X; x < b1.Max.X; x++ {
		for y := b1.Min.Y; y < b1.Max.Y; y++ {
			c1 := colorAsInt(img1.At(x, y))
			c2 := colorAsInt(img2.At(x, y))
			assert.Equal(t, c1, c2, "colors at (%d, %d) are unequal: 0x%08x != 0x%08x", x, y, c1, c2)
		}
	}
}
