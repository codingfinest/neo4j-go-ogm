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
	"sort"
	"strings"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type queryer struct {
	cypherExecuter *cypherExecuter
	graphFactory   graphFactory
	registry       *registry
}

func newQueryer(cypherExecutor *cypherExecuter, graphFactory graphFactory, registry *registry) *queryer {
	return &queryer{cypherExecutor, graphFactory, registry}
}

func (q *queryer) queryForObject(object interface{}, cypher string, parameters map[string]interface{}) error {
	var (
		err      error
		values   reflect.Value
		metadata metadata
		records  []neo4j.Record
	)
	//Object: **DomainObject
	domainObjectType := reflect.TypeOf(object).Elem()
	if metadata, err = q.registry.get(domainObjectType); err != nil {
		return err
	}
	var label string
	if label, err = metadata.getLabel(invalidValue); err != nil {
		return err
	}
	if records, err = neo4j.Collect(q.cypherExecuter.exec(cypher, parameters)); err != nil {
		return err
	}

	if len(records) == 1 {
		if values, err = q.getObjectsFromRecords(domainObjectType, metadata, label, []neo4j.Record{records[0]}); err != nil {
			return err
		}

		reflect.ValueOf(object).Elem().Set(values.Index(0))
	}

	if len(records) > 1 {
		return errors.New("result contains more than one record")
	}

	return nil

}

func (q *queryer) queryForObjects(objects interface{}, cypher string, parameters map[string]interface{}) error {

	var (
		err      error
		records  []neo4j.Record
		values   reflect.Value
		metadata metadata
	)
	//Object type is *[]*<DomaoinObject>
	domainObjectType := reflect.TypeOf(objects).Elem().Elem()
	if metadata, err = q.registry.get(domainObjectType); err != nil {
		return err
	}

	var label string
	if label, err = metadata.getLabel(invalidValue); err != nil {
		return err
	}

	if records, err = neo4j.Collect(q.cypherExecuter.exec(cypher, parameters)); err != nil {
		return err
	}

	if values, err = q.getObjectsFromRecords(domainObjectType, metadata, label, records); err != nil {
		return err
	}
	reflect.ValueOf(objects).Elem().Set(values)
	return nil
}

func (q *queryer) query(cypher string, parameters map[string]interface{}) (neo4j.Result, error) {
	return q.cypherExecuter.exec(cypher, parameters)
}

func (q *queryer) getObjectsFromRecords(domainObjectType reflect.Type, metadata metadata, label string, records []neo4j.Record) (reflect.Value, error) {

	var (
		err                     error
		g                       graph
		entityLabel             string
		internalGraphEntityType reflect.Type
	)

	if internalGraphEntityType, err = getInternalGraphType(domainObjectType.Elem()); err != nil {
		return invalidValue, err
	}

	sliceOfPtrToObjs := reflect.MakeSlice(reflect.SliceOf(domainObjectType), 0, 0)
	ptrToObjs := reflect.New(sliceOfPtrToObjs.Type())

	for _, record := range records {
		column0 := record.GetByIndex(0)
		newPtrToDomainObject := reflect.New(domainObjectType.Elem())

		if neo4jNode, isNeo4jNode := column0.(neo4j.Node); isNeo4jNode == true {

			if internalGraphEntityType != typeOfPrivateNode {
				return invalidValue, errors.New("Expecting a Relationship, but got a Node from the query response")
			}
			nodeMetadata := metadata.(*nodeMetadata)
			labels := neo4jNode.Labels()
			sort.Strings(labels)
			g = &node{
				ID:         neo4jNode.Id(),
				properties: neo4jNode.Props(),
				label:      strings.Join(labels, labelsDelim)}
			g.getProperties()[idPropertyName] = neo4jNode.Id()

			entityLabel = nodeMetadata.getStructLabel(g)
		}

		if neo4jRelationship, isNeo4jReleationship := column0.(neo4j.Relationship); isNeo4jReleationship == true {
			if internalGraphEntityType != typeOfPrivateRelationship {
				return invalidValue, errors.New("Unexpected graph type. Expecting a Node, but got a Relationship from the query response")
			}
			g = &relationship{
				ID:         neo4jRelationship.Id(),
				properties: neo4jRelationship.Props(),
				relType:    neo4jRelationship.Type()}
			g.getProperties()[idPropertyName] = neo4jRelationship.Id()
			entityLabel = neo4jRelationship.Type()
		}
		g.setValue(&newPtrToDomainObject)
		g.setLabel(label)

		if label != entityLabel {
			return invalidValue, errors.New("label '" + label + "' from `" + domainObjectType.String() + "` don't match with label `" + entityLabel + "` from query result")
		}

		ptrToObjs.Elem().Set(reflect.Append(ptrToObjs.Elem(), newPtrToDomainObject))
		driverPropertiesAsStructFieldValues(g.getProperties(), metadata.getPropertyStructFields())
		unloadGraphProperties(g, metadata.getPropertyStructFields())
	}

	return ptrToObjs.Elem(), nil
}

func (q *queryer) countEntitiesOfType(object interface{}) (int64, error) {

	var (
		value         = reflect.ValueOf(object)
		record        neo4j.Record
		count         int64
		cypherBuilder graphQueryBuilder
		graphs        []graph
		cypher        string
		parameters    map[string]interface{}
		err           error
	)

	//object: **DomainObject
	if graphs, err = q.graphFactory.get(reflect.New(value.Elem().Type()), map[int]bool{labels: true}); err != nil {
		return -1, err
	}

	if cypherBuilder, err = newCypherBuilder(graphs[0], q.registry, nil); err != nil {
		return -1, err
	}
	cypher, parameters = cypherBuilder.getCountEntitiesOfType()

	if cypher != emptyString {
		if record, err = neo4j.Single(q.cypherExecuter.exec(cypher, parameters)); err != nil {
			return -1, err
		}
		if record != nil {
			count = record.GetByIndex(0).(int64)
		}
	}

	return count, nil
}

func (q *queryer) count(cypher string, parameters map[string]interface{}) (int64, error) {
	var (
		record neo4j.Record
		err    error
	)
	if record, err = neo4j.Single(q.cypherExecuter.exec(cypher, parameters)); err != nil {
		return -1, err
	}
	return record.GetByIndex(0).(int64), nil
}
