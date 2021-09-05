package container

import (
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
)

// DiSelfCheck ()map[reflect.Type]interface{} {}
func (s *Container) DiSelfCheck(instanceName string) error {
	pool, ok := s.poolMap[instanceName]
	if !ok {
		return fmt.Errorf("instance `%v` not exist in pool map", instanceName)
	}
	instance := pool.Get()
	defer pool.Put(instance)
	t := reflect.TypeOf(instance)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("the object to be injected %v should be addressable", t)
	}
	instanceVal := reflect.Indirect(reflect.ValueOf(instance))
	for index := 0; index < instanceVal.NumField(); index++ {
		objectName := encodeObjectName(instance, index)
		if _, exist := getTagByName(instance, index, CONTAINER); !exist {
			s.c.Log.Debugf("%20v: instanceName: %v index: %v objectName: %v, the container tag not exist, skip inject", "di self check", instanceName, index, objectName)
			continue
		}
		resourceName, exist := getStringTagFromContainer(instance, index, RESOURCE)
		if !exist {
			return fmt.Errorf("self check error: instanceName: %v index: %v objectName: %v, the resource tag not exist in container", instanceName, index, objectName)
		}
		instancePool, exist := s.poolMap[resourceName]
		if !exist {
			return fmt.Errorf("self check error: instanceName: %v index: %v objectName: %v, resource name: %v not register in container ", instanceName, index, objectName, resourceName)
		}
		val := instanceVal.Field(index)
		autoWire := s.getAutoWireTag(instance, index)
		if !autoWire {
			if val.CanSet() {
				s.c.Log.Warnf("self check warning: instanceName: %v index: %v objectName: %v, autowired is false but the param can be injected ", instanceName, index, objectName)
			}
		}
		if !val.CanSet() {
			return fmt.Errorf("self check error: instanceName: %v index: %v objectName: %v, private param", instanceName, index, objectName)
		}
		if !val.IsZero() {
			return fmt.Errorf("self check error: instanceName: %v index: %v objectName: %v, the param to be injected is not null", instanceName, index, objectName)
		}
		switch val.Kind() {
		case reflect.Interface:
			instance := instancePool.Get()
			defer instancePool.Put(instance)
			instanceType := reflect.TypeOf(instance)
			if !instanceType.Implements(val.Type()) {
				return fmt.Errorf("self check error: instanceName: %v index: %v objectName: %v, resource name: %v type: %v not implement the interface %v", instanceName, index, objectName, resourceName, instanceType.Name(), val.Type().Name())
			}
		default:
			instance := instancePool.Get()
			defer instancePool.Put(instance)
			instanceType := reflect.TypeOf(instance)
			if val.Type() != instanceType {
				return fmt.Errorf("self check error: instanceName: %v index: %v objectName: %v, resource name: %v type not same, expected: %v actual: %v", instanceName, index, objectName, resourceName, val.Type(), instanceType)
			}
		}
	}
	return nil
}

func (s *Container) DiAllFields(dest interface{}, injectingMap map[string]interface{}) {
	destVal := reflect.Indirect(reflect.ValueOf(dest))
	for index := 0; index < destVal.NumField(); index++ {
		if _, exist := getTagByName(dest, index, CONTAINER); !exist {
			continue
		}
		resourceName, _ := getStringTagFromContainer(dest, index, RESOURCE)
		val := destVal.Field(index)
		if instance, exist := injectingMap[resourceName]; exist {
			val.Set(reflect.ValueOf(instance))
			continue
		}
		instance := s.GetInstance(resourceName, injectingMap)
		val.Set(reflect.ValueOf(instance))
	}
}

func DiFree(log logrus.FieldLogger, dest interface{}) {
	t := reflect.TypeOf(dest)
	switch t.Kind() {
	case reflect.Ptr:
		destVal := reflect.Indirect(reflect.ValueOf(dest))
		for index := 0; index < destVal.NumField(); index++ {
			objectName := encodeObjectName(dest, index)
			if _, exist := getTagByName(dest, index, CONTAINER); !exist {
				log.Debugf("objectName di free skipped => %v, container not exist", objectName)
				continue
			}
			val := destVal.Field(index)
			if val.CanSet() {
				val.Set(reflect.Zero(val.Type()))
			}
		}
	}
}
