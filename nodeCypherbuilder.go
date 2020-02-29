package gogm

import (
	"strconv"
)

type nodeQueryBuilder struct {
	n                              *node
	registry                       *registry
	deltaProperties                map[string]interface{}
	isLabelsDirty                  bool
	removedRelationships           map[int64]graph
	removedRelationshipsOtherNodes map[int64]graph
}

func (nqb nodeQueryBuilder) getGraph() graph {
	return nqb.n
}

func (nqb nodeQueryBuilder) getGraphField(relatedGraph graph) (*field, error) {
	var (
		metadata metadata
		field    *field
		err      error
	)
	if metadata, err = nqb.registry.get(nqb.n.getValue().Type()); err != nil {
		return nil, err
	}
	if field, err = metadata.getGraphField(nqb.n, relatedGraph); err != nil {
		return nil, err
	}
	return field, nil
}

func (nqb nodeQueryBuilder) diffNodeRelatedGraphs(ref graph) (addedGraphs map[int64]graph, removedGraphs map[int64]graph, _err error) {
	if ref == nil {
		return nil, nil, nil
	}

	var (
		field *field
		err   error
	)
	removedRelationships := map[int64]graph{}
	removedRelationshipsOtherNodes := map[int64]graph{}
	refNode := ref.(*node)

	for internalID, storedRelationship := range refNode.relationships {
		if nqb.n.relationships[internalID] == nil {
			if field, err = nqb.getGraphField(storedRelationship); err != nil {
				return nil, nil, err
			}

			if field == nil {
				//The node n, doesn't have the field for this relationship. Hence,
				//the relationship wasn't removed by a nil assignment to the field
				continue
			}

			otherNode := storedRelationship.getRelatedGraphs()[startNode]
			if nqb.n.getID() == otherNode.getID() {
				otherNode = storedRelationship.getRelatedGraphs()[endNode]
			}

			removedRelationships[internalID] = storedRelationship
			removedRelationshipsOtherNodes[internalID] = otherNode
		}
	}

	return removedRelationships, removedRelationshipsOtherNodes, nil
}

func newNodeCypherBuilder(n *node, registry *registry, stored graph) (*nodeQueryBuilder, error) {

	var err error
	nqb := &nodeQueryBuilder{
		n:               n,
		registry:        registry,
		deltaProperties: n.getProperties(),
		isLabelsDirty:   true}

	if stored != nil {
		nqb.deltaProperties = diffProperties(n.getProperties(), stored.getProperties())
		nqb.isLabelsDirty = n.getLabel() != stored.getLabel()
		nqb.removedRelationships, nqb.removedRelationshipsOtherNodes, err = nqb.diffNodeRelatedGraphs(stored)
		if err != nil {
			return nil, err
		}
	}

	return nqb, nil
}

func (nqb nodeQueryBuilder) getRemovedGraphs() (map[int64]graph, map[int64]graph) {
	return nqb.removedRelationships, nqb.removedRelationshipsOtherNodes
}

func (nqb nodeQueryBuilder) isGraphDirty() bool {
	return nqb.n.getID() < 0 || len(nqb.deltaProperties) > 0 || len(nqb.removedRelationships) > 0 || nqb.isLabelsDirty
}

func (nqb nodeQueryBuilder) getCreate() (string, string, map[string]interface{}, map[string]graph) {
	create := `CREATE (` + nqb.n.getSignature() + `)
	`
	return create, emptyString, nil, nil
}

func (nqb nodeQueryBuilder) getMatch() (string, map[string]interface{}, map[string]graph) {
	var (
		nSign                                       = nqb.n.getSignature()
		metadata, _                                 = nqb.registry.get(nqb.n.getValue().Type())
		customIDPropertyName, customIDPropertyValue = metadata.getCustomID(*nqb.n.getValue())
		idCQLRef                                    = nSign + "ID"
		parameters                                  = map[string]interface{}{idCQLRef: nqb.n.getID()}
	)
	match := `MATCH (` + nSign + `)
	`
	filter := `	WHERE ID(` + nSign + `) = $` + idCQLRef + `
	`
	if customIDPropertyName != emptyString {
		filter = `WHERE ` + nSign + `.` + customIDPropertyName + ` = $` + idCQLRef + `
		`
		parameters[idCQLRef] = customIDPropertyValue.Interface()
	}

	return match + filter, parameters, nil
}
func (nqb nodeQueryBuilder) getSet() (string, map[string]interface{}) {

	var (
		set        string
		nSign      = nqb.n.getSignature()
		properties = map[string]interface{}{}
		parameters = map[string]interface{}{}
		propCQLRef = nSign + "Properties"
	)
	for propertyName, propertyValue := range nqb.deltaProperties {
		if !metaProperties[propertyName] {
			properties[propertyName] = propertyValue
		}
	}

	if len(properties) > 0 {
		set += `SET ` + nSign + ` += $` + propCQLRef + `
		`
		parameters[propCQLRef] = properties
	}

	if nqb.isLabelsDirty {
		set += `SET ` + nSign + `:` + nqb.n.getLabel() + `
		`
	}

	return set, parameters
}

func (nqb nodeQueryBuilder) getLoadAll(IDs interface{}, lo *LoadOptions) (string, map[string]interface{}) {

	var (
		depth                   = strconv.Itoa(lo.Depth)
		metadata, _             = nqb.registry.get(nqb.n.getValue().Type())
		customIDPropertyName, _ = metadata.getCustomID(*nqb.n.getValue())
		parameters              = map[string]interface{}{}
	)
	if lo.Depth == infiniteDepth {
		depth = emptyString
	}

	match := `MATCH path = (n:` + nqb.n.getLabel() + `)-[*0..` + depth + `]-()
	`
	var filter string
	if IDs != nil {
		filter = `WHERE ID(n) IN $ids 
		`
		if customIDPropertyName != emptyString {
			filter = `WHERE n.` + customIDPropertyName + ` IN $ids 
			`
		}
		parameters["ids"] = IDs
	}

	end := `WITH n, path, range(0, length(path) - 1) as index
	WITH  n, path, index, [i in index | CASE WHEN nodes(path)[i] = startNode(relationships(path)[i]) THEN false ELSE true END] as isDirectionInverted
	RETURN path, ID(n), isDirectionInverted
	`

	return match + filter + end, parameters
}

func (nqb nodeQueryBuilder) getDelete() (string, map[string]interface{}, map[string]graph) {
	match, parameters, _ := nqb.getMatch()
	delete := `DETACH DELETE ` + nqb.n.getSignature() + ` RETURN ID(` + nqb.n.getSignature() + `)
	`
	return match + delete, parameters, nil
}

func (nqb nodeQueryBuilder) getDeleteAll() (string, map[string]interface{}) {
	return `MATCH (n:` + nqb.n.getLabel() + `) DETACH DELETE n RETURN ID(n)`, nil
}

func (nqb nodeQueryBuilder) getCountEntitiesOfType() (string, map[string]interface{}) {
	return `MATCH (n:` + nqb.n.getLabel() + `) RETURN count(n) as count`, nil
}
