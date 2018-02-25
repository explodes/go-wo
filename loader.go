package wo

import (
	"bytes"
	"image"
	"io"
	"io/ioutil"

	"github.com/faiface/pixel"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// Loader is a utility for loading and caching assets such as
// Sounds, Sprites, SpriteSheets, Fonts, and FontFaces.
type Loader interface {
	// Sprite loads a Sprite by name and applies any custom transformations to the image.
	Sprite(name string, transforms ...ImageTransformer) (*pixel.Sprite, error)

	// SpriteSheet loads a SpriteSheet with the given options and
	// applies any custom transformations to the underlying image.
	SpriteSheet(name string, opts SpriteSheetOptions, transforms ...ImageTransformer) (*SpriteSheet, error)

	// Sound loads a Sound for a given format ("wav"/"mp3").
	Sound(format string, name string) (*Sound, error)

	// Font loads a truetype Font by name.
	Font(name string) (*truetype.Font, error)

	// FontFace loads FontFace for a truetype Font by name.
	FontFace(name string, size float64) (font.Face, error)
}

var _ Loader = &simpleLoader{}

// simpleLoader makes a simple function perform the
// heavy lifting of loading.
type simpleLoader struct {
	reader AssetReader
}

// NewLoaderFromByteReader creates a new Loader from a function
// that can acquire bytes by name.
//
// Sprites and SpriteSheets can be in any image format as long
// as it is loaded ahead of time (import _ "image/png").
//
// Sounds can be in "mp3" or "wav" format, their format is
// specified when acquiring the asset.
//
// Fonts are expected to be in truetype format.
func NewLoaderFromByteReader(reader ByteReader) Loader {
	return NewLoader(func(name string) (io.Reader, error) {
		b, err := reader(name)
		if err != nil {
			return nil, err
		}
		return bytes.NewReader(b), nil
	})
}

// NewLoader creates a new Loader from a function
// that can acquire an io.Reader by name.
//
// Sprites and SpriteSheets can be in any image format as long
// as it is loaded ahead of time (import _ "image/png").
//
// Sounds can be in "mp3" or "wav" format, their format is
// specified when acquiring the asset.
//
// Fonts are expected to be in truetype format.
func NewLoader(reader AssetReader) Loader {
	return &simpleLoader{
		reader: reader,
	}
}

// readCloser creates a ReadCloser for the given name. This is so that
// any Reader created from the AssetReader gets closed appropriately,
// should that Reader support Close.
func (load *simpleLoader) readCloser(name string) (io.ReadCloser, error) {
	r, err := load.reader(name)
	if err != nil {
		return nil, err
	}
	return &readCloserWrapper{r}, nil
}

// bytesOf reads all bytes for a given name.
func (load *simpleLoader) bytesOf(name string) ([]byte, error) {
	r, err := load.readCloser(name)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}

func (load *simpleLoader) Sprite(name string, transforms ...ImageTransformer) (*pixel.Sprite, error) {
	r, err := load.readCloser(name)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	img, err = TransformImage(img, transforms...)
	if err != nil {
		return nil, err
	}
	pic := pixel.PictureDataFromImage(img)
	sprite := pixel.NewSprite(pic, pic.Bounds())
	return sprite, nil
}

func (load *simpleLoader) SpriteSheet(name string, opts SpriteSheetOptions, transforms ...ImageTransformer) (*SpriteSheet, error) {
	r, err := load.readCloser(name)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	img, err = TransformImage(img, transforms...)
	if err != nil {
		return nil, err
	}
	pic := pixel.PictureDataFromImage(img)
	sheet := NewSpriteSheet(pic, opts)
	return sheet, nil
}

func (load *simpleLoader) Sound(format string, name string) (*Sound, error) {
	b, err := load.bytesOf(name)
	if err != nil {
		return nil, err
	}
	return NewSound(format, b)
}

func (load *simpleLoader) Font(name string) (*truetype.Font, error) {
	b, err := load.bytesOf(name)
	if err != nil {
		return nil, err
	}
	f, err := truetype.Parse(b)
	if err != nil {
		return nil, err
	}
	return f, err
}

func (load *simpleLoader) FontFace(name string, size float64) (font.Face, error) {
	f, err := load.Font(name)
	if err != nil {
		return nil, err
	}
	return truetype.NewFace(f, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}
