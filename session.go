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

//Session provides access to the database
type Session interface {
	Load(object interface{}, ID interface{}, loadOptions *LoadOptions) error
	LoadAll(objects interface{}, IDs interface{}, loadOptions *LoadOptions) error
	Reload(objects ...interface{}) error
	Save(objects interface{}, saveOptions *SaveOptions) error
	Delete(object interface{}) error
	DeleteAll(object interface{}, deleteOptions *DeleteOptions) error
	PurgeDatabase() error
	Clear() error
	BeginTransaction() (*transaction, error)
	GetTransaction() *transaction
	QueryForObject(object interface{}, cypher string, parameters map[string]interface{}) error
	QueryForObjects(objects interface{}, cypher string, parameters map[string]interface{}) error
	Query(cypher string, parameters map[string]interface{}, objects ...interface{}) ([]map[string]interface{}, error)
	CountEntitiesOfType(object interface{}) (int64, error)
	Count(cypher string, parameters map[string]interface{}) (int64, error)
	RegisterEventListener(EventListener) error
	DisposeEventListener(EventListener) error
}
