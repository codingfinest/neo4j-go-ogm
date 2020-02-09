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
