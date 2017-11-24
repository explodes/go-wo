package wo

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/colornames"
	"image/color"
	"time"
)

const (
	defaultFps = 60
)

func Run(run func()) {
	pixelgl.Run(run)
}

type World struct {
	window *pixelgl.Window
	canvas *pixelgl.Canvas
	fit    pixel.Matrix

	Color color.Color

	fps *FpsLimiter

	scenes map[string]SceneFactory
}

func NewWorld(title string, width, height int, scenes map[string]SceneFactory) (*World, error) {

	w := float64(width)
	h := float64(height)

	cfg := pixelgl.WindowConfig{
		Title:     title,
		Bounds:    pixel.R(0, 0, w, h),
		VSync:     true,
		Resizable: false,
	}
	window, err := pixelgl.NewWindow(cfg)
	if err != nil {
		return nil, err
	}
	window.SetSmooth(true)

	canvas := pixelgl.NewCanvas(pixel.R(0, 0, w, h))

	world := &World{
		window: window,
		canvas: canvas,
		scenes: scenes,
		Color:  colornames.Black,
		fps:    NewFpsLimiter(defaultFps),
		fit:    FitAtZero(canvas.Bounds(), window.Bounds()),
	}
	return world, nil
}

func (w *World) createScene(name string) (Scene, error) {
	logrus.WithFields(logrus.Fields{
		"name": name,
	}).Debug("creating scene")

	factory, ok := w.scenes[name]
	if !ok {
		return nil, errors.Errorf("scene %s does not exist", name)
	}
	start := time.Now()
	scene, err := factory(w.canvas)

	logrus.WithFields(logrus.Fields{
		"creationDuration": time.Now().Sub(start),
		"err":              err,
		"name":             name,
	}).Debug("created scene")

	return scene, err
}

func (w *World) SetFps(maxFps float64) {
	w.fps.SetLimit(maxFps)
}

func (w *World) RunScene(name string) (SceneResult, error) {
	scene, err := w.createScene(name)
	if err != nil {
		return SceneResultError, errors.Errorf("unable to create scene %s: %v", name, err)
	}
	return w.completeScene(scene), nil
}

func (w *World) completeScene(scene Scene) SceneResult {
	w.fps.Reset()
	for !w.window.Closed() {
		dt := w.fps.StartFrame()
		result := scene.Update(dt, w.window)
		if result != SceneResultNone {
			return result
		}

		if w.Color != nil {
			w.canvas.Clear(w.Color)
		}
		w.canvas.SetMatrix(pixel.IM)
		scene.Draw(w.canvas)

		w.drawToWindow()
		w.fps.WaitForNextFrame()
	}
	return SceneResultWindowClosed
}

func (w *World) drawToWindow() {
	w.canvas.Draw(w.window, w.fit)
	w.window.Update()
}
