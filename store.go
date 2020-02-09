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

package gogm

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

//store is a container of graphs. An instance of it is used as a cache for the session and other instances
// of it are used for keeping track of graphs during a graph traversal.
type store interface {
	// all returns all graphs in the store
	all() []graph

	// get returns a graph in the store with the same ID and typeOf the graph
	get(graph) graph

	// save saves a graph in the store
	save(graph)

	// clear reset the store state
	clear() error

	// delete removes a graph from the store and returns graph that were deleted in the process
	// and graphs that were affected by the delete (updatedGraphs)
	delete(graph) (deletedGraphs []graph, updatedGraphs []graph)

	// clears the store, but returns all deleted graphs
	purge() []graph

	// node returns the node graph with ID
	node(int64) graph

	// relationship returns the relationship graph with ID
	relationship(int64) graph

	// getByCustomID returns the graph with whose custom ID is interface{}
	getByCustomID(reflect.Value, reflect.Type, interface{}) graph

	// print prints the store
	print()
}

type storeImpl struct {
	registry *registry

	nodes          map[int64]graph
	relationships  map[int64]graph
	relationshipsA map[int64]map[int64]*int64
	customIDs      map[string]map[interface{}]*int64
	nodesMu        sync.Mutex
}

func newstore(registry *registry) *storeImpl {
	return &storeImpl{registry, map[int64]graph{}, map[int64]graph{}, map[int64]map[int64]*int64{}, map[string]map[interface{}]*int64{}, sync.Mutex{}}
}

func (s *storeImpl) get(g graph) graph {
	var storedGraph graph
	switch t := reflect.TypeOf(g); t {
	case typeOfPrivateNode:
		s.nodesMu.Lock()
		defer s.nodesMu.Unlock()
		storedGraph = s.nodes[g.getID()]
		break
	case typeOfPrivateRelationship:
		internalID := g.getID()
		if g.getValue() != nil && !g.getValue().IsValid() {
			internalID = initialGraphID
			relatedGraphs := g.getRelatedGraphs()

			startIDPtrAddr := getIDAddr(relatedGraphs[startNode])
			endIDPtrAddr := getIDAddr(relatedGraphs[endNode])

			if startIDPtrAddr != nil && endIDPtrAddr != nil && *startIDPtrAddr != nil && *endIDPtrAddr != nil && s.relationshipsA[**startIDPtrAddr] != nil && s.relationshipsA[**startIDPtrAddr][**endIDPtrAddr] != nil {
				internalID = *s.relationshipsA[**startIDPtrAddr][**endIDPtrAddr]
			}
		}
		storedGraph = s.relationships[internalID]
		break
	}
	return storedGraph
}

func (s *storeImpl) save(g graph) {
	switch t := reflect.TypeOf(g); t {
	case typeOfPrivateNode:
		s.nodesMu.Lock()
		defer s.nodesMu.Unlock()
		s.nodes[g.getID()] = g
		break
	case typeOfPrivateRelationship:
		s.relationships[g.getID()] = g
		if g.getValue() != nil && !g.getValue().IsValid() {
			relatedGraphs := g.getRelatedGraphs()
			if s.relationshipsA[relatedGraphs[startNode].getID()] == nil {
				s.relationshipsA[relatedGraphs[startNode].getID()] = map[int64]*int64{}
			}
			internalID := g.getID()
			s.relationshipsA[relatedGraphs[startNode].getID()][relatedGraphs[endNode].getID()] = &internalID
		}
		break
	}

	if s.registry != nil && g.getValue() != nil && g.getValue().IsValid() {
		vType := g.getValue().Type()
		metadata, _ := s.registry.get(vType)
		customIDName, customIDValue := metadata.getCustomID(*g.getValue())
		if customIDName != emptyString {
			internalID := g.getID()
			if s.customIDs[vType.String()] == nil {
				s.customIDs[vType.String()] = map[interface{}]*int64{}
			}
			s.customIDs[vType.String()][customIDValue.Interface()] = &internalID
		}
	}
}

