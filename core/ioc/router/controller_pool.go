package router

import (
	"reflect"
	"sync"
)

type ControllerPool struct {
	mu                sync.RWMutex
	instanceMap       map[string]reflect.Type
	controllerMap     []string
	controllerFuncMap map[string]string
}
