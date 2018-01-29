package internal

import (
	"math/rand"
	"time"

	"github.com/explodes/go-wo"
	"github.com/explodes/go-wo/examples/chaos/res"
	"github.com/faiface/pixel/pixelgl"
	"github.com/sirupsen/logrus"
)

const (
	title  = "Chaos"
	width  = 920
	height = width
	fps    = 60
)

const (
	gotoTitle wo.SceneResult = iota
	gotoSierpinski
	gotoCircles
)

var (
	sceneBindings = map[pixelgl.Button]wo.SceneResult{
		pixelgl.Key0: gotoTitle,
		pixelgl.Key1: gotoSierpinski,
		pixelgl.Key2: gotoCircles,
	}
)

type World struct {
	debug          bool
	loader         wo.Loader
	rng            *rand.Rand
	inputWaitFrame bool
}

func NewWorld(debug bool) *World {
	return &World{
		debug:  debug,
		loader: wo.NewLoaderFromByteReader(res.Load),
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (w *World) Run() {
	scenes := map[string]wo.SceneFactory{
		"title":      w.newTitleScene,
		"sierpinski": w.newSierpinskiScene,
		"circles":    w.newCirclesScene,
	}

	world, err := wo.NewWorld(title, width, height, scenes)
	if err != nil {
		logrus.Fatalf("error starting world: %v", err)
	}
	world.SetFps(fps)

	currentScene := "title"

	for {
		log := logrus.WithField("currentScene", currentScene)

		result, err := world.RunScene(currentScene)
		if err != nil {
			log.Fatalf("failed to run scene: %v", err)
		}
		log = log.WithField("result", result)
		switch result {
		case wo.SceneResultError:
			log.Error("error running scene")
			return
		case wo.SceneResultWindowClosed:
			log.Info("goodbye!")
			return
		case gotoTitle:
			log.Info("Title")
			currentScene = "title"
		case gotoSierpinski:
			log.Info("Sierpinski")
			currentScene = "sierpinski"
		case gotoCircles:
			log.Info("Circles")
			currentScene = "circles"
		}
	}
}

func (w *World) maybeSelectScene(input wo.Input) wo.SceneResult {
	if w.inputWaitFrame {
		w.inputWaitFrame = false
		return wo.SceneResultNone
	}
	if input.Pressed(pixelgl.KeyEscape) {
		return wo.SceneResultWindowClosed
	}
	for key, scene := range sceneBindings {
		if input.JustPressed(key) {
			w.inputWaitFrame = true
			return scene
		}
	}
	return wo.SceneResultNone
}
