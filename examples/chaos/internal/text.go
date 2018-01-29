package internal

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

func drawText(canvas *pixelgl.Canvas, text *text.Text, lines ...string) {
	const (
		initialOffset float64 = 10
		lineSpacing   float64 = 5
	)
	offset := initialOffset
	for _, line := range lines {
		text.Clear()
		text.WriteRune('\r')
		text.WriteString(line)
		bounds := text.Bounds()
		text.Draw(canvas, pixel.IM.Moved(pixel.V(initialOffset, canvas.Bounds().H()-bounds.H()-offset)))
		offset += bounds.H() + lineSpacing
	}
}
