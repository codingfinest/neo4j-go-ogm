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
	"errors"
	"reflect"
	"sort"
	"strings"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type loader struct {
	cypherExecuter *cypherExecuter
	store          store
	eventer        eventer
	registry       *registry
	graphFactory   graphFactory
}

func newLoader(cypherExecuter *cypherExecuter, store store, eventer eventer, registry *registry, graphFactory graphFactory) *loader {
	return &loader{cypherExecuter, store, eventer, registry, graphFactory}
}

func (l *loader) load(object interface{}, ID interface{}, loadOptions *LoadOptions, reload bool) (store, error) {

	var (
		valueOfObject = reflect.ValueOf(object)
		valueOfID     = reflect.ValueOf(ID)
		graphs        []graph
		sliceOfIDs    = reflect.MakeSlice(reflect.SliceOf(valueOfID.Type()), 1, 1)
		ptrToSliceIDs = reflect.New(sliceOfIDs.Type())
		err           error
	)

	ptrToSliceIDs.Elem().Set(sliceOfIDs)
	ptrToSliceIDs.Elem().Index(0).Set(valueOfID)

	//object: **DomainObject
	if graphs, err = l.graphFactory.get(valueOfObject, map[int]bool{labels: true, relatedGraph: true}); err != nil {
		return nil, err
	}

	if loadOptions == nil {
		loadOptions = NewLoadOptions()
	}
	dummyValue := reflect.New(elem(reflect.TypeOf(object)).Elem())
	graphs[0].setValue(&dummyValue)
	sliceOfObjs, unloadedGraphs, err := l.loadAllOfGraphType(graphs[0], ptrToSliceIDs.Elem().Interface(), loadOptions, reload)

	if err != nil {
		return nil, err
	}

	if sliceOfObjs.Len() > 1 {
		return nil, errors.New("Got too many objects for ID " + valueOfID.String())
	} else if sliceOfObjs.Len() == 1 {
		valueOfObject.Elem().Set(sliceOfObjs.Index(0).Elem().Addr())
	}

	return unloadedGraphs, err
}

func (l *loader) loadAll(objects interface{}, IDs interface{}, loadOptions *LoadOptions) error {

	var (
		graphs      []graph
		sliceOfObjs reflect.Value
		err         error
	)

	//objects: *[]*DomainObject
	if graphs, err = l.graphFactory.get(reflect.ValueOf(objects), map[int]bool{labels: true, relatedGraph: true}); err != nil {
		return err
	}

	if loadOptions == nil {
		loadOptions = NewLoadOptions()
	}

	dummyValue := reflect.New(elem(reflect.TypeOf(objects)).Elem())
	graphs[0].setValue(&dummyValue)
	if sliceOfObjs, _, err = l.loadAllOfGraphType(graphs[0], IDs, loadOptions, false); err != nil {
		return err
	}

	valueOfSliceOfObjs := reflect.ValueOf(objects).Elem()
	valueOfSliceOfObjs.Set(reflect.AppendSlice(valueOfSliceOfObjs, sliceOfObjs))

	return nil
}

func (l *loader) reload(objects ...interface{}) error {
	var err error
	var graphs []graph
	var IDer = getIDer(nil, nil)
	var storedGraph graph
	for _, object := range objects {
		valueOfObject := reflect.ValueOf(object)
		//object: **DomainObject
		if graphs, err = l.graphFactory.get(valueOfObject, map[int]bool{labels: true, relatedGraph: true}); err != nil {
			return err
		}
		for _, graph := range graphs {
			IDer(graph)
		}
		if storedGraph = l.store.get(graphs[0]); storedGraph == nil || !storedGraph.getValue().IsValid() {
			continue
		}
		metadata, _ := l.registry.get(storedGraph.getValue().Type())

		ID := reflect.ValueOf(storedGraph.getID()).Interface()
		customIDName, customIDValue := metadata.getCustomID(*storedGraph.getValue())
		if customIDName != emptyString {
			ID = customIDValue.Interface()
		}

		lo := NewLoadOptions()
		if storedGraph.getDepth() != nil {
			lo.Depth = *storedGraph.getDepth() / 2
		}
		storedUnwound := unwind(storedGraph, lo.Depth)
		var unloadedGraphs store
		if unloadedGraphs, err = l.load(valueOfObject.Interface(), ID, lo, true); err != nil {
			return err
		}

		if unloadedGraphs != nil && unloadedGraphs.get(storedGraph) != nil {
			unloadedUnwound := unwind(unloadedGraphs.get(storedGraph), lo.Depth)
			for _, g := range storedUnwound.all() {
				if unloadedUnwound.get(g) == nil {
					deletedGraphs, updatedGraphs := l.store.delete(g)
					for _, updatedGraph := range updatedGraphs {
						notifyPostDelete(l.eventer, updatedGraph, UPDATE)
					}
					for _, deletedGraph := range deletedGraphs {
						notifyPostDelete(l.eventer, deletedGraph, DELETE)
					}
				}
			}
		}
	}
	return nil
}

