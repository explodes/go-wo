package wobj

import (
	"testing"

	"github.com/faiface/pixel"
	"github.com/stretchr/testify/assert"
)

func TestMovement(t *testing.T) {
	obj := newTestObject("")
	obj.obj.Pos = pixel.V(10, 10)
	obj.obj.Velocity = pixel.V(10, 10)

	Movement(obj.obj, 1)

	assert.Equal(t, pixel.V(20, 20), obj.obj.Pos)
}

func TestCollision_hit(t *testing.T) {
	obj1 := newTestObject("")
	obj2 := newTestObject("")
	obj1.obj.Pos = pixel.V(0, 0)
	obj1.obj.Size = pixel.V(10, 10)
	obj2.obj.Pos = pixel.V(0, 0)
	obj2.obj.Size = pixel.V(10, 10)

	hitTest := false
	Collision(obj2.obj, func(*Object, float64) {
		hitTest = true
	})(obj1.obj, 1)

	assert.True(t, hitTest)
}

func TestCollision_miss(t *testing.T) {
	obj1 := newTestObject("")
	obj2 := newTestObject("")
	obj1.obj.Pos = pixel.V(0, 0)
	obj1.obj.Size = pixel.V(10, 10)
	obj2.obj.Pos = pixel.V(1000, 1000)
	obj2.obj.Size = pixel.V(10, 10)

	hitTest := false
	Collision(obj2.obj, func(*Object, float64) {
		hitTest = true
	})(obj1.obj, 1)

	assert.False(t, hitTest)
}

func TestObjectCollision_hit(t *testing.T) {
	obj1 := newTestObject("")
	obj2 := newTestObject("")
	obj1.obj.Pos = pixel.V(0, 0)
	obj1.obj.Size = pixel.V(10, 10)
	obj2.obj.Pos = pixel.V(0, 0)
	obj2.obj.Size = pixel.V(10, 10)

	hitTest := false
	ObjectCollision(obj2.obj, func(*Object, *Object, float64) {
		hitTest = true
	})(obj1.obj, 1)

	assert.True(t, hitTest)
}

func TestObjectCollision_miss(t *testing.T) {
	obj1 := newTestObject("")
	obj2 := newTestObject("")
	obj1.obj.Pos = pixel.V(0, 0)
	obj1.obj.Size = pixel.V(10, 10)
	obj2.obj.Pos = pixel.V(1000, 1000)
	obj2.obj.Size = pixel.V(10, 10)

	hitTest := false
	ObjectCollision(obj2.obj, func(*Object, *Object, float64) {
		hitTest = true
	})(obj1.obj, 1)

	assert.False(t, hitTest)
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

func TestFaceDirectionOffset(t *testing.T) {
	obj := newTestObject("")
	obj.obj.Velocity = pixel.V(1, 1)

	FaceDirectionOffset(45)(obj.obj, 1)

	assert.Equal(t, pixel.V(1, 1).Angle()+45, obj.obj.Rot)
}
