package models

import "github.com/codingfinest/neo4j-go-ogm/gogm"

type TestRelationshipEntity struct {
	gogm.Relationship
	TestEntity
}

//Relationships

type SimpleRelationship struct {
	TestRelationshipEntity
	N5     *Node5 `gogm:"startNode"`
	N4     *Node4 `gogm:"endNode"`
	Name   string
	TestID string `gogm:"id"`
}

type SimpleRelationship2 struct {
	TestRelationshipEntity
	N5   *Node5 `gogm:"startNode"`
	N4   *Node4 `gogm:"endNode"`
	Name string
}
