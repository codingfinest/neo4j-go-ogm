// MIT License
//
// Copyright (c) 2020 codingfinest
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
