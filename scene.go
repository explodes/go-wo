package wo

import (
	"github.com/faiface/pixel/pixelgl"
)

// SceneResult gives an indication as to how a Scene.Update(...) went.
//
// Negative values are reserved by the worldorder framework.
type SceneResult int64

const (
	// SceneResultNone indicates that a Scene has not yet reached a result
	// and should continue
	SceneResultNone SceneResult = -iota - 1
	// SceneResultError indicates that an error has occurred
	SceneResultError
	// SceneResultWindowClosed indicates that the window was closed
	SceneResultWindowClosed
)

// Scene is an interface for describing a stage
// that updates over time and drawn onto an IMDraw
type Scene interface {
	// Update updates the game state for
	// a given time-delta in seconds
	Update(dt float64, input Input) SceneResult

	// Draw draws onto a Canvas
	Draw(canvas *pixelgl.Canvas)
}

// SceneFactory builds scenes when it is time to use it
type SceneFactory func(canvas *pixelgl.Canvas) (Scene, error)
