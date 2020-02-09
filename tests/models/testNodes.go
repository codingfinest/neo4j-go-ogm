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

package models

import (
	"time"

	"github.com/codingfinest/neo4j-go-ogm/gogm"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type TestNodeEntity struct {
	gogm.Node
	TestEntity
}

type SimpleNode struct {
	TestNodeEntity
	Prop1 string
}

//(node0)-->(node1)-->(node2)<--(node3)-->(node4)
type Node0 struct {
	TestNodeEntity
	Name             string
	MapProps         map[string]int
	AliasedMapProps1 map[string]*string `gogm:"name:aliasedMapProp"`
	InvalidIDMapProp map[int]int
	MapToInterfaces  map[string]interface{}
	N1               *Node1
}
type Node1 struct {
	TestNodeEntity
	N0   *Node0 `gogm:"direction:<-"`
	N2   *Node2
	Name string
}
type Node2 struct {
	TestNodeEntity
	N1   *Node1 `gogm:"direction:<-"`
	N3   *Node3 `gogm:"direction:<-"`
	Name string
}
type Node3 struct {
	TestNodeEntity
	N2   *Node2
	N4   *Node4
	Name string
}
type Node4 struct {
	TestNodeEntity
	N3   *Node3   `gogm:"direction:<-"`
	N5s  []*Node5 `gogm:"direction:<-"`
	R1   *SimpleRelationship
	R2s  []*SimpleRelationship2
	Name string
}

type Node5 struct {
	TestNodeEntity
	N4   *Node4
	R1   *SimpleRelationship
	Name string
}

type Node6 struct {
	TestNodeEntity
	N7   *Node7 `gogm:"reltype:REL1"`
	N8   *Node8 `gogm:"reltype:REL1"`
	Name string
}

type Node7 struct {
	TestNodeEntity
	N6   *Node6 `gogm:"reltype:REL1,direction:<-"`
	Name string
}

type Node8 struct {
	TestNodeEntity
	N6   *Node6 `gogm:"reltype:REL1,direction:<-"`
	Name string
}

type Node9 struct {
	TestNodeEntity
	Name   string `gogm:"unique"`
	TestId string `gogm:"id,name:IDs"`
}

type Node10 struct {
	TestNodeEntity
	ZeroTime     time.Time
	NilTime1     *time.Time
	ZeroDuration neo4j.Duration

	Time  time.Time
	Time1 *time.Time

	Duration  neo4j.Duration
	Duration1 *neo4j.Duration
}

type InvalidID struct {
	TestNodeEntity
	TestId *string `gogm:"id,name:IDs"`
}
