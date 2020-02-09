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

import (
	"errors"
	"reflect"
)

func notifyPreSaveGraph(g graph, eventer eventer, registry *registry) error {

	if g.getValue().IsValid() {
		for _, eventListener := range eventer.eventListeners {
			eventListener.OnPreSave(event{g.getValue(), -1})
		}

		var (
			metadata metadata
			err      error
		)
		if metadata, err = registry.get(g.getValue().Type()); err != nil {
			return err
		}
		var label string
		if label, err = metadata.getLabel(*g.getValue()); err != nil {
			return err
		}
		g.setLabel(label)
		g.setProperties(metadata.getProperties(*g.getValue()))

		customIDName, customIDValue := metadata.getCustomID(*g.getValue())
		if customIDName != emptyString && customIDValue.Type().Kind() == reflect.Ptr && customIDValue.IsNil() {
			return errors.New("Custom ID cannot be nil in " + g.getValue().Type().String())
		}

	}

	return nil
}

func notifyPostSave(eventer eventer, g graph, lifeCycle lifeCycle) error {
	if g == nil {
		return nil
	}
	if g.getValue().IsValid() {
		for _, eventListener := range eventer.eventListeners {
			eventListener.OnPostSave(event{g.getValue(), lifeCycle})
		}
	}
	return nil
}

func notifyPostDelete(eventer eventer, g graph, lifeCycle lifeCycle) error {
	if g == nil {
		return nil
	}

	g.setID(initialGraphID)
	if g.getValue().IsValid() {
		id := initialGraphID
		unloadGraphID(g, &id)
	}
	if g.getValue().IsValid() {
		for _, eventListener := range eventer.eventListeners {
			eventListener.OnPostDelete(event{g.getValue(), lifeCycle})
		}
	}

	return nil
}
