package wobj

import "github.com/faiface/pixel"

// Bounder is an interface for something that has a hit box
type Bounder interface {
	// Bounds returns the Rect that encompasses this object
	Bounds() pixel.Rect
}
