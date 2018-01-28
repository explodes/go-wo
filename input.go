package wo

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Input defines the mechanism to read user input from mice and keyboards.
type Input interface {
	// Pressed returns whether any Button is currently pressed down.
	Pressed(button ...pixelgl.Button) bool

	// JustPressed returns whether any Button has just been pressed down.
	JustPressed(button ...pixelgl.Button) bool

	// JustReleased returns whether any Button has just been released up.
	JustReleased(button ...pixelgl.Button) bool

	// Repeated returns whether a repeat event has been triggered on any button.
	//
	// Repeat event occurs repeatedly when a button is held down for some time.
	Repeated(button ...pixelgl.Button) bool

	// MousePosition returns the current mouse position in the Window's Bounds.
	MousePosition() pixel.Vec

	// MouseScroll returns the mouse scroll amount (in both axes) since the last call to Window.Update.
	MouseScroll() pixel.Vec

	// Typed returns the text typed on the keyboard since the last call to Window.Update.
	Typed() string
}

type windowInput struct {
	win *pixelgl.Window
}

func (w *windowInput) Pressed(button ...pixelgl.Button) bool {
	for _, b := range button {
		if w.win.Pressed(b) {
			return true
		}
	}
	return false
}

func (w *windowInput) JustPressed(button ...pixelgl.Button) bool {
	for _, b := range button {
		if w.win.JustPressed(b) {
			return true
		}
	}
	return false
}

func (w *windowInput) JustReleased(button ...pixelgl.Button) bool {
	for _, b := range button {
		if w.win.JustReleased(b) {
			return true
		}
	}
	return false
}

func (w *windowInput) Repeated(button ...pixelgl.Button) bool {
	for _, b := range button {
		if w.win.Repeated(b) {
			return true
		}
	}
	return false
}

func (w *windowInput) MousePosition() pixel.Vec {
	return w.win.MousePosition()
}

func (w *windowInput) MouseScroll() pixel.Vec {
	return w.win.MouseScroll()
}

func (w *windowInput) Typed() string {
	return w.win.Typed()
}
