package router

import (
	"net/http"
	"reflect"
	"sync"
)

var (
	_bootingControllers []bootingController
	_bootingInstances   []bootingInstance
	_bootingModels      []bootingModel
)

type bootingController struct {
	path            string
	instance        interface{}
	requestMappings []RequestMap
}

type bootingInstance struct {
	instanceName reflect.Type
	instancePool *sync.Pool
	instanceTags []string
}

type bootingModel struct {
	modelInstance interface{}
	defaultValues []interface{}
}

type RequestMap struct {
	Method   string
	SubPath  string
	FuncName string
	Handlers []func(http.Handler) http.Handler
}

type RouterRegister interface {
	Method(method, pattern string, handler http.Handler)
}

func RegisterController(path string, instance interface{}, requestMappings ...RequestMap) {
	newController := bootingController{
		path:            path,
		instance:        instance,
		requestMappings: requestMappings,
	}
	_bootingControllers = append(_bootingControllers, newController)
}

func NewRequestMapping(method string, subPath string, funcName string, handlers ...func(http.Handler) http.Handler) RequestMap {
	return RequestMap{
		Method:   method,
		SubPath:  subPath,
		FuncName: funcName,
		Handlers: handlers,
	}
}

func initController() {

}

func RegisterRouter(r RouterRegister) {
	// for _, controller := range _bootingControllers {
	// 	for _, mapping := range controller.requestMappings {
	// 		 r.Method(mapping.Method, controller.path+mapping.SubPath)
	// 	}
	// }

}
