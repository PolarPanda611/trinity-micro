package container

import (
	"fmt"
	"reflect"
)

type diSelfCheckResultCount struct {
	info    int
	warning int
}

// DiSelfCheck ()map[reflect.Type]interface{} {}
func (s *Container) DiSelfCheck(instanceType reflect.Type) {
	resultCount := new(diSelfCheckResultCount)
	pool, ok := s.poolMap[instanceType]
	if !ok {
		panic(fmt.Errorf("%v not exist in pool map", instanceType))
	}
	instance := pool.Get()
	defer pool.Put(instance)
	instanceVal := reflect.Indirect(reflect.ValueOf(instance))
	for index := 0; index < instanceVal.NumField(); index++ {
		objectName := encodeObjectName(instance, index)
		availableInjectInstance := 0
		var availableInjectType []reflect.Type
		val := instanceVal.Field(index)
		autoWire := s.getAutoWireTag(instance, index)
		if !autoWire {
			if val.CanSet() {
				fmt.Println(objectName, resultCount, "warning: autowired tag not set and the val can set ")
			}
			continue
		}
		if !val.CanSet() {
			fmt.Println(objectName, resultCount, "error : private param")
			continue
		}
		if !val.IsZero() {
			fmt.Println(objectName, resultCount, "error : not null param")
			continue
		}
		if val.Kind() == reflect.Struct {
			fmt.Println(objectName, resultCount, "error : should be addressable")
			continue
		}
		if val.Kind() == reflect.Ptr {
			for _, v := range s.GetInstanceTypeByTag(s.getResourceTag(instance, index)) {
				if val.Type() == v {
					availableInjectType = append(availableInjectType, v)
					availableInjectInstance++
				}
			}
			fmt.Println(availableInjectInstance, objectName, resultCount, availableInjectType)
			continue
		}
		if val.Kind() == reflect.Interface {
			for _, v := range s.GetInstanceTypeByTag(s.getResourceTag(instance, index)) {
				if v.Implements(val.Type()) {
					availableInjectType = append(availableInjectType, v)
					availableInjectInstance++
				}
			}
			fmt.Println(availableInjectInstance, objectName, resultCount, availableInjectType)
			continue
		}
	}
}

func (s *Container) DiAllFields(dest interface{}, injectingMap map[reflect.Type]interface{}) (map[reflect.Type]interface{}, map[reflect.Type]interface{}) {
	t := reflect.TypeOf(dest)
	if t.Kind() != reflect.Ptr {
		panic(fmt.Errorf("toFree object %v should be addressable", t))
	}
	sharedInstance := make(map[reflect.Type]interface{})
	toFreeInstance := make(map[reflect.Type]interface{})
	destVal := reflect.Indirect(reflect.ValueOf(dest))
	for index := 0; index < destVal.NumField(); index++ {
		autoWire := s.getAutoWireTag(dest, index)
		autoFree := s.getAutoFreeTag(dest, index)
		val := destVal.Field(index)
		if !autoWire {
			continue
		}
		objectName := encodeObjectName(dest, index)
		switch s.instanceMapping[objectName] {
		default:
			repo, sharedInstanceMap, toFreeInstanceMap := s.GetInstance(s.instanceMapping[objectName], injectingMap)
			for instanceType, instanceValue := range sharedInstanceMap {
				sharedInstance[instanceType] = instanceValue
			}
			for instanceType, instanceValue := range toFreeInstanceMap {
				toFreeInstance[instanceType] = instanceValue
			}
			val.Set(reflect.ValueOf(repo))
			sharedInstance[val.Type()] = repo
			if autoFree {
				toFreeInstance[val.Type()] = repo
			}
			break
		}
	}
	return sharedInstance, toFreeInstance
}

func (s *Container) DiFree(dest interface{}) {
	t := reflect.TypeOf(dest)
	if t.Kind() != reflect.Ptr {
		panic(fmt.Errorf("toFree object %v should be addressable", t))
	}
	destVal := reflect.Indirect(reflect.ValueOf(dest))
	for index := 0; index < destVal.NumField(); index++ {
		val := destVal.Field(index)
		autoWire := s.getAutoWireTag(dest, index)
		autoFree := s.getAutoFreeTag(dest, index)
		if !val.CanSet() {
			continue
		}
		if !autoWire {
			continue
		}
		if !autoFree {
			continue
		}
		if !val.IsZero() {
			val.Set(reflect.Zero(val.Type()))
		}
	}
}
