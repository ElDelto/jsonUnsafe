package jsonUnsafe

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

func MarshalJSON(input interface{}) ([]byte, error) {
	//m := map[string]interface{}{}

	return nil, nil
}

func UnmarshalJSON(data []byte, target interface{}) error {
	// TODO: if target is not a struct return standard json.Unmarshal

	m := map[string]interface{}{}
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	rtValue := reflect.ValueOf(target).Elem()

	for i, numField := 0, rtValue.Type().NumField(); i < numField; i++ {
		fieldKey := rtValue.Type().Field(i).Name

		mapValue, exists := findValue(m, fieldKey)
		if !exists {
			return fmt.Errorf("key '%v' does not exist in the given map", fieldKey)
		}

		fieldValue := rtValue.Field(i)
		err := setUnexportedField(&fieldValue, mapValue)
		if err != nil {
			return err
		}
	}

	return err
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
