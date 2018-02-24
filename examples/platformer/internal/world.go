package internal

import (
	"math/rand"

	"time"

	"errors"

	"github.com/explodes/go-wo"
	"github.com/explodes/go-wo/examples/platformer/res"
)

const (
	title  = "Platformer"
	width  = 600
	height = width
	fps    = 60
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
		"main": w.createMainScene,
	}
	world, err := wo.NewWorld(title, width, height, scenes)
	if err != nil {
		return err
	}
	world.SetFps(fps)
	w.input = world.Input()

	currentScene := "main"
	for {
		result, err := world.RunScene(currentScene)
		if err != nil {
			return err
		}
		switch result {
		case wo.SceneResultError:
			return errors.New("error running scene")
		case wo.SceneResultWindowClosed:
			return nil
		}
	}
}
