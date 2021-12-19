package container

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/sirupsen/logrus"
)

type Config struct {
	// AutoWired
	// the default value of autowired
	AutoWire         bool
	Log              logrus.FieldLogger
	ContainerKeyword Keyword
	AutowireKeyword  Keyword
	ResourceKeyword  Keyword
}

var (
	_DefaultConfig = &Config{
		AutoWire:         true,
		Log:              logrus.New(),
		ContainerKeyword: _CONTAINER,
		AutowireKeyword:  _AUTOWIRE,
		ResourceKeyword:  _RESOURCE,
	}
)

type Container struct {
	c       *Config
	poolMap map[string]*sync.Pool
	// poolTypeMap caching the type info
	poolTypeMap map[string]reflect.Type
}

// NewContainer get the new container instance
// if not passing the config , will init with the default config
func NewContainer(c ...Config) *Container {
	result := new(Container)
	result.poolMap = make(map[string]*sync.Pool)
	result.poolTypeMap = make(map[string]reflect.Type)
	if len(c) > 0 {
		if c[0].ContainerKeyword == "" {
			c[0].ContainerKeyword = _CONTAINER
		}
		if c[0].AutowireKeyword == "" {
			c[0].AutowireKeyword = _AUTOWIRE
		}
		if c[0].ResourceKeyword == "" {
			c[0].ResourceKeyword = _RESOURCE
		}
		result.c = &c[0]
	} else {
		result.c = _DefaultConfig
	}
	result.c.Log = result.c.Log.WithField("app", "container")
	return result
}

// RegisterInstance register new instance
// if instanceName is empty will fatal
// if instancePool is invalid , will fatal
func (s *Container) RegisterInstance(instanceName string, instancePool *sync.Pool) {
	if instanceName == "" {
		s.c.Log.Fatal("instance name cannot be empty")
	}
	if instancePool == nil {
		s.c.Log.Fatal("instance pool cannot be empty")
	}
	if _, ok := s.poolMap[instanceName]; ok {
		s.c.Log.Fatal(fmt.Errorf("instance name %v already existed", instanceName))
	}
	ins := instancePool.Get()
	defer instancePool.Put(ins)
	t := reflect.TypeOf(ins)
	s.poolMap[instanceName] = instancePool
	s.poolTypeMap[instanceName] = t
}

// CheckInstanceNameIfExist
// check instance name if exist
// if exist , return true
// if not exist , return false
func (s *Container) CheckInstanceNameIfExist(instanceName string) bool {
	_, ok := s.poolMap[instanceName]
	return ok
}

// InstanceDISelfCheck
// self check all the instance registered exist or not
func (s *Container) InstanceDISelfCheck() error {
	for k := range s.poolMap {
		if err := s.DiSelfCheck(k); err != nil {
			s.c.Log.Errorf("%-8v %-10v %-7v => %v, error: %v", "instance", "self-check", "failed", k, err)
			return err
		}
		s.c.Log.Infof("%-8v %-10v %-7v => %v", "instance", "self-check", "success", k)
	}
	return nil
}

// InstanceDISelfCheck
// get instance by instance name
// injectingMap , the dependency instance, will inject the instance in injectingMap as priority
func (s *Container) GetInstance(instanceName string, injectingMap map[string]interface{}) interface{} {
	if v, ok := injectingMap[instanceName]; ok {
		return v
	}
	pool, ok := s.poolMap[instanceName]
	if !ok {
		s.c.Log.Fatalf("instance not exist in container => %v", instanceName)
	}
	service := pool.Get()
	injectingMap[instanceName] = service
	s.DiAllFields(service, injectingMap)
	return service
}

// Release
// release the instance to instance pool
func (s *Container) Release(instanceName string, instance interface{}) {
	instancePool, ok := s.poolMap[instanceName]
	if !ok {
		s.c.Log.Errorf("instance release failed => %v, not exist in container", instanceName)
		return
	}
	if reflect.TypeOf(instance) != s.poolTypeMap[instanceName] {
		panic("")
	}
	s.DiFree(instance)
	instancePool.Put(instance)
}

func (s *Container) getAutoWireTag(obj interface{}, index int) bool {
	v, exist := getBoolTagFromContainer(obj, index, s.c.ContainerKeyword, s.c.AutowireKeyword)
	if exist {
		return v
	}
	return s.c.AutoWire
}
