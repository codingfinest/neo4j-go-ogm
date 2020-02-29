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
	"math"
	"reflect"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

const (
	idPropertyName            = "id"
	labelsDelim               = ":"
	emptyString               = ""
	spaceString               = " "
	initialGraphID      int64 = -1
	infiniteDepth             = -1
	maxDepth                  = math.MaxInt32 / 2
	defaultRelTypeDelim       = "_"
	indexDelim                = ","
	statementDelim            = ";\n"
	mapPropDelim              = "."
)

const (
	nodeCreateClause clause = iota
	relationshipCreateClause
	matchClause
	setClause
	deleteClause
)

const (
	relatedGraph int = iota
	labels
	properties
)

var typeOfPublicNode = reflect.TypeOf(Node{})
var typeOfPublicRelationship = reflect.TypeOf(Relationship{})
var typeOfPrivateNode = reflect.TypeOf(&node{})
var typeOfPrivateRelationship = reflect.TypeOf(&relationship{})
var typeOfNodeMetadata = reflect.TypeOf(&nodeMetadata{})

var invalidValue = reflect.ValueOf(nil)

var directionTags = map[string]direction{
	"<-": incoming,
	"->": outgoing,
	"--": undirected}

var clauseGroups = [5]clause{
	matchClause,
	nodeCreateClause,
	relationshipCreateClause,
	setClause,
	deleteClause}

type clause int
type clauses map[clause][]string

type LogLevel int

const (
	NONE    LogLevel = 0
	ERROR            = 1
	WARNING          = 2
	INFO             = 3
	DEBUG            = 4
)

var logLevels = map[LogLevel]neo4j.LogLevel{
	ERROR:   neo4j.ERROR,
	WARNING: neo4j.WARNING,
	INFO:    neo4j.INFO,
	DEBUG:   neo4j.DEBUG,
}

//Gogm is an instance of the OGM
type Gogm struct {
	uri      string
	username string
	password string
	logLevel LogLevel
	driver   neo4j.Driver
}

//New creates a new instance of the OGM
func New(uri string, username string, password string, logLevel LogLevel) *Gogm {
	return &Gogm{
		uri,
		username,
		password,
		logLevel,
		nil}
}

//NewSession creates a new session on an OGM instance
func (g *Gogm) NewSession(isWriteMode bool) (Session, error) {

	var err error
	var accessMode neo4j.AccessMode = neo4j.AccessModeRead
	if isWriteMode {
		accessMode = neo4j.AccessModeWrite
	}

	if g.driver == nil {
		if g.driver, err = getDriver(g.uri, g.username, g.password, g.logLevel); err != nil {
			return nil, err
		}
	}

	cypherExecutor := newCypherExecuter(g.driver, accessMode, nil)
	registry := newRegistry(*cypherExecutor)
	graphFactory := newGraphFactory(registry)
	transactioner := newTransactioner(accessMode)
	eventer := newEventer()
	store := newstore(registry)
	saver := newSaver(cypherExecutor, store, *eventer, registry, *graphFactory)
	loader := newLoader(cypherExecutor, store, *eventer, registry, *graphFactory)
	deleter := newDeleter(cypherExecutor, store, *eventer, registry, *graphFactory)
	queryer := newQueryer(cypherExecutor, *graphFactory, registry)

	return &sessionImpl{
		cypherExecutor,
		saver,
		loader,
		deleter,
		queryer,
		transactioner,
		store,
		registry,
		g.driver,
		eventer}, nil
}

func getDriver(uri string, username string, password string, logLevel LogLevel) (neo4j.Driver, error) {
	var (
		err    error
		driver neo4j.Driver
	)

	if driver, err = neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(config *neo4j.Config) {
		if logLevel != NONE {
			config.Log = neo4j.ConsoleLogger(logLevels[logLevel])
		}
	}); err != nil {
		return nil, err
	}

	return driver, err
}
