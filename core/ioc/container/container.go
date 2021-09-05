package container

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

type Config struct {
	// AutoWired
	// the default value of autowired
	AutoWired bool
	Log       logrus.FieldLogger
}

var (
	DefaultConfig = &Config{
		AutoWired: true,
		Log:       logrus.New(),
	}
)

type Container struct {
	c       *Config
	poolMap map[string]*sync.Pool
}

// NewContainer new pool with init map
func NewContainer(c ...Config) *Container {
	result := new(Container)
	result.poolMap = make(map[string]*sync.Pool)
	if len(c) > 0 {
		result.c = &c[0]
	} else {
		result.c = DefaultConfig
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
	s.poolMap[instanceName] = instancePool
}

// CheckInstanceNameIfExist check contain name if exist
func (s *Container) CheckInstanceNameIfExist(instanceName string) bool {
	_, ok := s.poolMap[instanceName]
	return ok
}

// InstanceDISelfCheck  self check di request registered func exist or not
func (s *Container) InstanceDISelfCheck() error {
	for k := range s.poolMap {
		if err := s.DiSelfCheck(k); err != nil {
			s.c.Log.Errorf("instance self check passed failed => %v, error: %v", k, err)
			return err
		}
		s.c.Log.Infof("instance self check passed => %v", "di self check", k)
	}
	return nil
}

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

func (s *Container) Release(instanceName string, instance interface{}) {
	instancePool, ok := s.poolMap[instanceName]
	if !ok {
		s.c.Log.Errorf("instance release failed => %v, not exist in container", instanceName)
		return
	}
	DiFree(s.c.Log, instance)
	instancePool.Put(instance)
}

func (s *Container) getAutoWireTag(obj interface{}, index int) bool {
	v, exist := getBoolTagFromContainer(obj, index, AUTOWIRE)
	if exist {
		return v
	}
	return s.c.AutoWired
}
