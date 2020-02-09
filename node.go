package gogm

import (
	"reflect"
	"strconv"
	"strings"
)

type node struct {
	ID            int64
	Value         *reflect.Value
	coordinate    *coordinate
	depth         *int
	label         string
	signature     string
	properties    map[string]interface{}
	relationships map[int64]graph
}

func (n *node) setDepth(depth *int) {
	n.depth = depth
}

func (n *node) getDepth() *int {
	return n.depth
}

func (n *node) getID() int64 {
	return n.ID
}

func (n *node) setID(ID int64) {
	n.ID = ID
	n.signature = strings.ReplaceAll("n"+strconv.FormatInt(ID, 10), "-", "_")
}

func (n *node) getValue() *reflect.Value {
	return n.Value
}

func (n *node) setValue(v *reflect.Value) {
	n.Value = v
}

func (n *node) getLabel() string {
	return n.label
}

func (n *node) setLabel(label string) {
	n.label = label
}

func (n *node) getProperties() map[string]interface{} {
	return n.properties
}

func (n *node) setProperties(p map[string]interface{}) {
	n.properties = p
}

func (n *node) getRelatedGraphs() map[int64]graph {
	return n.relationships
}

func (n *node) setRelatedGraph(g graph) {
	n.relationships[g.getID()] = g
}

func (n *node) getCoordinate() *coordinate {
	return n.coordinate
}

func (n *node) setCoordinate(c *coordinate) {
	n.coordinate = c
}

func (n *node) getSignature() string {
	return n.signature
}
