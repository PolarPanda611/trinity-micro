package container

import (
	"context"
	"net/http"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/PolarPanda611/trinity-micro/core/e"
	"github.com/PolarPanda611/trinity-micro/core/httpx"
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
	isRaw    bool
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
		isRaw:    false,
	}
}

func NewRawRequestMapping(method string, path string, funcName string, handlers ...http.Handler) RequestMap {
	return RequestMap{
		method:   method,
		subPath:  path,
		funcName: funcName,
		handlers: handlers,
		isRaw:    true,
	}
}

func InitInstance(container *Container) {
	for _, instance := range _bootingInstances {
		container.RegisterInstance(instance.instanceName, instance.instancePool)
		container.c.Log.Infof("instance registered => %v ", instance.instanceName)
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
			_, ok := reflect.TypeOf(instance).MethodByName(requestMap.funcName)
			if !ok {
				container.c.Log.Fatalf("instance router self check failed => %v.%v , func %v not exist ", controller.instanceName, requestMap.funcName, requestMap.funcName)
				continue
			}
			container.c.Log.Infof("instance router self check passed => %v.%v ", controller.instanceName, requestMap.funcName)
		}
	}

}

// multi instance di handler
func DIHandler(container *Container, instanceName string, funcName string, isRaw bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), httpx.HTTPRequest, r))
		r = r.WithContext(context.WithValue(r.Context(), httpx.HTTPStatus, new(int)))
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
		inParams, err := httpx.InvokeMethod(currentMethod.Type, r, instance, w)
		if err != nil {
			// e.Logging(sessionLogger, err)
			httpx.HttpResponseErr(w, err)
			return
		}
		responseValue := currentMethod.Func.Call(inParams)
		if isRaw {
			return
		}
		switch len(responseValue) {
		case 0:
			httpx.HttpResponse(w, httpx.GetHTTPStatusCode(r.Context(), httpx.DefaultHttpSuccessCode), nil)
			return
		case 1:
			if err, ok := responseValue[0].Interface().(error); ok {
				if err != nil {
					// e.Logging(sessionLogger, err)
					httpx.HttpResponseErr(w, err)
					return
				}
			}
			httpx.HttpResponse(w, httpx.GetHTTPStatusCode(r.Context(), httpx.DefaultHttpSuccessCode), responseValue[0].Interface())
			return
		case 2:
			if err, ok := responseValue[1].Interface().(error); ok {
				if err != nil {
					// e.Logging(sessionLogger, err)
					httpx.HttpResponseErr(w, err)
					return
				}
			}
			httpx.HttpResponse(w, httpx.GetHTTPStatusCode(r.Context(), httpx.DefaultHttpSuccessCode), responseValue[0].Interface())
			return
		default:
			err := e.NewError(e.Error, e.ErrInternalServer, "wrong res type , first out should be response value , second out should be error ")
			// e.Logging(sessionLogger, err)
			httpx.HttpResponseErr(w, err)
			return
		}
	}
}

type mux interface {
	MethodFunc(method, pattern string, handlerFn http.HandlerFunc)
}

func DIRouter(r mux, container *Container) {
	InitInstance(container)
	RouterSelfCheck(container)
	// register router
	for _, controller := range _bootingControllers {
		for _, requestMapping := range controller.requestMaps {
			urlPath := filepath.Join(controller.rootPath, requestMapping.subPath)
			r.MethodFunc(requestMapping.method, urlPath, DIHandler(container, controller.instanceName, requestMapping.funcName, requestMapping.isRaw))
			container.c.Log.Infof("request mapping: method: %-6s %-30s => handler: %v.%v ", requestMapping.method, urlPath, controller.instanceName, requestMapping.funcName)
		}
	}

}
