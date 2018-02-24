package internal

import (
	"math/rand"

	"time"

	"errors"

	"github.com/explodes/go-wo"
	"github.com/explodes/go-wo/templates/multiscene/res"
)

const (
	title  = "Multi-Scene"
	width  = 600
	height = 400
	fps    = 60
)

const (
	sceneResultGotoTitle wo.SceneResult = iota
	sceneResultGotoMain
)

const (
	layerBackground = iota
	layerObjects
	layerForeground
	numLayers
)

type World struct {
	loader wo.Loader
	debug  bool
	rng    *rand.Rand
	input  wo.Input
}

func NewWorld(debug bool) *World {
	return &World{
		loader: wo.NewLoaderFromByteReader(res.Load),
		debug:  debug,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (w *World) Run() error {
	scenes := map[string]wo.SceneFactory{
		"title": w.createTitleScene,
		"main":  w.createMainScene,
	}
	world, err := wo.NewWorld(title, width, height, scenes)
	if err != nil {
		return err
	}
	world.SetFps(fps)
	w.input = world.Input()

	currentScene := "title"
	for {
		result, err := world.RunScene(currentScene)
		if err != nil {
			return err
		}

		switch result {
		case sceneResultGotoTitle:
			currentScene = "title"
		case sceneResultGotoMain:
			currentScene = "main"
		case wo.SceneResultError:
			return errors.New("error running scene")
		case wo.SceneResultWindowClosed:
			return nil
		}
	}
}
