package gogm

import (
	"reflect"
)

type field struct {
	parent reflect.Value
	name   string
	tag    *tag
}

func (f *field) getValue() reflect.Value {
	return f.parent.FieldByName(f.name)
}

func (f *field) getStructField() reflect.StructField {
	sf, _ := f.parent.Type().FieldByName(f.name)
	return sf
}

func (f *field) isEntityMetadata(graphEntityType reflect.Type) bool {
	//TODO embedsType check only one layer.support embedded
	return !f.isIgnored() && f.getStructField().Anonymous && (f.getStructField().Type == graphEntityType || embedsType(f.getStructField().Type, graphEntityType))
}

func (f *field) isEntity(graphEntityType reflect.Type) bool {
	if f.isIgnored() || f.getStructField().Anonymous || elem(f.getStructField().Type).Kind() != reflect.Struct {
		return false
	}

	fieldType := elem2(f.getStructField().Type)

	if fieldType.Kind() != reflect.Ptr || fieldType.Elem().Kind() != reflect.Struct {
		return false
	}

	if internalGrpahType, err := getInternalGraphType(elem2(f.getStructField().Type).Elem()); internalGrpahType == graphEntityType && err == nil {
		return true
	}

	return false
}

func (f *field) isTagged(s string) bool {
	return f.tag.get(s) != nil
}

func (f *field) isIgnored() bool {
	return f.isTagged(tpc.IgnoreChar)
}

func (f *field) getEffectiveDirection() direction {
	direction := outgoing
	if taggedDirection := f.tag.get(directionTag); taggedDirection != nil && len(taggedDirection[0]) > 0 {
		direction = directionTags[taggedDirection[0]]
	}
	return direction
}
