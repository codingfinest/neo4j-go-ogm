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
	"reflect"
)

type graphFactory struct {
	registry *registry
}

func newGraphFactory(registry *registry) *graphFactory {
	return &graphFactory{
		registry}
}

func (graphFactory graphFactory) get(v reflect.Value, settings map[int]bool) ([]graph, error) {

	var metadata metadata
	var graphs []graph
	var domainObjectType = elem2(v.Type())
	var err error

	if metadata, err = graphFactory.registry.get(domainObjectType); err != nil {
		return nil, err
	}

	if settings == nil {
		settings = map[int]bool{relatedGraph: true, labels: true, properties: true}
	}

	values := []reflect.Value{v.Elem()}

	if v.Kind() == reflect.Ptr && (v.Elem().Kind() == reflect.Array || v.Elem().Kind() == reflect.Slice) {
		values = nil
		objects := v.Elem()
		for i := 0; i < objects.Len(); i++ {
			object := objects.Index(i)
			if !object.IsNil() && object.IsValid() {
				values = append(values, object.Elem().Addr())
			}
		}
	}
	var label string
	switch metadata.(type) {
	case *nodeMetadata:
		for i := 0; i < len(values); i++ {
			node := &node{ID: initialGraphID, Value: &values[i], relationships: map[int64]graph{}}
			if settings[labels] {
				if label, err = metadata.getLabel(values[i]); err != nil {
					return nil, err
				}
				node.setLabel(label)
			}
			if settings[properties] {
				node.setProperties(metadata.getProperties(values[i]))
			}
			graphs = append(graphs, node)
		}
		if len(graphs) == 0 {
			node := &node{ID: initialGraphID, relationships: map[int64]graph{}}
			if settings[labels] {
				if label, err = metadata.getLabel(invalidValue); err != nil {
					return nil, err
				}
				node.setLabel(label)
			}
			graphs = append(graphs, node)
		}
	case *relationshipMetadata:
		for i := 0; i < len(values); i++ {
			relationship := &relationship{ID: initialGraphID, Value: &values[i]}
			if settings[relatedGraph] {
				if relationship.nodes, err = metadata.loadRelatedGraphs(relationship, nil, graphFactory.registry); err != nil {
					return nil, err
				}
			}
			if settings[labels] {
				label, _ = metadata.getLabel(values[i])
				relationship.setLabel(label)
			}
			if settings[properties] {
				relationship.setProperties(metadata.getProperties(values[i]))
			}
			graphs = append(graphs, relationship)
		}
		if len(graphs) == 0 {
			relationship := &relationship{ID: initialGraphID}
			if settings[labels] {
				label, _ = metadata.getLabel(invalidValue)
				relationship.setLabel(label)
			}
			graphs = append(graphs, relationship)
		}
	}

	return graphs, nil
}
