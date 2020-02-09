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
