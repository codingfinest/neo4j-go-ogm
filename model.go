package gogm

type Object struct {
	ID *int64 `json:"id"`
}

type Relationship struct {
	Object
}

type Node struct {
	Object
}
