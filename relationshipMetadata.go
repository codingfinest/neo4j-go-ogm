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

	var relatedGraphField *field
	relatedGraphIndex := startNode
	if relatedGraph == g.getRelatedGraphs()[endNode] {
		relatedGraphIndex = endNode
	}

	if g.getValue().IsValid() {
		relatedGraphField = &field{
			parent: g.getValue().Elem(),
			name:   rm.endpoints[relatedGraphIndex].Name}
	} else {
		relatedGraphs := g.getRelatedGraphs()
		otherNode := relatedGraphs[(relatedGraphIndex+1)%relationshipRelatedGraphLen]
		relDirection := outgoing
		if otherNode == relatedGraphs[endNode] {
			relDirection = incoming
		}
		var (
			metadata metadata
			err      error
		)
		if metadata, err = rm.registry.get(otherNode.getValue().Type()); err != nil {
			return nil, err
		}

		nodeMetadata := metadata.(*nodeMetadata)
		relatedGraphStructLabel := nodeMetadata.getStructLabel(relatedGraph)
		otherNodeGraphStructLabel := nodeMetadata.getStructLabel(otherNode)
		var otherNodeStructField *reflect.StructField
		if relatedGraphStructLabel == otherNodeGraphStructLabel {
			otherNodeStructField = nodeMetadata.getSameEntityRelStructFields(g.getLabel(), relDirection)
		} else {
			fromLabel := otherNodeGraphStructLabel
			toLabel := relatedGraphStructLabel
			if relDirection == incoming {
				fromLabel = relatedGraphStructLabel
				toLabel = otherNodeGraphStructLabel
			}
			otherNodeStructField = nodeMetadata.getDifferentEntityRelStructFields(g.getLabel(), fromLabel, toLabel)
		}

		if otherNodeStructField != nil {
			relatedGraphField = &field{
				parent: otherNode.getValue().Elem(),
				name:   otherNodeStructField.Name}
		}

	}
	return relatedGraphField, nil
}
