package internal

import (
	"github.com/explodes/go-wo"
	"github.com/explodes/go-wo/examples/tanks/res"
	"github.com/sirupsen/logrus"
)

const (
	title  = "Tanks"
	width  = 800
	height = 400
	fps    = 60

	gotoBattle wo.SceneResult = 1
	gotoTitle  wo.SceneResult = 2
)

type World struct {
	loader wo.Loader
	debug  bool

	blueScore int
	redScore  int
}

func NewWorld(debug bool) *World {
	return &World{
		loader: wo.NewLoaderFromByteReader(res.Load),
		debug:  debug,
	}
}

func (w *World) Run() {
	scenes := map[string]wo.SceneFactory{
		"title": w.newTitleScene,
		"game":  w.newGameScene,
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
		switch result {
		case gotoTitle:
			currentScene = "title"
		case gotoBattle:
			currentScene = "game"
		case wo.SceneResultError:
			log.Error("error running scene")
			return
		case wo.SceneResultWindowClosed:
			log.Info("goodbye!")
			return
		}
	}
}
