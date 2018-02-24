package internal

import (
	"github.com/explodes/go-wo"
	"github.com/explodes/go-wo/wobj"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	layerBackground = iota
	layerObjects
	layerForeground
	numLayers
)

type mainScene struct {
	w *World

	bounds pixel.Rect
	time   float64

	layers wobj.Layers
}

func (w *World) createMainScene(canvas *pixelgl.Canvas) (wo.Scene, error) {
	layers := wobj.NewLayers(numLayers)

	scene := &mainScene{
		w:      w,
		bounds: canvas.Bounds(),
		layers: layers,
	}
	return scene, nil
}

func (s *mainScene) Update(dt float64, input wo.Input) wo.SceneResult {
	s.time += dt

	if input.JustPressed(pixelgl.KeyEscape) {
		return wo.SceneResultWindowClosed
	}

	s.layers.Update(dt)

	return wo.SceneResultNone
}

func (s *mainScene) Draw(canvas *pixelgl.Canvas) {
	s.layers.Draw(canvas)
}
