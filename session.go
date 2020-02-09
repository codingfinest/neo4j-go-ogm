package gogm

import "github.com/neo4j/neo4j-go-driver/neo4j"

type Session interface {
	Load(object interface{}, ID interface{}, loadOptions *LoadOptions) error
	LoadAll(objects interface{}, IDs interface{}, loadOptions *LoadOptions) error
	Reload(objects ...interface{}) error
	Save(objects interface{}, saveOptions *SaveOptions) error

	//Delete object(s) at depth depth
	Delete(object interface{}) error

	//Delete all entities of object type
	DeleteAll(object interface{}, deleteOptions *DeleteOptions) error

	PurgeDatabase() error
	Clear() error
	BeginTransaction() (*transaction, error)
	GetTransaction() *transaction
	QueryForObject(object interface{}, cypher string, parameters map[string]interface{}) error
	QueryForObjects(objects interface{}, cypher string, parameters map[string]interface{}) error
	Query(cypher string, parameters map[string]interface{}) (neo4j.Result, error)
	CountEntitiesOfType(object interface{}) (int64, error)
	Count(cypher string, parameters map[string]interface{}) (int64, error)
	RegisterEventListener(EventListener) error
	DisposeEventListener(EventListener) error
}
