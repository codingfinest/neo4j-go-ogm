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
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type sessionImpl struct {
	cypherExecuter *cypherExecuter
	saver          *saver
	loader         *loader
	deleter        *deleter
	queryer        *queryer
	transactioner  *transactioner
	store          store
	registry       *registry
	driver         neo4j.Driver
	eventer        *eventer
}

func (s *sessionImpl) Load(object interface{}, ID interface{}, loadOptions *LoadOptions) error {
	_, err := s.loader.load(object, ID, loadOptions, false)
	return err
}

func (s *sessionImpl) LoadAll(objects interface{}, IDs interface{}, loadOptions *LoadOptions) error {
	return s.loader.loadAll(objects, IDs, loadOptions)
}

func (s *sessionImpl) Reload(objects ...interface{}) error {
	return s.loader.reload(objects...)
}

func (s *sessionImpl) Save(objects interface{}, saveOptions *SaveOptions) error {
	return s.saver.save(objects, saveOptions)
}

func (s *sessionImpl) Delete(object interface{}) error {
	return s.deleter.delete(object)
}

func (s *sessionImpl) DeleteAll(objects interface{}, deleteOptions *DeleteOptions) error {
	return s.deleter.deleteAll(objects, deleteOptions)
}

func (s *sessionImpl) PurgeDatabase() error {
	var err error
	if err = s.deleter.purgeDatabase(); err != nil {
		return err
	}
	return s.store.clear()
}

func (s *sessionImpl) Clear() error {
	return s.store.clear()
}

func (s *sessionImpl) BeginTransaction() (*transaction, error) {
	return s.transactioner.beginTransaction(s)
}

func (s *sessionImpl) GetTransaction() *transaction {
	return s.transactioner.transaction
}

//Precondition:
// * object is a pointer to a pointer of domain object: **<domainObject>
// * cypher returns one record with a column of domain object(s)
// * database entity type - node/relationhip - returned by cypher matches the domain object type - node/relationship
// * it is the user's resposibility to make sure  the database object returned by cypher are unloadable into domain object
//
//Post condition:
//Polulated domain objects
func (s *sessionImpl) QueryForObject(object interface{}, cypher string, parameters map[string]interface{}) error {
	return s.queryer.queryForObject(object, cypher, parameters)
}

//Precondition:
// * objects is a pointer to slice of pointers to domain objects: *[]*<domainObject>
// * cypher returns one or more record(s) with a column of domain object(s)
// * database entity type - node/relationhip - returned by cypher matches the domain object type -node/relationship
// * it is the user's resposibility to make sure that database objects returned by cypher are unloadable into the domain object
//
//Post condition:
//Polulated domain objects
func (s *sessionImpl) QueryForObjects(objects interface{}, cypher string, parameters map[string]interface{}) error {
	return s.queryer.queryForObjects(objects, cypher, parameters)
}

func (s *sessionImpl) Query(cypher string, parameters map[string]interface{}, objects ...interface{}) ([]map[string]interface{}, error) {
	return s.queryer.query(cypher, parameters, objects...)
}

func (s *sessionImpl) CountEntitiesOfType(object interface{}) (int64, error) {
	return s.queryer.countEntitiesOfType(object)
}

func (s *sessionImpl) Count(cypher string, parameters map[string]interface{}) (int64, error) {
	return s.queryer.count(cypher, parameters)
}

func (s *sessionImpl) RegisterEventListener(eventListener EventListener) error {
	return s.eventer.registerEventListener(eventListener)
}
func (s *sessionImpl) DisposeEventListener(eventListener EventListener) error {
	return s.eventer.disposeEventListener(eventListener)
}
