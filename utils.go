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

//TODO use first field approach. //TODO Doc: To get the embedded node fiedl
func getInternalGraphType(t reflect.Type) (reflect.Type, error) {

	var (
		internalType   reflect.Type
		embeddedFields [][]*field
		err            error
	)

	if internalType = getInternalType(t); internalType != nil {
		return internalType, err
	}

	if embeddedFields, err = getFeilds(reflect.New(t).Elem(), isEmbeddedFieldFilter); err != nil {
		return nil, err
	}

	if len(embeddedFields[0]) == 0 {
		return nil, errors.New("No embedded field found. The internal graph type of " + t.String() + " can't be determined")
	}

	for _, embeddedField := range embeddedFields[0] {
		if internalType, err = getInternalGraphType(embeddedField.getStructField().Type); internalType != nil || err != nil {
			return internalType, err
		}
	}
	return nil, errors.New("The internal graph type of " + t.String() + " can't be determined")
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
