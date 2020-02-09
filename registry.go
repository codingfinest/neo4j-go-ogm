package gogm

import (
	"reflect"
	"sync"
)

type registry struct {
	objects        map[string]metadata
	cypherExecuter cypherExecuter
	objectsMu      sync.Mutex
}

func newRegistry(cypherExecuter cypherExecuter) *registry {
	objects := map[string]metadata{}
	return &registry{objects, cypherExecuter, sync.Mutex{}}
}

func (r *registry) register(name string, m metadata) {
	r.objects[name] = m
}

func (r *registry) get(t reflect.Type) (metadata, error) {
	var err error
	if r.getMetadata(t.String()) == nil {
		var m metadata
		if m, err = getMetadata(t, r); err != nil {
			return nil, err
		}
		r.setMetadata(t.String(), m)
		for _, statement := range getCreateSchemaStatement(r.objects[t.String()]) {
			if _, err = r.cypherExecuter.exec(statement, nil); err != nil {
				return nil, err
			}
		}
	}
	return r.objects[t.String()], err
}

func (r *registry) getMetadata(id string) metadata {
	r.objectsMu.Lock()
	defer r.objectsMu.Unlock()
	return r.objects[id]
}

func (r *registry) setMetadata(id string, metadata metadata) {
	r.objectsMu.Lock()
	defer r.objectsMu.Unlock()
	r.objects[id] = metadata
}
