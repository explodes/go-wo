package internal

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type IMDrawable struct {
	im     *imdraw.IMDraw
	bounds pixel.Rect
}

func (i *IMDrawable) Draw(target pixel.Target, mat pixel.Matrix) {
	c := target.(*pixelgl.Canvas)
	mat = mat.Moved(pixel.ZV.Sub(i.bounds.Center()))
	c.SetMatrix(mat)
	i.im.Draw(c)
	c.SetMatrix(pixel.IM)
}

func (i *IMDrawable) Bounds() pixel.Rect {
	return i.bounds
}
