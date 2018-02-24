package wobj

import (
	"testing"

	"github.com/faiface/pixel"
	"github.com/stretchr/testify/assert"
)

func TestBehaviors_Execute(t *testing.T) {
	count := 0
	behavior := func(source *Object, dt float64) {
		count++
	}
	b := MakeBehaviors(behavior, behavior)

	b.Execute(nil, 0)

	assert.Equal(t, 2, count)
}

func TestMovement(t *testing.T) {
	obj := newTestObject("")
	obj.obj.Pos = pixel.V(10, 10)
	obj.obj.Velocity = pixel.V(10, 10)

	Movement(obj.obj, 1)

	assert.Equal(t, pixel.V(20, 20), obj.obj.Pos)
}

func TestFaceDirection(t *testing.T) {
	obj := newTestObject("")
	obj.obj.Velocity = pixel.V(1, 1)

	FaceDirection(obj.obj, 1)

	assert.Equal(t, pixel.V(1, 1).Angle(), obj.obj.Rot)
}

func TestFaceDirection_zero_velocity(t *testing.T) {
	obj := newTestObject("")
	obj.obj.Velocity = pixel.V(0, 0)

	FaceDirection(obj.obj, 1)

	assert.Equal(t, 0.0, obj.obj.Rot)
}
