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
