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
	"strings"
)

var (
	relTypeTag      = "reltype"
	startNodeTag    = "startNode"
	endNodeTag      = "endNode"
	labelsTag       = "label"
	directionTag    = "direction"
	customIDTag     = "id"
	propertyNameTag = "name"
	uniqueTag       = "unique"
	indexTag        = "index"
)

var (
	tpc = tagParserConf{
		TagKey:     "gogm",
		Delim:      ",",
		AssignOp:   ":",
		IgnoreChar: "-",
		MultiPropsAllowed: map[string]bool{
			"label": true,
		},
	}
)

type tagParserConf struct {
	TagKey            string
	IgnoreChar        string
	Delim             string
	AssignOp          string
	MultiPropsAllowed map[string]bool
}

type tag struct {
	keyval       string
	mappedKeyval map[string][]string
}

func getNamespacedTag(fieldTag reflect.StructTag) *tag {
	tag := &tag{
		keyval: emptyString}
	if namespaceTag, ok := fieldTag.Lookup(tpc.TagKey); ok {
		tag.keyval = namespaceTag
	}
	return tag
}

func (t *tag) get(key string) []string {
	if t.mappedKeyval != nil {
		return t.mappedKeyval[key]
	}
	t.mappedKeyval = map[string][]string{}
	keyVals := strings.Split(t.keyval, tpc.Delim)
	for _, v := range keyVals {
		keyVal := strings.Split(v, tpc.AssignOp)
		if len(keyVal) == 2 {
			key := strings.Trim(keyVal[0], spaceString)
			value := strings.Trim(keyVal[1], spaceString)
			t.mappedKeyval[key] = append(t.mappedKeyval[key], value)
		}

		if len(keyVal) == 1 {
			key := strings.Trim(keyVal[0], spaceString)
			t.mappedKeyval[key] = append(t.mappedKeyval[key], emptyString)
		}
	}
	return t.mappedKeyval[key]
}
