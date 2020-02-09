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

var metaProperties = map[string]bool{"id": true}

func unloadGraphProperties(g graph, propertyStructFields map[string]*reflect.StructField) {
	if g.getValue().IsValid() {
		for backendName, structField := range propertyStructFields {
			propertyField := &field{
				parent: g.getValue().Elem(),
				name:   structField.Name}
			v := reflect.ValueOf(g.getProperties()[backendName])
			if g.getProperties()[backendName] == nil {
				v = reflect.Zero(structField.Type)
			}
			propertyField.getValue().Set(v)
		}
	}
}

func diffProperties(proposedProperties map[string]interface{}, storedProperties map[string]interface{}) map[string]interface{} {
	var properties = map[string]interface{}{}
	for name, property := range proposedProperties {
		if !reflect.DeepEqual(storedProperties[name], property) {
			properties[name] = proposedProperties[name]
		}
	}
	return properties
}

func getPropertyStructField(t reflect.Type) (map[string]*reflect.StructField, error) {

	var (
		propertyStructFields = map[string]*reflect.StructField{}
		fields               [][]*field
		err                  error
	)

	fields, _ = getFeilds(reflect.New(t).Elem(), propertyFilter)

	for _, field := range fields[0] {
		if field.getStructField().Type.Kind() == reflect.Struct && field.getStructField().Anonymous {
			var promotedStructFields map[string]*reflect.StructField
			if promotedStructFields, err = getPropertyStructField(field.getValue().Type()); err != nil {
				return nil, err
			}
			for promotedPropertyFieldName, structField := range promotedStructFields {
				propertyStructFields[promotedPropertyFieldName] = structField
			}
			continue
		}

		sf := field.getStructField()
		backendName := getBackendPropertyName(field)

		if strings.Contains(backendName, mapPropDelim) {
			return nil, errors.New("Backend property name for field '" + field.getStructField().Name + "' in domain object '" + t.String() + "' can't contain '.'")
		}
		if propertyStructFields[backendName] != nil {
			return nil, errors.New("Backend property name for field '" + field.getStructField().Name + "' in domain object '" + t.String() + "' conflicts with field '" + propertyStructFields[backendName].Name + "'")
		}
		propertyStructFields[backendName] = &sf

	}
	return propertyStructFields, err
}

func getBackendPropertyName(field *field) string {
	propertyName := strings.ToLower(field.getStructField().Name)
	taggedProp := field.tag.get(propertyNameTag)
	if taggedProp != nil && len(taggedProp[0]) != 0 {
		propertyName = strings.ToLower(taggedProp[0])
	}
	return propertyName
}

func driverPropertiesAsStructFieldValues(driverProperties map[string]interface{}, structFields map[string]*reflect.StructField) {
	mappedProperties := map[string]map[string]interface{}{}
	for key, property := range driverProperties {
		if strings.Contains(key, mapPropDelim) {
			mappedPropName := strings.Split(key, mapPropDelim)
			if structFields[mappedPropName[0]] != nil {
				if mappedProperties[mappedPropName[0]] == nil {
					mappedProperties[mappedPropName[0]] = map[string]interface{}{}
				}
				mappedProperties[mappedPropName[0]][mappedPropName[1]] = property
				continue
			}
			continue
		}

		if structFields[key] != nil {
			driverProperties[key] = driverValueAsType(property, structFields[key].Type)
		}
	}

	for backendName, mapp := range mappedProperties {
		structField := structFields[backendName]
		mapElem := structField.Type.Elem()
		mapValue := reflect.MakeMapWithSize(structField.Type, len(mapp))
		for key, value := range mapp {
			mapValue.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(driverValueAsType(value, mapElem)))

		}
		driverProperties[backendName] = mapValue.Interface()
	}

}

func getMapProperties(backendName string, structField *reflect.StructField, v reflect.Value) map[string]interface{} {
	mappedProperties := map[string]interface{}{}
	if mapValue := v.Elem().FieldByName(structField.Name); !mapValue.IsNil() {
		iter := mapValue.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			if k.Type().Kind() == reflect.String {
				mappedProperties[backendName+mapPropDelim+k.String()] = v.Interface()
			}
		}
	}
	return mappedProperties
}
