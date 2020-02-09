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
	claused := make(map[string]bool)
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
