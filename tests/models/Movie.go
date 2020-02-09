package models

import "github.com/codingfinest/neo4j-go-ogm/gogm"

type Movie struct {
	gogm.Node `gogm:"label:FILM,label:PICTURE"`

	Title    string
	Released int64
	Tagline  string

	Characters []*Character `gogm:"direction:INCOMING"`
	Directors  []*Director  `gogm:"direction:INCOMING"`
}

func (m *Movie) AddCharacter(c *Character) {
	m.Characters = append(m.Characters, c)
}
