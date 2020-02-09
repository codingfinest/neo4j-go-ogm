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

import "reflect"

func canSlice(k reflect.Kind) bool {
	return k == reflect.Array || k == reflect.Slice || k == reflect.String
}

func canTypeElem(k reflect.Kind) bool {
	return k == reflect.Array || k == reflect.Slice || k == reflect.Ptr || k == reflect.Map || k == reflect.Chan
}

func elem(_type reflect.Type) reflect.Type {
	var (
		typeToElem = _type
	)

	if typeToElem.Kind() == reflect.Struct {
		return typeToElem
	} else if (canSlice(typeToElem.Kind()) && canTypeElem(typeToElem.Kind())) || typeToElem.Kind() == reflect.Ptr {
		return elem(typeToElem.Elem())
	}

	return typeToElem
}

//Return the simple domain object type of a simple or compound
//type t. Example of a simple type is *DomainObject. Example of
//a compound type is *[]*DomainObject.
func elem2(t reflect.Type) reflect.Type {

	if !canTypeElem(t.Kind()) {
		return t
	} else if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		return t
	}
	return elem2(t.Elem())
}

func embedsType(container reflect.Type, t reflect.Type) bool {
	containerType := elem(container)
	for i := 0; i < containerType.NumField(); i++ {
		childField := &field{
			parent: reflect.New(containerType).Elem(),
			name:   containerType.Field(i).Name,
			// index:  i,
			tag: getNamespacedTag(containerType.Field(i).Tag),
		}
		if elem(childField.getStructField().Type) == t {
			return true
		}
	}
	return false
}