func (s *storeImpl) delete(g graph) ([]graph, []graph) {
	var deletedGraphs []graph
	var updatedGraphs []graph
	switch t := reflect.TypeOf(g); t {
	case typeOfPrivateNode:
		node := s.nodes[g.getID()]
		if node != nil {
			for _, relatedGraph := range node.getRelatedGraphs() {
				relatedDeletedGraphs, relatedUpdatedGraphs := s.delete(relatedGraph)
				if relatedUpdatedGraphs != nil {
					if len(relatedUpdatedGraphs) == 2 {
						if relatedUpdatedGraphs[0] == node {
							relatedUpdatedGraphs[0] = nil
						} else {
							relatedUpdatedGraphs[1] = nil
						}
					}
					updatedGraphs = append(updatedGraphs, relatedUpdatedGraphs...)
				}
				if relatedDeletedGraphs != nil {
					deletedGraphs = append(deletedGraphs, relatedDeletedGraphs...)
				}
			}
			delete(s.nodes, node.getID())
			deletedGraphs = append(deletedGraphs, node)
		}
		break
	case typeOfPrivateRelationship:
		relationship := s.relationships[g.getID()]
		if relationship != nil {
			if !relationship.getValue().IsValid() {
				relatedGraphs := relationship.getRelatedGraphs()
				delete(s.relationshipsA[relatedGraphs[startNode].getID()], relatedGraphs[endNode].getID())
			}

			delete(relationship.getRelatedGraphs()[startNode].getRelatedGraphs(), relationship.getID())
			delete(relationship.getRelatedGraphs()[endNode].getRelatedGraphs(), relationship.getID())
			delete(s.relationships, relationship.getID())
			updatedGraphs = append(updatedGraphs, relationship.getRelatedGraphs()[startNode], relationship.getRelatedGraphs()[endNode])
			deletedGraphs = append(deletedGraphs, relationship)
		}
		break
	}

	if s.registry != nil && g.getValue() != nil && g.getValue().IsValid() {
		vType := g.getValue().Type()
		metadata, _ := s.registry.get(vType)
		customIDName, customIDValue := metadata.getCustomID(*g.getValue())
		if customIDName != emptyString {
			if s.customIDs[vType.String()] != nil && s.customIDs[vType.String()][customIDValue.Interface()] != nil {
				delete(s.customIDs[vType.String()], customIDValue.Interface())
			}
		}
	}

	return deletedGraphs, updatedGraphs
}

func (s *storeImpl) clear() error {
	s.nodes = map[int64]graph{}
	s.relationships = map[int64]graph{}
	s.relationshipsA = map[int64]map[int64]*int64{}
	s.customIDs = map[string]map[interface{}]*int64{}
	return nil
}

func (s *storeImpl) purge() []graph {
	var deletedGraphs []graph
	for _, node := range s.nodes {
		nodeDeletedGraphs, _ := s.delete(node)
		deletedGraphs = append(deletedGraphs, nodeDeletedGraphs...)
	}
	return deletedGraphs
}

func (s *storeImpl) node(ID int64) graph {
	return s.nodes[ID]
}

func (s *storeImpl) relationship(ID int64) graph {
	return s.relationships[ID]
}

func (s *storeImpl) all() []graph {
	var allGraphs []graph
	for _, node := range s.nodes {
		allGraphs = append(allGraphs, node)
	}
	for _, relationship := range s.relationships {
		allGraphs = append(allGraphs, relationship)
	}
	return allGraphs
}

func (s *storeImpl) getByCustomID(v reflect.Value, typeOfRefGraph reflect.Type, idValue interface{}) graph {
	typeName := v.Type().String()
	mapToSearch := s.nodes
	if typeOfRefGraph == typeOfPrivateRelationship {
		mapToSearch = s.relationships
	}

	if s.customIDs[typeName] != nil && s.customIDs[typeName][idValue] != nil {
		return mapToSearch[*s.customIDs[typeName][idValue]]
	}
	return nil
}

func unwind(g graph, depth int) store {
	visited := newstore(nil)
	maxDepth := depth * 2
	g.setCoordinate(&coordinate{0, 0, 0})
	queue := []graph{g}
	for len(queue) > 0 {
		if reflect.TypeOf(queue[0]) == typeOfPrivateRelationship && queue[0].getCoordinate().depth > maxDepth {
			break
		}
		visited.save(queue[0])
		for _, relatedGraph := range queue[0].getRelatedGraphs() {
			if visited.get(relatedGraph) == nil {
				relatedGraph.setCoordinate(&coordinate{queue[0].getCoordinate().depth + 1, 0, 0})
				queue = append(queue, relatedGraph)
			}
		}
		queue[0].setCoordinate(nil)
		queue[0] = nil
		queue = queue[1:]
	}
	for _, queued := range queue {
		queued.setCoordinate(nil)
	}
	return visited
}

func (s *storeImpl) print() {

	gPrint := func(g graph) string {
		s := strconv.FormatInt(g.getID(), 10) + `|` + g.getLabel() + `
		`
		var rStrings []string
		for _, r := range g.getRelatedGraphs() {
			rStrings = append(rStrings, strconv.FormatInt(r.getID(), 10)+`|`+r.getLabel()+``)
		}
		s = s + strings.Join(rStrings, ", ") + `
		`
		return s
	}

	fmt.Println("**************** NODES ****************************")
	for _, node := range s.nodes {
		fmt.Println(gPrint(node))
	}

	fmt.Println("**************** RELATIONSHIPS ****************************")
	for _, relationship := range s.relationships {
		fmt.Println(gPrint(relationship))
	}
}
