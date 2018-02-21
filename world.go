package wo

import (
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/colornames"
)

const (
	defaultFps = 60
)

// Run runs a world loop in the OpenGL context using pixel.
func Run(run func()) {
	pixelgl.Run(run)
}

// World is a container of Scenes that are rendered onto
// a window. It is used as the highest level entry point
// into the graphical programming of an application.
type World struct {
	window *pixelgl.Window
	input  Input
	canvas *pixelgl.Canvas
	fit    pixel.Matrix

	Color color.Color

	fps *FpsLimiter

	scenes map[string]SceneFactory
}

// NewWorld creates a world with a displayed window.
// Use RunScene(name) to begin rendering a Scene.
//
// SceneFactories are used to load Scenes by name.
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
		input:  &windowInput{window},
		canvas: canvas,
		scenes: scenes,
		Color:  colornames.Black,
		fps:    NewFpsLimiter(defaultFps),
		fit:    FitAtZero(canvas.Bounds(), window.Bounds()),
	}
	return world, nil
}

// RunScene renders a Scene until that Scene returns
// a SceneResult other than SceneResultNone.
func (w *World) RunScene(name string) (SceneResult, error) {
	scene, err := w.createScene(name)
	if err != nil {
		return SceneResultError, errors.Errorf("unable to create scene %s: %v", name, err)
	}
	return w.runToCompletion(scene), nil
}

// Input gets the World's Input.
func (w *World) Input() Input {
	return w.input
}

// SetFps sets the target frames per second to render scenes.
func (w *World) SetFps(maxFps float64) {
	w.fps.SetLimit(maxFps)
}

// createScene loads a Scene by name using its respective SceneFactory.
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

// runToCompletion runs the update/draw cycle on a Scene until
// that Scene returns a result other than SceneResultNone.
func (w *World) runToCompletion(scene Scene) SceneResult {
	w.fps.Reset()
	for !w.window.Closed() {
		dt := w.fps.StartFrame()
		result := scene.Update(dt, w.input)
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

// drawToWindow renders the canvas onto the window.
func (w *World) drawToWindow() {
	w.canvas.Draw(w.window, w.fit)
	w.window.Update()
}
