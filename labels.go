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
	return []string{container.Name()}
}

func getRelationshipType(container reflect.Type) string {
	if taggedRelType := getNamespacedTag(container.FieldByIndex([]int{0}).Tag).get(relTypeTag); len(taggedRelType) > 0 && taggedRelType[0] != emptyString {
		return taggedRelType[0]
	}
	return strings.ToUpper(container.Name())
}

func (f *field) getRelType() string {
	direction := f.getEffectiveDirection()

	relType := f.parent.Type().Name() + defaultRelTypeDelim + elem(f.getStructField().Type).Elem().Name()
	if direction == incoming {
		relType = elem(f.getStructField().Type).Elem().Name() + defaultRelTypeDelim + f.parent.Type().Name()

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
