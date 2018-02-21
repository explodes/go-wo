package wobj

import (
	"image"

	"github.com/faiface/pixel"
)

type testObj struct {
	obj *Object

	preCount  int
	stepCount int
	postCount int

	drawable *testDrawable
}

type testDrawable struct {
	drawCount int
}

func (t *testDrawable) Draw(pixel.Target, pixel.Matrix) {
	t.drawCount++
}

func (t *testDrawable) Bounds() pixel.Rect {
	return pixel.R(0, 0, 10, 10)
}

func newTestObject(tag string) *testObj {
	drawable := &testDrawable{}

	testObject := &testObj{
		drawable: drawable,
	}

	testObject.obj = &Object{
		Tag:      tag,
		Drawable: drawable,
		PreSteps: MakeBehaviors(func(source *Object, dt float64) {
			testObject.preCount++
		}),
		Steps: MakeBehaviors(func(source *Object, dt float64) {
			testObject.stepCount++
		}),
		PostSteps: MakeBehaviors(func(source *Object, dt float64) {
			testObject.postCount++
		}),
	}

	return testObject
}

type testTarget struct{}

func (testTarget) MakeTriangles(pixel.Triangles) pixel.TargetTriangles {
	return nil
}

func (testTarget) MakePicture(pixel.Picture) pixel.TargetPicture {
	return nil
}

func newTestPicture() pixel.Picture {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	return pixel.PictureDataFromImage(img)
}
