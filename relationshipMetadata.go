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

const relationshipRelatedGraphLen = 2

type relationshipMetadata struct {
	commonMetadata
	endpoints map[int64]reflect.StructField
}

func newRelationshipMetadata() *relationshipMetadata {
	return &relationshipMetadata{
		endpoints: map[int64]reflect.StructField{}}
}

func (rm *relationshipMetadata) loadRelatedGraphs(g graph, ID func(graph), registry *registry) (map[int64]graph, error) {

	if g.getValue().IsZero() || len(g.getRelatedGraphs()) == relationshipRelatedGraphLen {
		return g.getRelatedGraphs(), nil
	}

	value := (*g.getValue()).Elem()
	relatedGraphs := map[int64]graph{}

	startValue := value.FieldByName(rm.endpoints[startNode].Name)
	endValue := value.FieldByName(rm.endpoints[endNode].Name)

	if startValue.IsNil() {
		return nil, errors.New("start node for relationship is nil. Expected a non-nil start node")
	}
	if endValue.IsNil() {
		return nil, errors.New("end node for relationship is nil. Expected a non-nil end node")
	}
	v1 := startValue.Elem().Addr()
	v2 := endValue.Elem().Addr()

	relatedGraphs[startNode] = &node{Value: &v1, relationships: map[int64]graph{}}
	relatedGraphs[endNode] = &node{Value: &v2, relationships: map[int64]graph{}}

	var (
		metadata metadata
		err      error
	)
	if metadata, err = registry.get(v1.Type()); err != nil {
		return nil, err
	}
	var label string
	if label, err = metadata.getLabel(v1); err != nil {
		return nil, err
	}
	relatedGraphs[startNode].setLabel(label)
	relatedGraphs[startNode].setProperties(metadata.getProperties(v1))

	if metadata, err = registry.get(v2.Type()); err != nil {
		return nil, err
	}
	if label, err = metadata.getLabel(v2); err != nil {
		return nil, err
	}
	relatedGraphs[endNode].setLabel(label)
	relatedGraphs[endNode].setProperties(metadata.getProperties(v2))

	if ID != nil {
		ID(relatedGraphs[startNode])
		ID(relatedGraphs[endNode])
	}
	return relatedGraphs, nil
}

func (rm *relationshipMetadata) getLabel(v reflect.Value) (string, error) {
	return rm.structLabel, nil
}

func (rm *relationshipMetadata) getGraphField(g graph, relatedGraph graph) (*field, error) {

	relatedGraphIndex := startNode
	if relatedGraph == g.getRelatedGraphs()[endNode] {
		relatedGraphIndex = endNode
	}

	return &field{
		parent: g.getValue().Elem(),
		name:   rm.endpoints[relatedGraphIndex].Name}, nil
}
