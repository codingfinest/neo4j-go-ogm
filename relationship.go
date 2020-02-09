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
	"strconv"
	"strings"
)

type relationship struct {
	Value      *reflect.Value
	coordinate *coordinate
	depth      *int

	ID         int64
	relType    string
	properties map[string]interface{}

	signature string

	nodes map[int64]graph
}

type direction int

const (
	outgoing direction = iota
	incoming
	undirected
)

const (
	startNode int64 = iota
	endNode
)

func (r *relationship) setDepth(depth *int) {
	r.depth = depth
}

func (r *relationship) getDepth() *int {
	return r.depth
}

func (r *relationship) getRelatedGraphs() map[int64]graph {
	return r.nodes
}

func (r *relationship) setRelatedGraph(g graph) {
	//TODO don't compare by value. compare by id?
	if *g.getValue() == *r.nodes[startNode].getValue() {
		r.nodes[startNode] = g
	} else if *g.getValue() == *r.nodes[endNode].getValue() {
		r.nodes[endNode] = g

	}
	g.setRelatedGraph(r)
}

func (r *relationship) getID() int64 {
	return r.ID
}

func (r *relationship) setID(ID int64) {
	r.ID = ID
	r.signature = strings.ReplaceAll("r"+strconv.FormatInt(ID, 10), "-", "_")
}

func (r *relationship) getValue() *reflect.Value {
	return r.Value
}

func (r *relationship) setValue(v *reflect.Value) {
	r.Value = v
}

func (r *relationship) getLabel() string {
	return r.getType()
}

func (r *relationship) setLabel(label string) {
	r.relType = label
}

func (r *relationship) getType() string {
	return r.relType
}

func (r *relationship) getProperties() map[string]interface{} {
	return r.properties
}

func (r *relationship) setProperties(p map[string]interface{}) {
	r.properties = p
}

func (r *relationship) getCoordinate() *coordinate {
	return r.coordinate
}

func (r *relationship) setCoordinate(c *coordinate) {
	r.coordinate = c
}

func (r *relationship) getSignature() string {
	return r.signature
}
