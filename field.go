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

type field struct {
	parent reflect.Value
	name   string
	tag    *tag
}

func (f *field) getValue() reflect.Value {
	return f.parent.FieldByName(f.name)
}

func (f *field) getStructField() reflect.StructField {
	sf, _ := f.parent.Type().FieldByName(f.name)
	return sf
}

func (f *field) isEntity(graphEntityType reflect.Type) bool {
	if f.isIgnored() || f.getStructField().Anonymous || elem(f.getStructField().Type).Kind() != reflect.Struct {
		return false
	}

	fieldType := elem2(f.getStructField().Type)

	if fieldType.Kind() != reflect.Ptr || fieldType.Elem().Kind() != reflect.Struct {
		return false
	}

	if internalGrpahType, err := getInternalGraphType(elem2(f.getStructField().Type).Elem()); internalGrpahType == graphEntityType && err == nil {
		return true
	}

	return false
}

func (f *field) isTagged(s string) bool {
	return f.tag.get(s) != nil
}

func (f *field) isIgnored() bool {
	return f.isTagged(tpc.IgnoreChar)
}

func (f *field) getEffectiveDirection() direction {
	direction := outgoing
	if taggedDirection := f.tag.get(directionTag); taggedDirection != nil && len(taggedDirection[0]) > 0 {
		direction = directionTags[taggedDirection[0]]
	}
	return direction
}
