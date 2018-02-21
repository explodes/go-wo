package wobj

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLayers(t *testing.T) {
	layers := NewLayers(10)

	assert.Len(t, layers, 10)
}

func TestLayers_Draw(t *testing.T) {
	testObj := newTestObject("tag")
	layers := NewLayers(1)
	layers[0].Add(testObj.obj)

	layers.Draw(testTarget{})

	assert.Equal(t, 1, testObj.drawable.drawCount)
}

func TestLayers_Draw_no_drawable(t *testing.T) {
	testObj := newTestObject("tag")
	oldDrawable := testObj.drawable
	testObj.drawable = nil
	layers := NewLayers(1)
	layers[0].Add(testObj.obj)

	layers.Draw(testTarget{})

	assert.Equal(t, 1, oldDrawable.drawCount)
}

func TestLayers_Update(t *testing.T) {
	testObj := newTestObject("tag")
	layers := NewLayers(1)
	layers[0].Add(testObj.obj)

	assert.Equal(t, 0, testObj.preCount)
	assert.Equal(t, 0, testObj.stepCount)
	assert.Equal(t, 0, testObj.postCount)

	layers.Update(1)

	assert.Equal(t, 1, testObj.preCount)
	assert.Equal(t, 1, testObj.stepCount)
	assert.Equal(t, 1, testObj.postCount)
}

func TestNewObjects(t *testing.T) {
	objects := NewObjects()

	assert.NotNil(t, objects)
}

func TestObjects_Add_tagged(t *testing.T) {
	testObj := newTestObject("tag")
	objects := NewObjects()

	assert.False(t, objects.All().Contains(testObj.obj))
	assert.False(t, objects.Tagged("tag").Contains(testObj.obj))

	objects.Add(testObj.obj)

	assert.True(t, objects.All().Contains(testObj.obj))
	assert.True(t, objects.Tagged("tag").Contains(testObj.obj))
}

func TestObjects_Add_untagged(t *testing.T) {
	testObj := newTestObject("")
	objects := NewObjects()

	objects.Add(testObj.obj)

	assert.True(t, objects.All().Contains(testObj.obj))
	assert.Nil(t, objects.Tagged(""))
}

func TestObjects_Size(t *testing.T) {
	testObj := newTestObject("tag")
	objects := NewObjects()

	assert.Equal(t, 0, objects.Size())

	objects.Add(testObj.obj)

	assert.Equal(t, 1, objects.Size())
}

func TestObjects_Remove(t *testing.T) {
	testObj := newTestObject("tag")
	objects := NewObjects()

	assert.Equal(t, 0, objects.Size())

	objects.Add(testObj.obj)
	objects.Remove(testObj.obj)

	assert.Equal(t, 0, objects.Size())
}

func TestObjects_Contains(t *testing.T) {
	testObj := newTestObject("tag")
	objects := NewObjects()

	assert.False(t, objects.Contains(testObj.obj))

	objects.Add(testObj.obj)

	assert.True(t, objects.Contains(testObj.obj))
}

func TestObjectSet_Contains(t *testing.T) {
	testObj := newTestObject("tag")
	set := make(ObjectSet)

	assert.False(t, set.Contains(testObj.obj))

	set.add(testObj.obj)

	assert.True(t, set.Contains(testObj.obj))
}
