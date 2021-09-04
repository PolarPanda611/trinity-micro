package container

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
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

func getStringTagFromContainer(obj interface{}, index int, tagName Keyword) (value string, isExist bool) {
	v, exist := getTagByName(obj, index, CONTAINER)
	if exist {
		resourceValue, ok := decodeTag(v, tagName)
		if ok {
			return resourceValue, true
		}
	}
	return "", false
}

// getBoolTag
func getBoolTagFromContainer(obj interface{}, index int, tagName Keyword) (value bool, isExist bool) {
	v, exist := getTagByName(obj, index, CONTAINER)
	if exist {
		autoFreeOption, ok := decodeTag(v, tagName)
		if ok {
			if autoFreeOption == "" {
				return true, true
			}
			b, _ := strconv.ParseBool(autoFreeOption)
			return b, true
		}
	}
	return false, false
}

func decodeTag(value string, key Keyword) (string, bool) {
	kvStr := strings.Split(strings.Trim(value, TAG_SPLITER), TAG_SPLITER)
	t := make(map[string]string)
	for _, v := range kvStr {
		if v == "" {
			continue
		}
		index := strings.Index(v, string(TAG_KV_SPLITER))
		if index == 0 {
			continue
		} else if index >= 0 {
			t[v[:index]] = v[index+1:]
		} else {
			t[v] = ""
		}

	}
	v, ok := t[string(key)]
	return v, ok
}
