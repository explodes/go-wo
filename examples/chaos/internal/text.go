package internal

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

func drawText(canvas *pixelgl.Canvas, topLeft pixel.Vec, text *text.Text, lines ...string) {
	const (
		lineSpacing float64 = 2
	)
	offset := 0.0
	for _, line := range lines {
		text.Clear()
		text.WriteRune('\r')
		text.WriteString(line)
		bounds := text.Bounds()
		text.Draw(canvas, pixel.IM.Moved(pixel.V(topLeft.X, topLeft.Y-bounds.H()-offset)))
		offset += bounds.H() + lineSpacing
	}
}
