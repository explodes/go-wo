package wo

import (
	"bytes"
	"image"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/nfnt/resize"
)

// DecodePicture decodes PictureData from bytes. The image
// can also be transformed ahead of time using ImageTransformers.
func DecodePicture(b []byte, transforms ...ImageTransformer) (*pixel.PictureData, error) {
	r := bytes.NewReader(b)
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	img, err = transformImage(img, transforms...)
	if err != nil {
		return nil, err
	}
	pic := pixel.PictureDataFromImage(img)
	return pic, nil
}

// ImageTransformer transforms one image into another.
type ImageTransformer func(image.Image) (image.Image, error)

// AlphaKeyTransformer transforms an image into an image
// where pixels of a certain color become transparent.
func AlphaKeyTransformer(key color.Color) ImageTransformer {
	return func(base image.Image) (image.Image, error) {
		transformed := &alphaKeyImage{
			Image: base,
			key:   key,
		}
		return transformed, nil
	}
}

func ResizeTransformer(width, height uint) ImageTransformer {
	return func(base image.Image) (image.Image, error) {
		transformed := resize.Resize(width, height, base, resize.Bicubic)
		return transformed, nil
	}
}

// transformImage applies a chain of image transformations
// to a given image.
func transformImage(base image.Image, transforms ...ImageTransformer) (image.Image, error) {
	var err error
	for _, transform := range transforms {
		base, err = transform(base)
		if err != nil {
			return nil, err
		}
	}
	return base, nil
}

// transparent is a constant for completely transparent colors.
var transparent = color.Alpha{A: 0}

// alphaKeyImage wraps an image and turns pixels of
// a certain color into transparent pixels.
type alphaKeyImage struct {
	image.Image
	key color.Color
}

// At overrides Image's At an returns a transparent color
// for pixels that match the alpha key.
func (a *alphaKeyImage) At(x, y int) color.Color {
	base := a.Image.At(x, y)
	cr, cg, cb, ca := base.RGBA()
	kr, kg, kb, ka := a.key.RGBA()
	if cr == kr &&
		cg == kg &&
		cb == kb &&
		ca == ka {
		return transparent
	}
	return base
}
