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
