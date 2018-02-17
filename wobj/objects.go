package wobj

import "github.com/faiface/pixel"

type Objects struct {
	all    ObjectSet
	tagged objectTagMap
}

func NewObjects() *Objects {
	return &Objects{
		all:    make(ObjectSet),
		tagged: make(objectTagMap),
	}
}

func (o *Objects) Size() int {
	return len(o.all)
}

func (o *Objects) All() ObjectSet {
	return o.all
}

func (o *Objects) Tagged(tag string) ObjectSet {
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
	for object := range o.all {
		object.Draw(target)
	}
}

func (o *Objects) PreStep(dt float64) {
	for object := range o.all {
		object.PreSteps.Execute(object, dt)
	}
}

func (o *Objects) Step(dt float64) {
	for object := range o.all {
		object.Steps.Execute(object, dt)
	}
}

func (o *Objects) PostStep(dt float64) {
	for object := range o.all {
		object.PostSteps.Execute(object, dt)
	}
}

type ObjectSet map[*Object]struct{}

func (os ObjectSet) Contains(obj *Object) bool {
	if os == nil {
		return false
	}
	_, exists := os[obj]
	return exists
}

func (os ObjectSet) add(obj *Object) {
	os[obj] = struct{}{}
}

func (os ObjectSet) remove(obj *Object) {
	delete(os, obj)
}

type objectTagMap map[string]ObjectSet

func (m objectTagMap) add(key string, obj *Object) {
	set := m[key]
	if set == nil {
		set = make(ObjectSet)
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
