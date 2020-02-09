package gogm

import (
	"strconv"
)

type relationshipQueryBuilder struct {
	r               *relationship
	registry        *registry
	deltaProperties map[string]interface{}
}

func (rqb relationshipQueryBuilder) getGraph() graph {
	return rqb.r
}

func newRelationshipCypherBuilder(r *relationship, registry *registry, stored graph) relationshipQueryBuilder {
	deltaProperties := r.getProperties()
	if stored != nil {
		deltaProperties = diffProperties(deltaProperties, stored.getProperties())
	}
	return relationshipQueryBuilder{
		r,
		registry,
		deltaProperties}
}

func (rqb relationshipQueryBuilder) getRemovedGraphs() (map[int64]graph, map[int64]graph) {
	return nil, nil
}

func (rqb relationshipQueryBuilder) isGraphDirty() bool {
	return rqb.r.getID() < 0 || len(rqb.deltaProperties) > 0
}

func (rqb relationshipQueryBuilder) getCreate() (string, string, map[string]interface{}, map[string]graph) {
	var (
		r         = rqb.r
		startSign = r.nodes[startNode].getSignature()
		endSign   = r.nodes[endNode].getSignature()
		rSign     = r.getSignature()
	)
	create := `CREATE (` + startSign + `)-[` + rSign + `:` + r.getType() + `]->(` + endSign + `)
	`
	return "", create, nil, map[string]graph{startSign: r.nodes[startNode], endSign: r.nodes[endNode]}
}

func (rqb relationshipQueryBuilder) getMatch() (string, map[string]interface{}, map[string]graph) {
	var (
		r         = rqb.r
		startSign = r.nodes[startNode].getSignature()
		endSign   = r.nodes[endNode].getSignature()
		rSign     = r.getSignature()
		match     = `MATCH (` + startSign + `)-[` + rSign + `:` + r.getType() + `]->(` + endSign + `)
		`
	)
	return match, nil, map[string]graph{startSign: r.nodes[startNode], endSign: r.nodes[endNode]}
}

func (rqb relationshipQueryBuilder) getSet() (string, map[string]interface{}) {
	var (
		r          = rqb.r
		rSign      = r.getSignature()
		properties = map[string]interface{}{}
		parameters = map[string]interface{}{}
		propCQLRef = rSign + "Properties"
		set        string
	)
	for propertyName, propertyValue := range r.getProperties() {
		if !metaProperties[propertyName] {
			properties[propertyName] = propertyValue
		}
	}

	if len(properties) > 0 {
		set += `SET ` + rSign + ` += $` + propCQLRef + `
		`
		parameters[propCQLRef] = properties
	}

	return set, parameters
}

func (rqb relationshipQueryBuilder) getLoadAll(IDs interface{}, lo *LoadOptions) (string, map[string]interface{}) {

	var (
		depth                   = strconv.Itoa(lo.Depth)
		metadata, _             = rqb.registry.get(rqb.r.getValue().Type())
		customIDPropertyName, _ = metadata.getCustomID(*rqb.r.getValue())
		parameters              = map[string]interface{}{}
	)

	if lo.Depth == infiniteDepth {
		depth = ""
	}

	match := `MATCH path = ()-[*0..` + depth + `]-()-[r:` + rqb.r.getLabel() + `]-()-[*0..` + depth + `]-()
	`

	var filter string
	if IDs != nil {
		filter = `WHERE ID(r) IN $ids 
		`
		if customIDPropertyName != emptyString {
			filter = `WHERE r.` + customIDPropertyName + ` IN $ids 
			`
		}
		parameters["ids"] = IDs
	}

	end := `WITH r, path, range(0, length(path) - 1) as index
	WITH  r, path, index, [i in index | CASE WHEN nodes(path)[i] = startNode(relationships(path)[i]) THEN false ELSE true END] as isDirectionInverted
	RETURN path, ID(r), isDirectionInverted
	`

	return match + filter + end, parameters
}

func (rqb relationshipQueryBuilder) getDeleteAll() (string, map[string]interface{}) {
	return `MATCH ()-[r:` + rqb.r.getType() + `]-()
	DELETE r
	RETURN ID(r)`, nil
}

func (rqb relationshipQueryBuilder) getDelete() (string, map[string]interface{}, map[string]graph) {
	rSign := rqb.r.getSignature()
	delete, _, depedencies := rqb.getMatch()
	delete += `DELETE ` + rSign + ` RETURN ID(` + rSign + `)
	`
	return delete, nil, depedencies
}

func (rqb relationshipQueryBuilder) getCountEntitiesOfType() (string, map[string]interface{}) {
	return `MATCH ()-[r:` + rqb.r.getType() + `]->() RETURN count(r) as count`, nil
}
