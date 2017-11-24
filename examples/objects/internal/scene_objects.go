package internal

import (
	"github.com/explodes/go-wo"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type scene struct {
	time   float64
	bounds pixel.Rect

	obj *Obj
}

var _ wo.Scene = &scene{}

func (w *World) newSpaceScene(canvas *pixelgl.Canvas) (wo.Scene, error) {

	obj, err := newObj(w.loader)
	if err != nil {
		return nil, err
	}
	scene := &scene{
		bounds: canvas.Bounds(),
		obj:    obj,
	}
	return scene, nil
}

func (s *scene) Update(dt float64, input wo.Input) wo.SceneResult {
	s.time += dt
	s.obj.update(dt, input)
	return wo.SceneResultNone
}

func (s *scene) Draw(canvas *pixelgl.Canvas) {
	s.obj.draw(canvas)
}
