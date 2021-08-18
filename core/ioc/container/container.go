package container

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type Config struct {
	AutoFree  bool
	AutoWired bool
}

var (
	DefaultConfig = &Config{
		AutoFree:  true,
		AutoWired: false,
	}
)

type Container struct {
	c *Config
	// map[instanceName] = instancePool
	poolMap map[reflect.Type]*sync.Pool
	// instance list
	instanceTypeList []reflect.Type
	// instance tags maps instanceName
	poolTags map[string]reflect.Type
	// instanceMapping instance mapping
	// caching the instance di instance relation during the self check
	instanceMapping map[string]reflect.Type
}

// NewContainer new pool with init map
func NewContainer(c ...Config) *Container {
	result := new(Container)
	result.poolMap = make(map[reflect.Type]*sync.Pool)
	result.poolTags = make(map[string]reflect.Type)
	result.instanceMapping = make(map[string]reflect.Type)
	if len(c) > 0 {
		result.c = &c[0]
	} else {
		result.c = DefaultConfig
	}
	return result
}

// RegisterInstance register new instance
func (s *Container) RegisterInstance(instance interface{}, instanceTag ...string) {
	instanceType := reflect.TypeOf(instance)
	switch instanceType.Kind() {
	case reflect.Struct:
		s.newInstance(instanceType, &sync.Pool{
			New: func() interface{} {
				return reflect.New(reflect.TypeOf(instance)).Interface()
			},
		}, instanceTag...)
	case reflect.Func:
		f, ok := instance.(func() interface{})
		if !ok {
			log.Fatal("The instance func should be  func () interface{}")
		}
		s.newInstance(reflect.TypeOf(f()), &sync.Pool{
			New: func() interface{} {
				return f()
			},
		}, instanceTag...)
	case reflect.Ptr:
		s.RegisterInstance(reflect.Indirect(reflect.ValueOf(instance)).Interface())
	default:
		panic("The instance should be struct or func () interface{}")
	}

}

// newInstance new instance
func (s *Container) newInstance(instanceType reflect.Type, instancePool *sync.Pool, instanceTag ...string) {
	if _, ok := s.poolMap[instanceType]; ok {
		return
	}
	s.poolMap[instanceType] = instancePool
	s.instanceTypeList = append(s.instanceTypeList, instanceType)
	if len(instanceTag) > 0 {
		if instanceTag[0] != "" {
			if _, ok := s.poolTags[instanceTag[0]]; ok {
				panic(fmt.Errorf("tag %v already existed", instanceTag[0]))
			}
			s.poolTags[instanceTag[0]] = instanceType
		}
	}
}

// GetInstanceType get all service type by tag
// if no tag provide , return all type
// if tags provide , will return the types of the tags
func (s *Container) GetInstanceTypeByTag(tags ...string) []reflect.Type {
	if len(tags) > 0 {
		var types []reflect.Type
		for _, v := range tags {
			if instance, ok := s.poolTags[v]; ok {
				types = append(types, instance)
			}
		}
		return types
	}
	return s.instanceTypeList
}

// CheckInstanceNameIfExist check contain name if exist
func (s *Container) CheckInstanceNameIfExist(instanceName reflect.Type) bool {
	_, ok := s.poolMap[instanceName]
	return ok
}

// InstanceMapping get instance mapping
// return the copy of the instance mapping
func (s *Container) InstanceMapping() map[string]reflect.Type {
	instanceMap := make(map[string]reflect.Type, len(s.instanceMapping))
	for k, v := range s.instanceMapping {
		instanceMap[k] = v
	}
	return instanceMap
}

// InstanceDISelfCheck  self check di request registered func exist or not
func (s *Container) InstanceDISelfCheck() {
	for _, v := range s.instanceTypeList {
		s.DiSelfCheck(v)
	}
	return

}

func (s *Container) GetInstance(instanceType reflect.Type, injectingMap map[reflect.Type]interface{}) (interface{}, map[reflect.Type]interface{}, map[reflect.Type]interface{}) {
	if v, ok := injectingMap[instanceType]; ok {
		return v, nil, nil
	}
	pool, ok := s.poolMap[instanceType]
	if !ok {
		panic("unknown service name")
	}
	service := pool.Get()
	injectingMap[instanceType] = service
	sharedInstance, toFreeInstance := s.DiAllFields(service, injectingMap)
	return service, sharedInstance, toFreeInstance
}

func (s *Container) Release(instance interface{}) {
	t := reflect.TypeOf(instance)
	syncpool, ok := s.poolMap[t]
	if !ok {
		return
	}
	s.DiFree(instance)
	syncpool.Put(instance)
}
func (s *Container) getResourceTag(obj interface{}, index int) string {
	v, exist := getTagByName(obj, index, CONTAINER)
	if exist {
		resourceValue, ok := decodeTag(v, RESOURCE)
		if ok {
			return resourceValue
		}
	}
	return ""
}
func (s *Container) getAutoWireTag(obj interface{}, index int) bool {
	v, exist := getTagByName(obj, index, CONTAINER)
	if exist {
		autoWireOption, ok := decodeTag(v, AUTOWIRE)
		if ok {
			if autoWireOption == "" {
				return true
			}
			b, _ := strconv.ParseBool(autoWireOption)
			return b
		}
	}
	return s.c.AutoWired
}

func (s *Container) getAutoFreeTag(obj interface{}, index int) bool {
	v, exist := getTagByName(obj, index, CONTAINER)
	if exist {
		autoFreeOption, ok := decodeTag(v, AUTOFREE)
		if ok {
			if autoFreeOption == "" {
				return true
			}
			b, _ := strconv.ParseBool(autoFreeOption)
			return b
		}
	}
	return s.c.AutoFree
}

func decodeTag(value string, key Keyword) (string, bool) {
	kvStr := strings.Split(strings.Trim(value, ";"), TAG_SPLITER)
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
