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
