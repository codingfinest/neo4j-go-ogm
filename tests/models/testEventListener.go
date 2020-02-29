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
	"sort"
	"time"

	gogm "github.com/codingfinest/neo4j-go-ogm"
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
