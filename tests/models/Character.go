package models

import "github.com/codingfinest/neo4j-go-ogm/gogm"

type Character struct {
	gogm.Relationship `gogm:"reltype:ACTED_IN"`

	Roles          []string
	Name           string
	NumberOfScenes int64
	Actor          *Actor `gogm:"startNode"`
	Movie          *Movie `gogm:"endNode"`
}
