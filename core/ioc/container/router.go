package container

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"reflect"
	"sync"
	"trinity-micro/core/httpx"
)

var (
	injectMapPool = &sync.Pool{
		New: func() interface{} {
			return make(map[string]interface{})
		},
	}
)

type bootingController struct {
	rootPath     string
	instanceName string
	requestMaps  []RequestMap
}

type bootingInstance struct {
	instanceName string
	instancePool *sync.Pool
}

var (
	_bootingControllers []bootingController
	_bootingInstances   []bootingInstance
)

type RequestMap struct {
	method   string
	subPath  string
	funcName string
	handlers []http.Handler
}

func RegisterController(rootPath string, instanceName string, requestMaps ...RequestMap) {
	newController := bootingController{
		rootPath:     rootPath,
		instanceName: instanceName,
		requestMaps:  requestMaps,
	}
	_bootingControllers = append(_bootingControllers, newController)
}

func RegisterInstance(instanceName string, instancePool *sync.Pool) {
	newInstance := bootingInstance{
		instanceName: instanceName,
		instancePool: instancePool,
	}
	_bootingInstances = append(_bootingInstances, newInstance)
}

func NewRequestMapping(method string, path string, funcName string, handlers ...http.Handler) RequestMap {
	return RequestMap{
		method:   method,
		subPath:  path,
		funcName: funcName,
		handlers: handlers,
	}
}

func InitInstance(container *Container) {
	for _, instance := range _bootingInstances {
		log.Printf("register instance: instanceName: %v ", instance.instanceName)
		container.RegisterInstance(instance.instanceName, instance.instancePool)
	}
	container.InstanceDISelfCheck()
}

func RouterSelfCheck(container *Container) {
	for _, controller := range _bootingControllers {
		for _, requestMap := range controller.requestMaps {
			injectMap := injectMapPool.Get().(map[string]interface{})
			instance := container.GetInstance(controller.instanceName, injectMap)
			defer func() {
				for k, v := range injectMap {
					container.Release(k, v)
					delete(injectMap, k)
				}
				injectMapPool.Put(injectMap)
			}()
			currentMethod, ok := reflect.TypeOf(instance).MethodByName(requestMap.funcName)
			if !ok {
				panic(fmt.Sprintf("instance %v method %v not exist ", controller.instanceName, currentMethod))
			}

		}
	}

}

func DIHandler(container *Container, instanceName string, funcName string) func(w http.ResponseWriter, r *http.Request) {
	injectMap := injectMapPool.Get().(map[string]interface{})
	instance := container.GetInstance(instanceName, injectMap)
	defer func() {
		for k, v := range injectMap {
			container.Release(k, v)
			delete(injectMap, k)
		}
		injectMapPool.Put(injectMap)
	}()
	currentMethod, _ := reflect.TypeOf(instance).MethodByName(funcName)
	return httpx.DIParamMethod(currentMethod, instance)
}

type Mux interface {
	MethodFunc(method, pattern string, handlerFn http.HandlerFunc)
}

func DIRouter(r Mux, container *Container) {
	InitInstance(container)
	RouterSelfCheck(container)
	// register router
	for _, controller := range _bootingControllers {
		for _, requestMapping := range controller.requestMaps {
			urlPath := filepath.Join(controller.rootPath, requestMapping.subPath)
			r.MethodFunc(requestMapping.method, urlPath, DIHandler(container, controller.instanceName, requestMapping.funcName))
			log.Printf("request mapping: method: %-6s %-30s => handler: %v.%v ", requestMapping.method, urlPath, controller.instanceName, requestMapping.funcName)
		}
	}

}
