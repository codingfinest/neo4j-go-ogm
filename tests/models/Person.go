package models

type Person struct {
	TestNodeEntity
	Name    string
	Born    int64
	Follows []*Person `gogm:"reltype:FOLLOWS,direction:<-"`
	Tags    []string  `gogm:"label"`
}

type Person2 struct {
	TestNodeEntity
	Name    string
	Born    int64
	Follows []*Person2 `gogm:"reltype:FOLLOWS,direction:--"`
	Tags    []string   `gogm:"label"`
}
