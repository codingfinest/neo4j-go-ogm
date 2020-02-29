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
	"sort"
	"strings"
)

type nodeMetadata struct {
	commonMetadata
	thisStructLabel                []string
	disallowedRuntimeLabels        map[string]bool
	runtimeLabelsStructField       *reflect.StructField
	sameEntityRelStructFields      map[string]map[direction]*reflect.StructField
	differentEntityRelStructFields map[string]map[string]map[string]*reflect.StructField
	relationshipAStructFields      []reflect.StructField
	relationshipBStructFields      []reflect.StructField
}

func newNodeMetadata() *nodeMetadata {
	return &nodeMetadata{
		disallowedRuntimeLabels:        map[string]bool{},
		sameEntityRelStructFields:      map[string]map[direction]*reflect.StructField{},
		differentEntityRelStructFields: map[string]map[string]map[string]*reflect.StructField{}}
}

func (nm *nodeMetadata) getLabel(v reflect.Value) (string, error) {
	if !v.IsValid() || !v.Elem().IsValid() || v.IsZero() || v.Elem().IsZero() || nm.runtimeLabelsStructField == nil {
		return nm.structLabel, nil
	}

	var runtimeLabels []string
	for _, label := range v.Elem().FieldByName(nm.runtimeLabelsStructField.Name).Interface().([]string) {
		if nm.disallowedRuntimeLabels[label] {
			return emptyString, errors.New("Runtime label '" + label + "' is't allowed. Either the the object of type " + v.Elem().Type().String() + " has this label or a node in one of its relationship does. Use another label.")
		}
		runtimeLabels = append(runtimeLabels, label)
	}

	runtimeLabels = append(runtimeLabels, strings.Split(nm.structLabel, labelsDelim)...)
	sort.Strings(runtimeLabels)
	return strings.Join(runtimeLabels, labelsDelim), nil
}

func (nm *nodeMetadata) loadRelatedGraphs(g graph, ID func(graph), registry *registry) (map[int64]graph, error) {

	relatedGraphs := map[int64]graph{}
	value := *g.getValue()
	relationshipAFieldExtractor := extractRelationshipTypeA(g.(*node), ID, registry)
	relationshipBFieldExtractor := extractRelationshipTypeB(g.(*node), ID, registry)

	for _, relationshipAStructField := range nm.relationshipAStructFields {
		sf, _ := value.Elem().Type().FieldByName(relationshipAStructField.Name)
		f := &field{
			parent: value.Elem(),
			name:   relationshipAStructField.Name,
			tag:    getNamespacedTag(sf.Tag)}

		var relationships []*relationship
		var err error
		if relationships, err = relationshipAFieldExtractor(f); err != nil {
			return nil, err
		}
		for _, relatedGraph := range relationships {
			relatedGraphs[relatedGraph.getID()] = relatedGraph
		}
	}

	for _, relationshipBStructField := range nm.relationshipBStructFields {
		sf, _ := value.Elem().Type().FieldByName(relationshipBStructField.Name)
		f := &field{
			parent: value.Elem(),
			name:   relationshipBStructField.Name,
			tag:    getNamespacedTag(sf.Tag)}

		var relationships []*relationship
		var err error
		if relationships, err = relationshipBFieldExtractor(f); err != nil {
			return nil, err
		}

		for _, relatedGraph := range relationships {
			relatedGraphs[relatedGraph.getID()] = relatedGraph
		}
	}

	return relatedGraphs, nil
}

func (nm *nodeMetadata) getGraphField(g graph, relatedGraph graph) (*field, error) {
	var relatedGraphField *field
	var relatedGraphFieldStructField *reflect.StructField
	fromNode := relatedGraph.getRelatedGraphs()[startNode]
	toNode := relatedGraph.getRelatedGraphs()[endNode]

	fromStructLabel := nm.filterStructLabel(fromNode)
	toStructLabel := nm.filterStructLabel(toNode)

	if fromStructLabel == toStructLabel {
		direction := outgoing
		if g == toNode {
			direction = incoming
		}

		if nm.sameEntityRelStructFields[relatedGraph.getLabel()] != nil && nm.sameEntityRelStructFields[relatedGraph.getLabel()][direction] != nil {
			relatedGraphFieldStructField = nm.sameEntityRelStructFields[relatedGraph.getLabel()][direction]
			relatedGraphField = &field{
				parent: g.getValue().Elem(),
				name:   relatedGraphFieldStructField.Name}
		}
	} else {
		if nm.differentEntityRelStructFields[relatedGraph.getLabel()] != nil && nm.differentEntityRelStructFields[relatedGraph.getLabel()][fromStructLabel] != nil && nm.differentEntityRelStructFields[relatedGraph.getLabel()][fromStructLabel][toStructLabel] != nil {
			relatedGraphFieldStructField = nm.differentEntityRelStructFields[relatedGraph.getLabel()][fromStructLabel][toStructLabel]
			relatedGraphField = &field{
				parent: g.getValue().Elem(),
				name:   relatedGraphFieldStructField.Name}
		}
	}

	return relatedGraphField, nil
}

