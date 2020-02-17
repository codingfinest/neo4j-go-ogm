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

type fieldFilter func(*field) bool

func getFeilds(v reflect.Value, fieldFilters ...fieldFilter) ([][]*field, error) {

	var fields = make([][]*field, len(fieldFilters))
	vType := v.Type()

	for fieldIdx := 0; fieldIdx < vType.NumField(); fieldIdx++ {
		field := &field{
			parent: v,
			name:   vType.Field(fieldIdx).Name,
			tag:    getNamespacedTag(vType.Field(fieldIdx).Tag),
		}
		for filterIdx, fieldFilter := range fieldFilters {
			if match := fieldFilter(field); match {
				fields[filterIdx] = append(fields[filterIdx], field)
			}
		}
	}
	return fields, nil
}

func getEntitiesFromField(f *field) []reflect.Value {
	values := []reflect.Value{}
	kind := f.getStructField().Type.Kind()
	//TODO  rm support for reflect.Array
	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < f.getValue().Len(); i++ {
			if f.getValue().Index(i).IsNil() || !f.getValue().Index(i).Elem().IsValid() {
				continue
			}
			values = append(values, f.getValue().Index(i).Elem().Addr())
		}
	} else if !f.getValue().IsNil() && f.getValue().Elem().IsValid() {
		values = append(values, f.getValue().Elem().Addr())
	}
	return values
}

func addDomainObject(f *field, value reflect.Value) {
	kind := f.getStructField().Type.Kind()
	if kind == reflect.Slice {
		f.getValue().Set(reflect.Append(f.getValue(), value))
	} else {
		f.getValue().Set(value)
	}
}

//Filter
func propertyFilter(f *field) bool {
	return !f.isIgnored() && !f.isEntity(typeOfPrivateNode) && !f.isEntity(typeOfPrivateRelationship) && f.getValue().CanInterface()
}

func isEmbeddedFieldFilter(f *field) bool {
	return f.getStructField().Type.Kind() == reflect.Struct && f.getStructField().Anonymous
}

func isRelationshipFieldFilter(_type reflect.Type) fieldFilter {
	return func(f *field) bool {
		fType := f.getStructField().Type
		kind := fType.Kind()
		return f.isEntity(_type) &&
			(((kind == reflect.Slice || kind == reflect.Array) && fType.Elem().Kind() == reflect.Ptr && fType.Elem().Elem().Kind() == reflect.Struct) ||
				(kind == reflect.Ptr && fType.Elem().Kind() == reflect.Struct))
	}
}

func isRelationshipEndPointFieldFilter(endpoint string) fieldFilter {
	return func(f *field) bool {
		return !f.isIgnored() && !f.getStructField().Anonymous &&
			f.getValue().Kind() == reflect.Ptr && elem(f.getStructField().Type).Kind() == reflect.Struct &&
			f.isTagged(endpoint)
	}
}
