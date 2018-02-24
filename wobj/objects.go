package wobj

import (
	"github.com/faiface/pixel"
)

type Layers []*Objects

func NewLayers(n int) Layers {
	layers := make(Layers, 0, n)
	for i := 0; i < n; i++ {
		layers = append(layers, NewObjects())
	}
	return layers
}

func (ly Layers) Update(dt float64) {
	for _, layer := range ly {
		layer.Update(dt)
	}
}

func (ly Layers) Draw(target pixel.Target) {
	for _, layer := range ly {
		layer.Draw(target)
	}
}

type Objects struct {
	all    *ObjectSet
	tagged objectTagMap
}

func NewObjects() *Objects {
	return &Objects{
		all:    NewObjectSet(),
		tagged: make(objectTagMap),
	}
}

func (o *Objects) Size() int {
	return o.all.Size()
}

func (o *Objects) All() *ObjectSet {
	return o.all
}

func (o *Objects) Tagged(tag string) *ObjectSet {
	return o.tagged[tag]
}

func (o *Objects) Add(obj *Object) {
	o.all.add(obj)
	if obj.Tag != "" {
		o.tagged.add(obj.Tag, obj)
	}
}

func (o *Objects) Remove(obj *Object) {
	o.all.remove(obj)
	if obj.Tag != "" {
		o.tagged.remove(obj.Tag, obj)
	}
}

func (o *Objects) Contains(obj *Object) bool {
	return o.all.Contains(obj)
}

func (o *Objects) Update(dt float64) {
	o.PreStep(dt)
	o.Step(dt)
	o.PostStep(dt)
}

func (o *Objects) Draw(target pixel.Target) {
	for _, object := range o.all.Iterable() {
		object.Draw(target)
	}
}

func (o *Objects) PreStep(dt float64) {
	for _, object := range o.all.Iterable() {
		object.PreSteps.Execute(object, dt)
	}
}

func (o *Objects) Step(dt float64) {
	for _, object := range o.all.Iterable() {
		object.Steps.Execute(object, dt)
	}
}

func (o *Objects) PostStep(dt float64) {
	for _, object := range o.all.Iterable() {
		object.PostSteps.Execute(object, dt)
	}
}

type ObjectSet struct {
	set []*Object
}

func NewObjectSet() *ObjectSet {
	return &ObjectSet{
		set: make([]*Object, 0, 32),
	}
}

func (os *ObjectSet) Iterable() []*Object {
	if os == nil {
		return nil
	}
	return os.set
}

func (os *ObjectSet) Size() int {
	if os == nil {
		return 0
	}
	return len(os.set)
}

func (os *ObjectSet) Contains(obj *Object) bool {
	return os.index(obj) != -1
}

func (os *ObjectSet) index(obj *Object) int {
	if os == nil {
		return -1
	}
	for index, o := range os.set {
		if o == obj {
			return index
		}
	}
	return -1
}

func (os *ObjectSet) add(obj *Object) {
	if !os.Contains(obj) {
		os.set = append(os.set, obj)
	}
}

func (os *ObjectSet) remove(obj *Object) {
	if index := os.index(obj); index != -1 {
		os.set = append(os.set[:index], os.set[index+1:]...)
	}
}

type objectTagMap map[string]*ObjectSet

func (m objectTagMap) add(key string, obj *Object) {
	set := m[key]
	if set == nil {
		set = NewObjectSet()
		m[key] = set
	}
	set.add(obj)
}

func (m objectTagMap) remove(key string, obj *Object) {
	set := m[key]
	if set != nil {
		set.remove(obj)
	}
}
