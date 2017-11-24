package wo

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func DecodeFontFace(ttf []byte, size float64) (font.Face, error) {
	f, err := truetype.Parse(ttf)
	if err != nil {
		return nil, err
	}
	return truetype.NewFace(f, &truetype.Options{Size: size}), nil
}
func MustDecodeFontFace(ttf []byte, size float64) font.Face {
	ff, err := DecodeFontFace(ttf, size)
	if err != nil {
		panic(err)
	}
	return ff
}
