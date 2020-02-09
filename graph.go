package gogm

import "reflect"

type graph interface {
	getSignature() string

	getDepth() *int
	setDepth(*int)

	setID(int64)
	getID() int64

	setValue(*reflect.Value)
	getValue() *reflect.Value

	setProperties(map[string]interface{})
	getProperties() map[string]interface{}

	setLabel(label string)
	getLabel() string

	setCoordinate(*coordinate)
	getCoordinate() *coordinate

	setRelatedGraph(graph)
	getRelatedGraphs() map[int64]graph
}
