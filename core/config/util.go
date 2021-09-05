package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func getParamName(instance interface{}, index int) string {
	instanceType := reflect.TypeOf(instance)
	paramName := instanceType.Elem().Field(index).Name
	return paramName

}
func encodeObjectName(instance interface{}, index int) string {
	instanceType := reflect.TypeOf(instance)
	instanceVal := reflect.Indirect(reflect.ValueOf(instance))
	paramName := getParamName(instance, index)
	paramType := instanceVal.Field(index).Type()
	return fmt.Sprintf("%v.%v.(%v)", instanceType, paramName, paramType)

}

// GetTagsValue get tags value by key
func GetTagsValue(object interface{}, index int, key string) string {
	objectType := reflect.TypeOf(object)
	if objectType.Kind() == reflect.Struct {
		return reflect.TypeOf(object).Field(index).Tag.Get(key)
	}
	return reflect.TypeOf(object).Elem().Field(index).Tag.Get(key)

}

func typeConverter(val string, destType reflect.Type) (interface{}, error) {
	v := reflect.New(destType).Elem()
	switch destType.Kind() {
	case reflect.Ptr:
		return nil, errors.New("unsupported pointer type")
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int:
		convertedValue, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("key %v convert to type : %v  error  %v", val, destType.Kind(), err)
		}
		v.SetInt(convertedValue)
	case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint:
		convertedValue, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("key %v convert to type : %v  error  %v", val, destType.Kind(), err)
		}
		v.SetUint(convertedValue)
	case reflect.Float32, reflect.Float64:
		convertedValue, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, fmt.Errorf("key %v convert to type : %v  error  %v", val, destType.Kind(), err)
		}
		v.SetFloat(convertedValue)
	case reflect.String:
		v.SetString(val)
	case reflect.Bool:
		convertedValue, err := strconv.ParseBool(val)
		if err != nil {
			return nil, fmt.Errorf("key %v convert to type : %v  error  %v", val, destType.Kind(), err)
		}
		v.SetBool(convertedValue)
	case reflect.Struct:
		s := reflect.New(destType).Interface()
		if err := json.Unmarshal([]byte(val), s); err != nil {
			return nil, err
		}
		return reflect.Indirect(reflect.ValueOf(s)).Interface(), nil
	default:
		return nil, fmt.Errorf("type %v not support", destType)
	}
	return v.Interface(), nil

}