func (nm *nodeMetadata) filterStructLabel(g graph) string {
	var labels []string
	for _, label := range strings.Split(g.getLabel(), labelsDelim) {
		if nm.disallowedRuntimeLabels[label] {
			labels = append(labels, label)
		}
	}
	sort.Strings(labels)
	return strings.Join(labels, labelsDelim)
}

func (nm *nodeMetadata) setSameEntityRelStructFields(typeOfObject reflect.Type, relType string, relDirection direction, relationshipStructField *reflect.StructField) error {
	if nm.sameEntityRelStructFields[relType] == nil {
		nm.sameEntityRelStructFields[relType] = map[direction]*reflect.StructField{}
	}
	if nm.sameEntityRelStructFields[relType][relDirection] != nil && nm.sameEntityRelStructFields[relType][relDirection].Name != relationshipStructField.Name {
		return errors.New("Ambiguous relationship detected between field '" + relationshipStructField.Name + "' and field '" + nm.sameEntityRelStructFields[relType][relDirection].Name + "' in domain object '" + typeOfObject.Elem().String() + "'")
	}

	if nm.sameEntityRelStructFields[relType][relDirection] == nil {
		nm.sameEntityRelStructFields[relType][relDirection] = relationshipStructField
	}
	return nil
}

func (nm *nodeMetadata) setDifferentEntityRelStructFields(typeOfObject reflect.Type, relType string, fromLabel string, toLabel string, relationshipStructField *reflect.StructField) error {
	if nm.differentEntityRelStructFields[relType] == nil {
		nm.differentEntityRelStructFields[relType] = map[string]map[string]*reflect.StructField{}
	}
	if nm.differentEntityRelStructFields[relType][fromLabel] == nil {
		nm.differentEntityRelStructFields[relType][fromLabel] = map[string]*reflect.StructField{}
	}

	if nm.differentEntityRelStructFields[relType][fromLabel][toLabel] != nil && nm.differentEntityRelStructFields[relType][fromLabel][toLabel].Name != relationshipStructField.Name {
		return errors.New("Ambiguous relationship detected between field '" + relationshipStructField.Name + "' and field '" + nm.differentEntityRelStructFields[relType][fromLabel][toLabel].Name + "' in domain object '" + typeOfObject.Elem().String() + "'")
	}

	nm.differentEntityRelStructFields[relType][fromLabel][toLabel] = relationshipStructField

	return nil
}

func (nm *nodeMetadata) blacklistLabels(labels []string) {
	for _, label := range labels {
		nm.disallowedRuntimeLabels[label] = true
	}
}

func extractRelationshipTypeA(n *node, ID func(graph), registry *registry) func(*field) ([]*relationship, error) {
	return func(f *field) ([]*relationship, error) {

		var (
			relationships []*relationship
			value         = invalidValue
			relType       = f.getRelType()
			values        = getEntitiesFromField(f)
			direction     = f.getEffectiveDirection()
		)

		for i := 0; i < len(values); i++ {
			var relationshipA *relationship
			var node = &node{ID: initialGraphID, Value: &values[i], relationships: make(map[int64]graph)}

			var metadata metadata
			var label string
			var err error
			if metadata, err = registry.get(values[i].Type()); err != nil {
				return nil, err
			}
			if label, err = metadata.getLabel(values[i]); err != nil {
				return nil, err
			}
			node.setLabel(label)
			node.setProperties(metadata.getProperties(values[i]))

			if direction == incoming {
				relationshipA = &relationship{
					Value:      &value,
					nodes:      map[int64]graph{startNode: node, endNode: n},
					relType:    relType,
					properties: make(map[string]interface{})}
			}

			if direction == outgoing || direction == undirected {
				relationshipA = &relationship{
					Value:      &value,
					nodes:      map[int64]graph{startNode: n, endNode: node},
					relType:    relType,
					properties: make(map[string]interface{})}
			}
			ID(relationshipA.nodes[startNode])
			ID(relationshipA.nodes[endNode])
			ID(relationshipA)

			relationships = append(relationships, relationshipA)
		}

		return relationships, nil
	}
}

func extractRelationshipTypeB(n *node, ID func(graph), registry *registry) func(*field) ([]*relationship, error) {
	return func(f *field) ([]*relationship, error) {
		var (
			metadata metadata
			err      error
		)
		if metadata, err = registry.get(elem2(f.getValue().Type())); err != nil {
			return nil, err
		}

		var relationships []*relationship
		values := getEntitiesFromField(f)
		for i := 0; i < len(values); i++ {
			var label string
			relationship := &relationship{ID: initialGraphID, Value: &values[i]}
			if relationship.nodes, err = metadata.loadRelatedGraphs(relationship, ID, registry); err != nil {
				return nil, err
			}
			if label, err = metadata.getLabel(values[i]); err != nil {
				return nil, err
			}
			relationship.setLabel(label)
			relationship.setProperties(metadata.getProperties(values[i]))
			ID(relationship)
			relationships = append(relationships, relationship)
		}
		return relationships, nil
	}
}
