package models

type Director struct {
	Person
	Movies []*Movie
}
