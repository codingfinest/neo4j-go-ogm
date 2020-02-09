package gogm

import (
	"reflect"
	"strings"
)

var (
	relTypeTag      = "reltype"
	startNodeTag    = "startNode"
	endNodeTag      = "endNode"
	labelsTag       = "label"
	directionTag    = "direction"
	customIDTag     = "id"
	propertyNameTag = "name"
	uniqueTag       = "unique"
	indexTag        = "index"
)

var (
	tpc = tagParserConf{
		TagKey:     "gogm",
		Delim:      ",",
		AssignOp:   ":",
		IgnoreChar: "-",
		MultiPropsAllowed: map[string]bool{
			"label": true,
		},
	}
)

type tagParserConf struct {
	TagKey            string
	IgnoreChar        string
	Delim             string
	AssignOp          string
	MultiPropsAllowed map[string]bool
}

type tag struct {
	keyval       string
	mappedKeyval map[string][]string
}

func getNamespacedTag(fieldTag reflect.StructTag) *tag {
	tag := &tag{
		keyval: emptyString}
	if namespaceTag, ok := fieldTag.Lookup(tpc.TagKey); ok {
		tag.keyval = namespaceTag
	}
	return tag
}

func (t *tag) get(key string) []string {
	if t.mappedKeyval != nil {
		return t.mappedKeyval[key]
	}
	t.mappedKeyval = map[string][]string{}
	keyVals := strings.Split(t.keyval, tpc.Delim)
	for _, v := range keyVals {
		keyVal := strings.Split(v, tpc.AssignOp)
		if len(keyVal) == 2 {
			key := strings.Trim(keyVal[0], spaceString)
			value := strings.Trim(keyVal[1], spaceString)
			t.mappedKeyval[key] = append(t.mappedKeyval[key], value)
		}

		if len(keyVal) == 1 {
			key := strings.Trim(keyVal[0], spaceString)
			t.mappedKeyval[key] = append(t.mappedKeyval[key], emptyString)
		}
	}
	return t.mappedKeyval[key]
}
