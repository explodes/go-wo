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
	testObj.obj.Drawable = nil
	layers := NewLayers(1)
	layers[0].Add(testObj.obj)

	layers.Draw(testTarget{})

	assert.Equal(t, 0, testObj.drawable.drawCount)
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

func TestObjects_Len(t *testing.T) {
	testObj := newTestObject("tag")
	objects := NewObjects()

	assert.Equal(t, 0, objects.Len())

	objects.Add(testObj.obj)

	assert.Equal(t, 1, objects.Len())
}

func TestObjects_Remove(t *testing.T) {
	testObj := newTestObject("tag")
	objects := NewObjects()

	assert.Equal(t, 0, objects.Len())

	objects.Add(testObj.obj)
	objects.Remove(testObj.obj)

	assert.Equal(t, 0, objects.Len())
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
	set := NewObjectSet()

	assert.False(t, set.Contains(testObj.obj))

	set.add(testObj.obj)

	assert.True(t, set.Contains(testObj.obj))
}

func TestObjectSet_Iterator_nil(t *testing.T) {
	var set *ObjectSet = nil

	iter := set.Iterator()
	iterSize := countIterator(iter)

	assert.Equal(t, 0, iterSize)
}

func TestObjectSet_Len_nil(t *testing.T) {
	var set *ObjectSet = nil

	assert.Equal(t, 0, set.Len())
}

func TestObjectSet_Contains_nil(t *testing.T) {
	var set *ObjectSet = nil

	assert.False(t, set.Contains(nil))
}

func countIterator(iter ObjectIterator) int {
	count := 0
	for _, ok := iter(); ok; _, ok = iter() {
		count++
	}
	return count
}

func TestLayers_Iterator(t *testing.T) {
	layers := NewLayers(3)
	layers[0].Add(newTestObject("a").obj)
	layers[1].Add(newTestObject("b").obj)
	layers[2].Add(newTestObject("b").obj)

	iter := layers.Iterator()
	iterSize := countIterator(iter)

	assert.Equal(t, 3, iterSize)
}

func TestLayers_Iterator_skipLayer(t *testing.T) {
	layers := NewLayers(3)
	layers[0].Add(newTestObject("a").obj)
	layers[2].Add(newTestObject("b").obj)

	iter := layers.Iterator()
	iterSize := countIterator(iter)

	assert.Equal(t, 2, iterSize)
}

func TestLayers_TagIterator(t *testing.T) {
	layers := NewLayers(3)
	layers[0].Add(newTestObject("a").obj)
	layers[1].Add(newTestObject("b").obj)
	layers[2].Add(newTestObject("b").obj)

	iter := layers.TagIterator("b")
	iterSize := countIterator(iter)

	assert.Equal(t, 2, iterSize)
}
