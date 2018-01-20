package wo

import (
	"testing"

	"image"

	"errors"

	"github.com/stretchr/testify/assert"
)

func TestDecodePicture(t *testing.T) {
	p, err := DecodePicture(testImageBytes)

	assert.Nil(t, err)
	assert.NotNil(t, p)

	img := p.Image()
	assertImageEquals(t, NewTestImage(), img)
}

func TestAlphaKeyTransformer(t *testing.T) {
	// transform all solid red pixels to transparent pixels
	transformer := AlphaKeyTransformer(intAsColor(0xffff0000))
	base := NewTestImage()

	out, err := transformer(base)
	assert.Nil(t, err)

	expected := NewTestImagePixels([][]int{
		{0x00000000, 0x00000000, 0x00000000, 0x00000000},
		{0xff00ff00, 0xff00ff00, 0xff00ff00, 0xff00ff00},
		{0xff0000ff, 0xff0000ff, 0xff0000ff, 0xff0000ff},
		{0xff000000, 0x00000000, 0xffffffff, 0x00ffffff},
	})

	assertImageEquals(t, expected, out)
}

func TestTintTransformer(t *testing.T) {
	// transform all pixels to color #decba9 with original alpha
	transformer := TintTransformer(intAsColor(0xffdecba9))
	base := NewTestImage()

	out, err := transformer(base)
	assert.Nil(t, err)

	expected := NewTestImagePixels([][]int{
		{0xffdecba9, 0xffdecba9, 0xffdecba9, 0xffdecba9},
		{0xffdecba9, 0xffdecba9, 0xffdecba9, 0xffdecba9},
		{0xffdecba9, 0xffdecba9, 0xffdecba9, 0xffdecba9},
		{0xffdecba9, 0x00decba9, 0xffdecba9, 0x00decba9},
	})

	assertImageEquals(t, expected, out)
}

func TestResizeTransformer(t *testing.T) {
	transformer := ResizeTransformer(2, 2)
	base := NewTestImage()

	out, err := transformer(base)
	assert.Nil(t, err)

	expectedRect := image.Rect(0, 0, 2, 2)

	assert.Equal(t, expectedRect, out.Bounds())
}

func TestTransformImage(t *testing.T) {

	noOpTransform := func(i image.Image) (image.Image, error) {
		return i, nil
	}
	base := NewTestImage()

	out, err := TransformImage(base, noOpTransform, noOpTransform, noOpTransform)
	assert.Nil(t, err)

	assertImageEquals(t, base, out)
}

func TestTransformImage_withError(t *testing.T) {

	noOpTransform := func(i image.Image) (image.Image, error) {
		return i, nil
	}
	expectedFailure := errors.New("expected failure")
	errTransform := func(i image.Image) (image.Image, error) {
		return nil, expectedFailure
	}
	base := NewTestImage()

	out, err := TransformImage(base, noOpTransform, noOpTransform, errTransform, noOpTransform)
	if assert.Error(t, err) {
		assert.Equal(t, expectedFailure, err)
	}
	assert.Nil(t, out)
}
