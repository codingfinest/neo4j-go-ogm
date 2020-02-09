package models

import (
	"sort"
	"time"

	"github.com/codingfinest/neo4j-go-ogm/gogm"
)

type TestEventListener struct{}

func (e *TestEventListener) OnPreSave(event gogm.Event) {
	switch object := event.GetObject().(type) {
	case *SimpleNode, *SimpleRelationship, *Node0, *Node1, *Node2, *Node3, *Node4, *Node5, *Node6, *Node7, *Node8, *Node10:
		testObj := object.(TestObject)
		testObj.ClearMetaTimestamps()
	}
}

func (e *TestEventListener) OnPostSave(event gogm.Event) {
	switch object := event.GetObject().(type) {
	case *Person2:
		switch event.GetLifeCycle() {
		case gogm.CREATE:
			object.SetCreatedTime(time.Now().Unix())
		case gogm.UPDATE:
			object.SetUpdatedTime(time.Now().Unix())
		}
		sort.SliceStable(object.Follows, func(i, j int) bool { return object.Follows[i].Name < object.Follows[j].Name })
	case *Person:
		switch event.GetLifeCycle() {
		case gogm.CREATE:
			object.SetCreatedTime(time.Now().Unix())
		case gogm.UPDATE:
			object.SetUpdatedTime(time.Now().Unix())
		}
		sort.SliceStable(object.Follows, func(i, j int) bool { return object.Follows[i].Name < object.Follows[j].Name })
	case *SimpleNode, *SimpleRelationship, *Node0, *Node1, *Node2, *Node3, *Node4, *Node5, *Node6, *Node7, *Node8, *Node10:
		testObj := object.(TestObject)
		switch event.GetLifeCycle() {
		case gogm.CREATE:
			testObj.SetCreatedTime(time.Now().Unix())
		case gogm.UPDATE:
			testObj.SetUpdatedTime(time.Now().Unix())
		}
	}
}

func (e *TestEventListener) OnPostLoad(event gogm.Event) {
	switch object := event.GetObject().(type) {
	case *Person2:
		sort.SliceStable(object.Follows, func(i, j int) bool { return object.Follows[i].Name < object.Follows[j].Name })
	case *Person:
		sort.SliceStable(object.Follows, func(i, j int) bool { return object.Follows[i].Name < object.Follows[j].Name })
	case *Movie:
		sort.SliceStable(object.Characters, func(i, j int) bool { return object.Characters[i].Name < object.Characters[j].Name })
	}
}
func (e *TestEventListener) OnPreDelete(event gogm.Event) {}

func (e *TestEventListener) OnPostDelete(event gogm.Event) {
	switch object := event.GetObject().(type) {
	case *SimpleNode, *SimpleRelationship, *Node0, *Node1, *Node2, *Node3, *Node4, *Node5, *Node6, *Node7, *Node8, Node10:
		testObj := object.(TestObject)
		switch event.GetLifeCycle() {
		case gogm.UPDATE:
			testObj.SetUpdatedTime(time.Now().Unix())
		case gogm.DELETE:
			testObj.SetDeletedTime(time.Now().Unix())
		}
	}
}
