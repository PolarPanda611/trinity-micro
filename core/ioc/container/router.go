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
			injectMap := make(map[string]interface{})
			instance := container.GetInstance(controller.instanceName, injectMap)
			defer func() {
				for k, v := range injectMap {
					container.Release(k, v)
				}
			}()
			currentMethod, ok := reflect.TypeOf(instance).MethodByName(requestMap.funcName)
			if !ok {
				panic(fmt.Sprintf("instance %v method %v not exist ", controller.instanceName, currentMethod))
			}

		}
	}

}

func DIHandler(container *Container, instanceName string, funcName string) func(w http.ResponseWriter, r *http.Request) {
	injectMapping := make(map[string]interface{})
	instance := container.GetInstance(instanceName, injectMapping)
	currentMethod, ok := reflect.TypeOf(instance).MethodByName(funcName)
	if !ok {
		panic("controller has no method ")
	}
	return httpx.DIParamHandler(currentMethod.Func.Interface())
}

type Mux interface {
	Method(method, pattern string, handler http.Handler)
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
			log.Printf("router register: method: %v , path: %v handler: %v.%v ", requestMapping.method, urlPath, controller.instanceName, requestMapping.funcName)
		}
	}

}