func (l *loader) loadAllOfGraphType(refGraph graph, IDs interface{}, loadOptions *LoadOptions, reload bool) (reflect.Value, store, error) {

	var (
		typeOfRefGraph = reflect.TypeOf(refGraph)

		metadata, err = l.registry.get(refGraph.getValue().Type())

		sliceOfPtrToObjs = reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(refGraph.getValue().Type().Elem())), 0, 0)
		ptrToObjs        = reflect.New(sliceOfPtrToObjs.Type())

		records     []neo4j.Record
		storedGraph graph
	)

	if loadOptions == nil {
		loadOptions = NewLoadOptions()
	}

	if err != nil {
		return invalidValue, nil, err
	}
	customIDName, _ := metadata.getCustomID(*refGraph.getValue())

	ptrToObjs.Elem().Set(sliceOfPtrToObjs)

	var ids interface{}
	if IDs == nil {
		ids = nil
	} else {
		valueOfIDs := reflect.ValueOf(IDs)
		sliceOfIDsToLoad := reflect.MakeSlice(reflect.SliceOf(valueOfIDs.Type().Elem()), 0, 0)
		IDsToLoad := reflect.New(sliceOfIDsToLoad.Type())
		IDsToLoad.Elem().Set(sliceOfIDsToLoad)

		if loadOptions.Depth <= -1 || reload {
			IDsToLoad.Elem().Set(valueOfIDs)
		} else {
			for i := 0; i < valueOfIDs.Len(); i++ {
				ID := valueOfIDs.Index(i)
				if customIDName != emptyString {
					storedGraph = l.store.getByCustomID(*refGraph.getValue(), typeOfRefGraph, ID.Interface())
				} else {
					var id int64
					var ok bool
					if id, ok = ID.Interface().(int64); !ok {
						return invalidValue, nil, errors.New("Unexpected type of ID on load. In the absence of a custom ID field in " + refGraph.getValue().Type().String() + ", expected an ID of int type for domain object on load")
					}
					refGraph.setID(id)
					storedGraph = l.store.get(refGraph)
				}
				if storedGraph != nil && storedGraph.getDepth() != nil && loadOptions.Depth*2 <= *storedGraph.getDepth() {
					ptrToObjs.Elem().Set(reflect.Append(ptrToObjs.Elem(), *storedGraph.getValue()))
				} else {
					IDsToLoad.Elem().Set(reflect.Append(IDsToLoad.Elem(), ID))
				}
			}
		}

		if IDsToLoad.Elem().Len() == 0 {
			return ptrToObjs.Elem(), nil, nil
		}

		ids = IDsToLoad.Elem().Interface()
	}

	var cypherBuilder graphQueryBuilder
	if cypherBuilder, err = newCypherBuilder(refGraph, l.registry, nil); err != nil {
		return invalidValue, nil, err
	}

	if records, err = neo4j.Collect(l.cypherExecuter.exec(cypherBuilder.getLoadAll(ids, loadOptions))); err != nil {
		return invalidValue, nil, err
	}

	toUnLoad := newstore(nil)
	visitedGraphs := newstore(nil)
	unloadedGrahps := newstore(nil)
	relatedValues := make(map[reflect.Type]map[int64]map[int64]bool)
	relatedValues[typeOfPrivateNode] = map[int64]map[int64]bool{}
	relatedValues[typeOfPrivateRelationship] = map[int64]map[int64]bool{}

	for _, record := range records {
		refGraph.setID(record.GetByIndex(1).(int64))
		toUnLoad.save(l.getGraphToLoadFromDBResult(record.GetByIndex(0).(neo4j.Path), record.GetByIndex(2).([]interface{}), refGraph, visitedGraphs, loadOptions.Depth))
	}

	for _, g := range toUnLoad.all() {
		g.setCoordinate(&coordinate{0, 0})
		var loadDepth = -1
		if loadDepth, err = l.unloadDBObject(g, unloadedGrahps, loadOptions.Depth, relatedValues); err != nil {
			return invalidValue, nil, err
		}

		if loadDepth > -1 {
			g.setDepth(&loadDepth)
		}

	}

	for _, g := range unloadedGrahps.all() {
		g.setCoordinate(nil)
		isRoot := toUnLoad.get(g) != nil
		if stored := l.store.get(g); !reload && stored != nil && stored.getDepth() != nil && g.getDepth() != nil && *stored.getDepth() >= *g.getDepth() {
			if stored.getValue().IsValid() {
				for _, eventListener := range l.eventer.eventListeners {
					eventListener.OnPostLoad(event{object: stored.getValue()})
				}
			}
			if isRoot {
				ptrToObjs.Elem().Set(reflect.Append(ptrToObjs.Elem(), *stored.getValue()))
			}
			continue
		}

		l.store.save(g)
		if g.getValue().IsValid() {
			for _, eventListener := range l.eventer.eventListeners {
				eventListener.OnPostLoad(event{object: g.getValue()})
			}
		}

		if isRoot {
			ptrToObjs.Elem().Set(reflect.Append(ptrToObjs.Elem(), *g.getValue()))
		}
	}

	return ptrToObjs.Elem(), toUnLoad, nil
}

