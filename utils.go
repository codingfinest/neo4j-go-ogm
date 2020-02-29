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

func getInternalType(t reflect.Type) reflect.Type {
	var internalType reflect.Type
	switch t {
	case typeOfPublicNode:
		internalType = typeOfPrivateNode
	case typeOfPublicRelationship:
		internalType = typeOfPrivateRelationship
	}
	return internalType
}

func getInternalGraphType(container reflect.Type) reflect.Type {

	index := []int{0}
	f := container.FieldByIndex(index)

	for field0 := &f; field0 != nil; {
		if internalType := getInternalType(field0.Type); internalType != nil {
			return internalType
		}

		if field0.Type.Kind() == reflect.Struct && field0.Type.NumField() > 0 && field0.Type.FieldByIndex(index).Anonymous {
			f = field0.Type.FieldByIndex(index)
			field0 = &f
		} else {
			field0 = nil
		}
	}

	return nil
}

func flattenParamters(parameters []map[string]interface{}) map[string]interface{} {
	var flattenedParamters = map[string]interface{}{}
	for _, parameter := range parameters {
		for key, value := range parameter {
			flattenedParamters[key] = value
		}
	}
	return flattenedParamters
}

func getCyhperFromClauses(cypherClauses map[clause][]string) string {
	cypher := ``
	claused := map[string]bool{}
	for _, clauseGroup := range clauseGroups {
		for _, clause := range cypherClauses[clauseGroup] {
			if !claused[clause] {
				cypher += clause
				claused[clause] = true
			}
		}
	}
	return cypher
}

func indexOfString(slice []string, target string) int {
	var index = -1
	for i, s := range slice {
		if s == target {
			index = i
			break
		}
	}
	return index
}

func removeStringAt(slice []string, index int) []string {
	slice[index] = slice[len(slice)-1]
	slice[len(slice)-1] = ""
	return slice[:len(slice)-1]
}

//Return the simple domain object type of a simple or compound
//type t. Example of a simple type is *DomainObject. Example of
//a compound type is *[]*DomainObject.
func elem(t reflect.Type) reflect.Type {
	if t.Kind() != reflect.Array && t.Kind() != reflect.Slice && t.Kind() != reflect.Ptr {
		return t
	} else if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		return t
	}
	return elem(t.Elem())
}
