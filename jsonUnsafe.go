// jsonUnsafe uses the "unsafe" package from the standard library to
// encode/decode unexported struct fields.
package jsonUnsafe

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

// Marshal has the same signature as "encoding/json#Marshal". If the given
// input value is not a struct then the result of "encoding/json#Marshal"
// will be returnedl
func Marshal(input interface{}) ([]byte, error) {
	if !isStruct(input) {
		return json.Marshal(&input)
	}

	m := map[string]string{}
	f := func(fieldKey string, fieldValue *reflect.Value) error {
		m[fieldKey] = fieldValue.String()
		return nil
	}

	err := forEachStructField(input, f)
	if err != nil {
		return nil, err
	}

	output, err := json.Marshal(&m)
	return output, err
}

// Unmarshal has the samge signature as "encoding/json#Unmarshal". if the given
// target value is not a struct then the result of "encoding/json#Unmarshal"
// will be returned.
func Unmarshal(data []byte, target interface{}) error {
	if !isStruct(target) {
		return json.Unmarshal(data, &target)
	}

	m := map[string]interface{}{}
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	f := func(fieldKey string, fieldValue *reflect.Value) error {
		mapValue, exists := findValue(m, fieldKey)
		if !exists {
			return fmt.Errorf("key '%v' does not exist in the given map", fieldKey)
		}

		return setUnexportedField(fieldValue, mapValue)
	}

	return forEachStructField(target, f)
}

func findValue(m map[string]interface{}, key string) (interface{}, bool) {
	value, exists := m[key]
	if exists {
		return value, exists
	}

	key = strings.ToLower(key)
	for k, v := range m {
		if key == strings.ToLower(k) {
			return v, true
		}
	}

	return nil, false
}

func setUnexportedField(field *reflect.Value, value interface{}) error {
	rf := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()

	fieldType := field.Type()
	valueType := reflect.ValueOf(value).Type()
	if fieldType != valueType {
		return fmt.Errorf("cannot assign %v to %v", valueType, fieldType)
	}

	rf.Set(reflect.ValueOf(value))

	return nil
}

func isStruct(s interface{}) bool {
	return reflect.TypeOf(s).Kind() == reflect.Ptr &&
		reflect.TypeOf(reflect.ValueOf(s).Elem()).Kind() == reflect.Struct
}

func forEachStructField(s interface{}, f func(fieldKey string, fieldValue *reflect.Value) error) error {
	if !isStruct(s) {
		return fmt.Errorf("%v is not of type struct", s)
	}
	rtValue := reflect.ValueOf(s).Elem()

	for i, numField := 0, rtValue.Type().NumField(); i < numField; i++ {
		fieldKey := rtValue.Type().Field(i).Name
		fieldValue := rtValue.Field(i)

		err := f(fieldKey, &fieldValue)
		if err != nil {
			return err
		}
	}

	return nil
}
