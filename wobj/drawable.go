package wobj

import (
	wo "github.com/explodes/go-wo"
	"github.com/faiface/pixel"
)

type Drawable interface {
	Draw(pixel.Target, pixel.Matrix)
	Bounds() pixel.Rect
}

type SpriteDrawable struct {
	Sprite *pixel.Sprite
}

func (s *SpriteDrawable) Bounds() pixel.Rect {
	return s.Sprite.Picture().Bounds()
}

func (s *SpriteDrawable) Draw(target pixel.Target, matrix pixel.Matrix) {
	s.Sprite.Draw(target, matrix)
}

type SpriteSheetDrawable struct {
	Sheet *wo.SpriteSheet
}

func (s *SpriteSheetDrawable) Bounds() pixel.Rect {
	frame := s.Sheet.Sprite().Frame()
	return pixel.R(0, 0, frame.W(), frame.H())
}

func (s *SpriteSheetDrawable) Draw(target pixel.Target, mat pixel.Matrix) {
	s.Sheet.Sprite().Draw(target, mat)
}
