package gogm

import "reflect"

type lifeCycle int

const (
	//CREATE is an event life cycle indicating creation of an object
	CREATE lifeCycle = iota
	LOAD             //READ
	UPDATE
	DELETE
)

type Event interface {
	GetObject() interface{}
	GetLifeCycle() lifeCycle
}

type event struct {
	object    *reflect.Value
	lifeCycle lifeCycle
}

func (e event) GetObject() interface{} {
	return e.object.Interface()
}

func (e event) GetLifeCycle() lifeCycle {
	return e.lifeCycle
}

type EventListener interface {
	OnPreSave(event Event)
	OnPostSave(event Event)
	OnPostLoad(event Event)
	OnPreDelete(event Event)
	OnPostDelete(event Event)
}

type eventer struct {
	eventListeners map[reflect.Value]EventListener
}

func newEventer() *eventer {
	return &eventer{
		eventListeners: map[reflect.Value]EventListener{}}
}

func (e *eventer) registerEventListener(eventListener EventListener) error {
	e.eventListeners[reflect.ValueOf(eventListener)] = eventListener
	return nil
}

func (e *eventer) disposeEventListener(eventListener EventListener) error {
	delete(e.eventListeners, reflect.ValueOf(eventListener))
	return nil
}