func (l *loader) getGraphToLoadFromDBResult(path neo4j.Path, isDirectionInverted []interface{}, refGraph graph, visitedGraphs store, depth int) graph {

	nodes := path.Nodes()
	relationships := path.Relationships()
	internalGraphType := reflect.TypeOf(refGraph)
	graphToLoadType := refGraph.getValue().Type().Elem()
	ID := refGraph.getID()

	graphToLoad := visitedGraphs.get(refGraph)

	for index, neoRelationship := range relationships {
		from := nodes[index]
		to := nodes[index+1]

		if visitedGraphs.relationship(neoRelationship.Id()) == nil {

			var fromNode, toNode graph
			if fromNode = visitedGraphs.node(from.Id()); fromNode == nil {
				labels := from.Labels()
				sort.Strings(labels)
				fromNode = &node{
					properties:    from.Props(),
					label:         strings.Join(labels, labelsDelim),
					relationships: map[int64]graph{}}
				fromNode.setID(from.Id())
				fromNode.getProperties()[idPropertyName] = from.Id()
				visitedGraphs.save(fromNode)
			}
			if toNode = visitedGraphs.node(to.Id()); toNode == nil {
				labels := to.Labels()
				sort.Strings(labels)
				toNode = &node{
					properties:    to.Props(),
					label:         strings.Join(labels, labelsDelim),
					relationships: map[int64]graph{}}
				toNode.setID(to.Id())
				toNode.getProperties()[idPropertyName] = to.Id()
				visitedGraphs.save(toNode)
			}

			nodes := map[int64]graph{startNode: fromNode, endNode: toNode}
			if isDirectionInverted[index].(bool) {
				nodes = map[int64]graph{startNode: toNode, endNode: fromNode}
			}

			fromNodeToNode := &relationship{
				relType:    neoRelationship.Type(),
				properties: neoRelationship.Props(),
				nodes:      nodes}
			fromNodeToNode.setID(neoRelationship.Id())
			fromNodeToNode.getProperties()[idPropertyName] = neoRelationship.Id()
			visitedGraphs.save(fromNodeToNode)

			fromNode.setRelatedGraph(fromNodeToNode)
			toNode.setRelatedGraph(fromNodeToNode)

			if graphToLoad == nil {
				if internalGraphType == typeOfPrivateNode {
					if from.Id() == ID {
						graphToLoad = fromNode
					} else if to.Id() == ID {
						graphToLoad = toNode
					}
				} else if internalGraphType == typeOfPrivateRelationship && neoRelationship.Id() == ID {
					graphToLoad = fromNodeToNode
				}
			}
		}
	}

	if graphToLoad == nil && len(relationships) == 0 {
		node := &node{
			properties:    nodes[0].Props(),
			label:         strings.Join(nodes[0].Labels(), labelsDelim),
			relationships: map[int64]graph{}}
		node.setID(nodes[0].Id())
		node.getProperties()[idPropertyName] = nodes[0].Id()
		visitedGraphs.save(node)
		graphToLoad = node
	}

	if graphToLoad.getValue() == nil {
		v := reflect.New(graphToLoadType)
		graphToLoad.setValue(&v)
	}

	return graphToLoad
}

