package container

import (
	"fmt"
	"reflect"
)

func getTagByName(object interface{}, index int, name Keyword) (string, bool) {
	objectType := reflect.TypeOf(object)
	switch objectType.Kind() {
	case reflect.Struct:
		return objectType.Field(index).Tag.Lookup(string(name))
	case reflect.Ptr:
		return objectType.Elem().Field(index).Tag.Lookup(string(name))
	default:
		panic("wrong type , must be struct or ptr ")
	}
}

func encodeObjectName(instance interface{}, index int) string {
	instanceType := reflect.TypeOf(instance)
	instanceVal := reflect.Indirect(reflect.ValueOf(instance))
	paramName := instanceType.Elem().Field(index).Name
	paramType := instanceVal.Field(index).Type()
	return fmt.Sprintf("%v.%v.(%v)", instanceType, paramName, paramType)

}
