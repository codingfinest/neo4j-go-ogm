package gogm

import (
	"reflect"
	"sort"
	"strings"
)

func getNodeLabels(container reflect.Type) []string {
	index := []int{0}
	var labels []string

	for field0 := container.FieldByIndex(index); true; {
		if taggedLabels := getNamespacedTag(field0.Tag).get(labelsTag); len(taggedLabels) > 0 {
			labels = append(labels, taggedLabels...)
			break
		}
		labels = append(labels, container.Name())
		if field0.Type == typeOfPublicNode {
			break
		}
		container = field0.Type
		field0 = container.FieldByIndex(index)
	}
	sort.Strings(labels)
	return labels
}

func getThisStructLabels(container reflect.Type) []string {
	if taggedLabels := getNamespacedTag(container.FieldByIndex([]int{0}).Tag).get(labelsTag); len(taggedLabels) > 0 {
		sort.Strings(taggedLabels)
		return taggedLabels
	}
	return []string{strings.ToUpper(container.Name())}
}

func getRelationshipType(container reflect.Type) string {
	if taggedRelType := getNamespacedTag(container.FieldByIndex([]int{0}).Tag).get(relTypeTag); len(taggedRelType) > 0 && taggedRelType[0] != emptyString {
		return taggedRelType[0]
	}
	return strings.ToUpper(container.Name())
}

func (f *field) getRelType() string {
	direction := f.getEffectiveDirection()

	relType := f.parent.Type().Name() + defaultRelTypeDelim + elem(f.getStructField().Type).Name()
	if direction == incoming {
		relType = elem(f.getStructField().Type).Name() + defaultRelTypeDelim + f.parent.Type().Name()

	}

	if relTypes := f.tag.get(relTypeTag); relTypes != nil && len(relTypes[0]) > 0 {
		relType = relTypes[0]
	}
	return relType
}

func getRuntimeLabelsStructFeild(propertyStructFields map[string]*reflect.StructField) *reflect.StructField {
	for _, structField := range propertyStructFields {
		if getNamespacedTag(structField.Tag).get(labelsTag) != nil {
			return structField
		}
	}
	return nil
}