func (l *loader) unloadDBObject(g graph, unloadedGrahps store, depth int, relatedValues map[reflect.Type]map[int64]map[int64]bool) (int, error) {

	var (
		err                               error
		graphfield                        *field
		firstMetadata, graphFieldMetadata metadata

		queue       = []graph{g}
		first       graph
		loadedDepth = -1
	)

	depthToLoad := maxDepth
	if depth > infiniteDepth {
		depthToLoad = depth * 2
	}

	for len(queue) > 0 {
		first = queue[0]
		queue[0] = nil
		queue = queue[1:]

		if reflect.TypeOf(first) == typeOfPrivateRelationship && first.getCoordinate().depth > depthToLoad {
			break
		}

		if first.getValue().IsValid() {
			if firstMetadata, err = l.registry.get(first.getValue().Type()); err != nil {
				return -1, err
			}
		}

		if unloadedGrahps.get(first) == nil {
			if first.getValue().IsValid() {
				driverPropertiesAsStructFieldValues(first.getProperties(), firstMetadata.getPropertyStructFields())
				unloadGraphProperties(first, firstMetadata.getPropertyStructFields())
			}
			unloadedGrahps.save(first)
		}

		loadedDepth = first.getCoordinate().depth

		if relatedValues[reflect.TypeOf(first)][first.getID()] == nil {
			relatedValues[reflect.TypeOf(first)][first.getID()] = map[int64]bool{}
		}

		for _, relatedGraph := range first.getRelatedGraphs() {
			if relatedValues[reflect.TypeOf(relatedGraph)][relatedGraph.getID()] == nil {
				relatedValues[reflect.TypeOf(relatedGraph)][relatedGraph.getID()] = map[int64]bool{}
			}

			if relatedGraph.getCoordinate() == nil {
				relatedGraph.setCoordinate(&coordinate{first.getCoordinate().depth + 1, 0})
			}

			// Related graph Value has not been added to the current graph Value
			if !relatedValues[reflect.TypeOf(first)][first.getID()][relatedGraph.getID()] {

				if graphfield, err = firstMetadata.getGraphField(first, relatedGraph); err != nil {
					return -1, err
				}

				if graphfield == nil {
					continue
				}

				typeOfGraphField := elem(graphfield.getStructField().Type)
				if graphFieldMetadata, err = l.registry.get(typeOfGraphField); err != nil {
					return -1, err
				}

				value := reflect.New(typeOfGraphField.Elem())

				addDomainObject(graphfield, value)
				relatedGraph.setValue(&value)

				//first is a Node and related graph field is a Node  field. relatedGraph is relationship A
				if first.getValue().IsValid() && reflect.TypeOf(firstMetadata) == reflect.TypeOf(graphFieldMetadata) {
					relatedGraph.setValue(&invalidValue)
					otherNode := relatedGraph.getRelatedGraphs()[startNode]
					if otherNode.getID() == first.getID() {
						otherNode = relatedGraph.getRelatedGraphs()[endNode]
					}
					otherNode.setValue(&value)

					relatedValues[reflect.TypeOf(relatedGraph)][relatedGraph.getID()][otherNode.getID()] = true
					if relatedValues[reflect.TypeOf(otherNode)][otherNode.getID()] == nil {
						relatedValues[reflect.TypeOf(otherNode)][otherNode.getID()] = map[int64]bool{}
					}
					relatedValues[reflect.TypeOf(otherNode)][otherNode.getID()][relatedGraph.getID()] = true
				}
				relatedValues[reflect.TypeOf(first)][first.getID()][relatedGraph.getID()] = true
				relatedValues[reflect.TypeOf(relatedGraph)][relatedGraph.getID()][first.getID()] = true

			}

			if unloadedGrahps.get(relatedGraph) == nil {
				queue = append(queue, relatedGraph)
			}
		}
	}
	return loadedDepth, nil
}
