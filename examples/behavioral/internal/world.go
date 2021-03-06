package internal

import (
	"github.com/explodes/go-wo"
	"github.com/explodes/go-wo/examples/behavioral/res"
	"github.com/sirupsen/logrus"
)

const (
	title  = "Behavioral Objects Example"
	width  = 500
	height = 500
	fps    = 60
)

type World struct {
	loader wo.Loader
	debug  bool
}

func NewWorld(debug bool) *World {
	return &World{
		loader: wo.NewLoaderFromByteReader(res.Load),
		debug:  debug,
	}
}

func (w *World) Run() {
	scenes := map[string]wo.SceneFactory{
		"scene": w.newObjectsScene,
	}

	world, err := wo.NewWorld(title, width, height, scenes)
	if err != nil {
		logrus.Fatalf("error starting world: %v", err)
	}
	world.SetFps(fps)

	currentScene := "scene"

	for {
		log := logrus.WithField("currentScene", currentScene)

		result, err := world.RunScene(currentScene)
		if err != nil {
			log.Fatalf("failed to run scene: %v", err)
		}
		switch result {
		case wo.SceneResultError:
			log.Error("error running scene")
			return
		case wo.SceneResultWindowClosed:
			log.Info("goodbye!")
			return
		}
	}
}
