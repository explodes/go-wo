package wo

import (
	"github.com/faiface/pixel"
	"math"
)

// Limits returns the limits of a rectangle
func Limits(rect pixel.Rect) (xMin, xMax, yMin, yMax float64) {
	return rect.Min.X, rect.Max.X, rect.Min.Y, rect.Max.Y
}

// Shape returns the shape of a rectangle
func Shape(rect pixel.Rect) (x, y, w, h float64) {
	return rect.Min.X, rect.Min.Y, rect.W(), rect.H()
}

// NegV rotates a vector 180 degrees
func NegV(v pixel.Vec) pixel.Vec {
	return pixel.V(-v.X, -v.Y)
}

// fit returns the Matrix that will transform a source Rect
// into the dest Rect
func fit(source, dest pixel.Rect) pixel.Matrix {
	xscale := dest.W() / source.W()
	yscale := dest.H() / source.H()
	scaleV := pixel.V(xscale, yscale)
	return pixel.IM.Moved(NegV(source.Min)).ScaledXY(pixel.ZV, scaleV).Moved(dest.Min)
}

// FitAtZero returns the Matrix that will transform a source Rect
// into the destination Rect and moves the source rectangle up and
// to the right by 50% of the source size which transforms
// the source origin from the center to the bottom left.
func FitAtZero(source, dest pixel.Rect) pixel.Matrix {
	return fit(source, dest).Moved(dest.Center())
}

// HitBox returns the hit box of an object based on its sprite and current position
func HitBox(sprite *pixel.Sprite, pos pixel.Vec) pixel.Rect {
	frame := sprite.Picture().Bounds()
	return frame.Moved(NegV(frame.Center())).Moved(pos)
}

// Collision returns if two rectangles intersect
func Collision(r1, r2 pixel.Rect) bool {
	if r1.Min.X > r2.Max.X || r2.Min.X > r1.Max.X {
		return false
	}
	if r1.Min.Y > r2.Max.Y || r2.Min.Y > r1.Max.Y {
		return false
	}
	return true
}

// DegToRad converts degrees to radians
func DegToRad(deg float64) (rad float64) {
	return deg * math.Pi / 180
}

// RadToDeg converts radians to degrees
func RadToDeg(rad float64) (deg float64) {
	return rad * 180 / math.Pi
}

// TranslateRect returns a rectangle that has
// been translated by a given offset
func TranslateRect(r pixel.Rect, delta pixel.Vec) pixel.Rect {
	return pixel.R(
		r.Min.X+delta.X,
		r.Min.Y+delta.Y,
		r.Max.X+delta.X,
		r.Max.Y+delta.Y,
	)
}

// ScaleRect returns a rectangle that has
// been scaled by a given percentage
func ScaleRect(r pixel.Rect, p float64) pixel.Rect {
	width := r.W() * p
	height := r.H() * p
	return SizedRect(r.Center(), pixel.V(width, height))
}

// SizedRect returns a rectangle that has
// a given center and given width and height
func SizedRect(center, size pixel.Vec) pixel.Rect {
	width := size.X / 2
	height := size.Y / 2
	return pixel.R(
		center.X-width,
		center.Y-height,
		center.X+width,
		center.Y+height,
	)
}
