package wo

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Input defines the mechanism to read user input from mice and keyboards.
// It is taken directly from the behavior of a pixel.Window.
type Input interface {
	// Pressed returns whether the Button is currently pressed down.
	Pressed(button pixelgl.Button) bool

	// JustPressed returns whether the Button has just been pressed down.
	JustPressed(button pixelgl.Button) bool

	// JustReleased returns whether the Button has just been released up.
	JustReleased(button pixelgl.Button) bool

	// Repeated returns whether a repeat event has been triggered on button.
	//
	// Repeat event occurs repeatedly when a button is held down for some time.
	Repeated(button pixelgl.Button) bool

	// MousePosition returns the current mouse position in the Window's Bounds.
	MousePosition() pixel.Vec

	// MouseScroll returns the mouse scroll amount (in both axes) since the last call to Window.Update.
	MouseScroll() pixel.Vec

	// Typed returns the text typed on the keyboard since the last call to Window.Update.
	Typed() string
}
