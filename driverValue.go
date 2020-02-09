package gogm

import (
	"reflect"
)

func driverValueAsType(driverValue interface{}, structFieldType reflect.Type) interface{} {

	switch driverValue.(type) {
	case []interface{}:
		return sliceAsType(driverValue, structFieldType)
	case int64:
		return int64AsType(driverValue, structFieldType)
	case float64:
		return float64AsType(driverValue, structFieldType)
	case []byte:
		return byteAsType(driverValue, structFieldType)
	default:
		return valueAsType(driverValue, structFieldType)
	}
}

func sliceAsType(driverValue interface{}, structFieldType reflect.Type) interface{} {
	switch structFieldType.Kind() {
	case reflect.Slice:
		values := reflect.ValueOf(driverValue)
		slice := reflect.MakeSlice(structFieldType, 0, 0)
		ptr := reflect.New(slice.Type())
		ptr.Elem().Set(slice)
		for i := 0; i < values.Len(); i++ {
			ptr.Elem().Set(reflect.Append(ptr.Elem(), reflect.ValueOf(driverValueAsType(values.Index(i).Interface(), structFieldType.Elem()))))
		}
		return ptr.Elem().Interface()
	case reflect.Ptr:
		ptr := reflect.New(structFieldType.Elem())
		ptr.Elem().Set(reflect.ValueOf(sliceAsType(driverValue, structFieldType.Elem())))
		return ptr.Interface()
	default:
		return driverValue
	}
}

func valueAsType(driverValue interface{}, structFieldType reflect.Type) interface{} {
	switch structFieldType.Kind() {
	case reflect.Ptr:
		ptr := reflect.New(structFieldType.Elem())
		ptr.Elem().Set(reflect.ValueOf(valueAsType(driverValue, structFieldType.Elem())))
		return ptr.Interface()
	default:
		return driverValue
	}
}

func int64AsType(driverValue interface{}, structFieldType reflect.Type) interface{} {

	switch structFieldType.Kind() {
	case reflect.Int:
		return int(driverValue.(int64))
	case reflect.Int8:
		return int8(driverValue.(int64))
	case reflect.Int16:
		return int16(driverValue.(int64))
	case reflect.Int32:
		return int32(driverValue.(int64))
	case reflect.Ptr:
		ptr := reflect.New(structFieldType.Elem())
		ptr.Elem().Set(reflect.ValueOf(int64AsType(driverValue, structFieldType.Elem())))
		return ptr.Interface()
	default:
		return driverValue
	}
}

func float64AsType(driverValue interface{}, structFieldType reflect.Type) interface{} {
	switch structFieldType.Kind() {
	case reflect.Float32:
		return float32(driverValue.(float64))
	case reflect.Ptr:
		ptr := reflect.New(structFieldType.Elem())
		ptr.Elem().Set(reflect.ValueOf(float64AsType(driverValue, structFieldType.Elem())))
		return ptr.Interface()
	default:
		return driverValue
	}
}

func byteAsType(driverValue interface{}, structFieldType reflect.Type) interface{} {
	// switch structFieldType.Kind() {

	//TODO complete
	// values := reflect.ValueOf(driverValue)
	// slice := reflect.MakeSlice(structFieldType, 0, 0)
	// ptr := reflect.New(slice.Type())
	// ptr.Elem().Set(slice)
	// for i := 0; i < values.Len(); i++ {
	// 	ptr.Elem().Set(reflect.Append(ptr.Elem(), reflect.ValueOf(driverValueAsType(values.Index(i).Interface(), structFieldType.Elem()))))
	// }
	// return ptr.Elem().Interface()

	// case reflect.Slice:
	// 	value := reflect.ValueOf(driverValue)
	// 	slice := reflect.MakeSlice(structFieldType, 0, 0)
	// 	ptr := reflect.New(slice.Type())
	// 	// ptr.Elem().Elem().Set(byteAsType())
	// 	return ptr.Elem().Interface()
	// case reflect.Ptr:
	// 	ptr := reflect.New(structFieldType.Elem())
	// 	ptr.Elem().Set(reflect.ValueOf(byteAsType(driverValue, structFieldType.Elem())))
	// 	return ptr.Interface()
	// default:
	// 	return driverValue
	// }

	return nil
}
