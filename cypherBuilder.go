package gogm

import (
	"reflect"
	"strings"
)

type graphQueryBuilder interface {
	getGraph() graph
	getCreate() (string, string, map[string]interface{}, map[string]graph)
	getMatch() (string, map[string]interface{}, map[string]graph)
	getSet() (string, map[string]interface{})
	getDelete() (string, map[string]interface{}, map[string]graph)
	getLoadAll(IDs interface{}, lo *LoadOptions) (string, map[string]interface{})
	getDeleteAll() (string, map[string]interface{})
	getCountEntitiesOfType() (string, map[string]interface{})

	isGraphDirty() bool
	getRemovedGraphs() (map[int64]graph, map[int64]graph)
}

func newCypherBuilder(g graph, registry *registry, store store) (graphQueryBuilder, error) {
	var (
		qGraphBuilder graphQueryBuilder
		stored        graph
		err           error
	)
	if store != nil {
		stored = store.get(g)
	}
	switch v := g.(type) {
	case *node:
		if qGraphBuilder, err = newNodeCypherBuilder(v, registry, stored); err != nil {
			return nil, err
		}
	case *relationship:
		qGraphBuilder = newRelationshipCypherBuilder(v, registry, stored)
	}
	return qGraphBuilder, nil
}

func getCreateSchemaStatement(metadata metadata) []string {

	var indexes []string
	var unique = map[string]bool{}
	var statements []string
	var objectMetadata *nodeMetadata
	var ok bool

	if typeOfNodeMetadata != reflect.TypeOf(metadata) {
		return statements
	}

	if objectMetadata, ok = metadata.(*nodeMetadata); !ok {
		return nil
	}

	for name, structField := range metadata.getPropertyStructFields() {
		namespaceTag := getNamespacedTag(structField.Tag)
		if (len(namespaceTag.get(uniqueTag)) > 0 || len(namespaceTag.get(customIDTag)) > 0) && !unique[name] {
			unique[name] = true
		}
	}

	for name, structField := range metadata.getPropertyStructFields() {
		namespaceTag := getNamespacedTag(structField.Tag)
		if len(namespaceTag.get(indexTag)) > 0 && !unique[name] {
			indexes = append(indexes, name)
		}
	}

	for name := range unique {
		for _, label := range objectMetadata.thisStructLabel {
			statements = append(statements, `CREATE CONSTRAINT ON (a:`+label+`) ASSERT a.`+name+` IS UNIQUE`)
		}
	}

	compositeIndexes := strings.Join(indexes, indexDelim)
	if compositeIndexes != emptyString {
		for _, label := range objectMetadata.thisStructLabel {
			statements = append(statements, `CREATE INDEX ON :`+label+`(`+compositeIndexes+`)`)

		}
	}
	return statements
}
