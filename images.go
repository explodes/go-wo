package wo

import (
	"bytes"
	"github.com/faiface/pixel"
	"image"
	"image/color"
)

type ImageTransformer func(image.Image) (image.Image, error)

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

func AlphaKeyTransformer(key color.Color) ImageTransformer {
	return func(base image.Image) (image.Image, error) {
		transformed := &alphaKeyImage{
			Image: base,
			key:   key,
		}
		return transformed, nil
	}
}

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

var alpha = color.Alpha{A: 0}

type alphaKeyImage struct {
	image.Image
	key color.Color
}

func (a *alphaKeyImage) At(x, y int) color.Color {
	base := a.Image.At(x, y)
	cr, cg, cb, ca := base.RGBA()
	kr, kg, kb, ka := a.key.RGBA()
	if cr == kr &&
		cg == kg &&
		cb == kb &&
		ca == ka {
		return alpha
	}
	return base
}