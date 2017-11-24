package internal

import (
	"github.com/explodes/go-wo"
	"github.com/sirupsen/logrus"
)

const (
	title  = "Flappy Gopher"
	width  = 1019
	height = 384
	fps    = 60

	noPreviousScore = -1
)

const (
	SceneResultGoToTitle wo.SceneResult = iota
	SceneResultGoToGame
)

type FlappyWorld struct {
	lastScore int
	log       *logrus.Logger
}

func NewFlappyWorld() *FlappyWorld {
	return &FlappyWorld{
		log:       logrus.StandardLogger(),
		lastScore: noPreviousScore,
	}
}

func (g *FlappyWorld) Run() {

	scenes := map[string]wo.SceneFactory{
		"title":       g.createTitleScene,
		"FlappyWorld": g.createPlayScene,
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
		case wo.SceneResultError:
			log.Error("error running scene")
			return
		case wo.SceneResultWindowClosed:
			log.Info("goodbye!")
			return
		case SceneResultGoToGame:
			log.Info("let the games begin...")
			currentScene = "FlappyWorld"
		case SceneResultGoToTitle:
			log.Info("FlappyWorld over")
			currentScene = "title"
		}
	}
}
