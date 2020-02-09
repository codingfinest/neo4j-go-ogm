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
	"strings"
)

type internalIDGenerator struct {
	currentID int64
}

func (generator *internalIDGenerator) new() int64 {
	generator.currentID--
	return generator.currentID
}

func getCustomIDBackendName(structFields map[string]*reflect.StructField) (string, error) {
	for backendName, structField := range structFields {
		if len(getNamespacedTag(structField.Tag).get(customIDTag)) > 0 {

			switch structField.Type.Kind() {
			case reflect.Bool,
				reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Float32, reflect.Float64, reflect.String:
				return backendName, nil
			default:
				return emptyString, errors.New("Invalid custom ID type")

			}

		}
	}
	return emptyString, nil
}

func unloadGraphID(g graph, id *int64) {
	if id == nil {
		g.setID(initialGraphID)
		delete(g.getProperties(), idPropertyName)
	} else {
		g.setID(*id)
		g.getProperties()[idPropertyName] = id
	}
	if g.getValue() != nil && g.getValue().IsValid() {
		*getIDAddr(g) = id
	}
}

func getIDAddr(g graph) **int64 {
	if g.getValue() != nil && g.getValue().IsValid() {
		internalIDField := g.getValue().Elem().FieldByName(strings.ToUpper(idPropertyName))
		return internalIDField.Addr().Interface().(**int64)
	}
	return nil
}

func getIDer(idGenerator *internalIDGenerator, store store) func(graph) {
	var ID func(g graph)
	ID = func(g graph) {
		if g.getValue().IsValid() {
			internalIDAddr := getIDAddr(g)
			if *internalIDAddr == nil && idGenerator != nil {
				id := idGenerator.new()
				*internalIDAddr = &id
				g.getProperties()[idPropertyName] = id
			}
			if *internalIDAddr != nil {
				g.setID(**internalIDAddr)
			}
		} else {
			for _, relatedGraph := range g.getRelatedGraphs() {
				ID(relatedGraph)
			}
			relationshipA := store.get(g)
			if relationshipA == nil {
				g.setID(idGenerator.new())
			} else {
				g.setID(relationshipA.getID())
			}
		}
	}
	return ID
}
