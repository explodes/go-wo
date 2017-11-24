package wo

import (
	"bytes"
	"github.com/faiface/pixel"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"io"
	"io/ioutil"
)

type ByteReader func(name string) ([]byte, error)

type AssetReader func(name string) (io.Reader, error)

type Loader struct {
	reader       AssetReader
	sprites      map[string]*pixel.Sprite
	spriteSheets map[string]*SpriteSheet
	fonts        map[string]*truetype.Font
}

func NewLoaderFromByteReader(reader ByteReader) *Loader {
	return NewLoader(func(name string) (io.Reader, error) {
		b, err := reader(name)
		if err != nil {
			return nil, err
		}
		return bytes.NewReader(b), nil
	})
}

func NewLoader(reader AssetReader) *Loader {
	return &Loader{
		reader:       reader,
		sprites:      make(map[string]*pixel.Sprite),
		spriteSheets: make(map[string]*SpriteSheet),
		fonts:        make(map[string]*truetype.Font),
	}
}

func (load *Loader) bytesOf(name string) ([]byte, error) {
	r, err := load.reader(name)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}

func (load *Loader) Free(name string) {
	delete(load.sprites, name)
	delete(load.spriteSheets, name)
	delete(load.fonts, name)
}

func (load *Loader) Sprite(name string, transforms ...ImageTransformer) (*pixel.Sprite, error) {
	if sprite, ok := load.sprites[name]; ok {
		return sprite, nil
	}
	r, err := load.reader(name)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	img, err = transformImage(img, transforms...)
	if err != nil {
		return nil, err
	}
	pic := pixel.PictureDataFromImage(img)
	sprite := pixel.NewSprite(pic, pic.Bounds())
	load.sprites[name] = sprite
	return sprite, nil
}

func (load *Loader) SpriteSheet(name string, opts SpriteSheetOptions, transforms ...ImageTransformer) (*SpriteSheet, error) {
	if sheet, ok := load.spriteSheets[name]; ok {
		return sheet, nil
	}
	r, err := load.reader(name)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	img, err = transformImage(img, transforms...)
	if err != nil {
		return nil, err
	}
	pic := pixel.PictureDataFromImage(img)
	sheet := NewSpriteSheet(pic, opts)
	load.spriteSheets[name] = sheet
	return sheet, nil
}

func (load *Loader) Font(name string) (*truetype.Font, error) {
	if f, ok := load.fonts[name]; ok {
		return f, nil
	}
	b, err := load.bytesOf(name)
	if err != nil {
		return nil, err
	}
	f, err := truetype.Parse(b)
	if err != nil {
		return nil, err
	}
	load.fonts[name] = f
	return f, err
}

func (load *Loader) FontFace(name string, size float64) (font.Face, error) {
	f, err := load.Font(name)
	if err != nil {
		return nil, err
	}
	return truetype.NewFace(f, &truetype.Options{Size: size}), nil
}
