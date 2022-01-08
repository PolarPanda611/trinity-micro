// Author: Daniel TAN
// Date: 2021-10-02 00:36:09
// LastEditors: Daniel TAN
// LastEditTime: 2021-12-18 00:11:08
// FilePath: /trinity-micro/di_router.go
// Description:
package trinity

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"sync"
	"syscall"

	"github.com/PolarPanda611/trinity-micro/core/e"
	"github.com/PolarPanda611/trinity-micro/core/httpx"
	"github.com/PolarPanda611/trinity-micro/core/ioc/container"
)

type bootingController struct {
	rootPath     string
	instanceName string
	requestMaps  []RequestMap
}

var (
	// booting buffer params
	_bootingControllers []bootingController
)

var (
	// booting cache
	injectMapPool = &sync.Pool{
		New: func() interface{} {
			return make(map[string]interface{})
		},
	}
)

type RequestMap struct {
	method   string
	subPath  string
	funcName string
	handlers []func(http.Handler) http.Handler
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

func NewRequestMapping(method string, path string, funcName string, handlers ...func(http.Handler) http.Handler) RequestMap {
	return RequestMap{
		method:   method,
		subPath:  path,
		funcName: funcName,
		handlers: handlers,
		isRaw:    false,
	}
}

func NewRawRequestMapping(method string, path string, funcName string, handlers ...func(http.Handler) http.Handler) RequestMap {
	return RequestMap{
		method:   method,
		subPath:  path,
		funcName: funcName,
		handlers: handlers,
		isRaw:    true,
	}
}

func (t *Trinity) initInstance() {
	for _, instance := range _bootingInstances {
		t.container.RegisterInstance(instance.instanceName, instance.instancePool)
		t.log.Infof("%-8v %-10v %-7v => %v ", "instance", "register", "success", instance.instanceName)
	}
	if err := t.container.InstanceDISelfCheck(); err != nil {
		t.log.Fatalf("%-10v %-10v %-7v, err: %v", "instance", "self-check", "failed", err)
	}
}

// multi instance di handler
func DIHandler(container *container.Container, instanceName string, funcName string, isRaw bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), httpx.HttpxContext, httpx.NewContext(r, 0)))
		injectMap := injectMapPool.Get().(map[string]interface{})
		instance := container.GetInstance(instanceName, injectMap)
		defer func() {
			for k, v := range injectMap {
				container.Release(k, v)
				delete(injectMap, k)
			}
			injectMapPool.Put(injectMap)
		}()
		currentMethod, ok := reflect.TypeOf(instance).MethodByName(funcName)
		if !ok {
			panic("method not registered, please ensure your run thee RouterSelfCheck before start your service")
		}
		inParams, err := httpx.InvokeMethod(currentMethod.Type, r, instance, w)
		if err != nil {
			// e.Logging(sessionLogger, err)
			httpx.HttpResponseErr(r.Context(), w, err)
			return
		}
		responseValue := currentMethod.Func.Call(inParams)
		if isRaw {
			return
		}
		switch len(responseValue) {
		case 0:
			httpx.HttpResponse(r.Context(), w, httpx.GetHTTPStatusCode(r.Context(), httpx.DefaultHttpSuccessCode), nil)
			return
		case 1:
			if err, ok := responseValue[0].Interface().(error); ok {
				if err != nil {
					// e.Logging(sessionLogger, err)
					httpx.HttpResponseErr(r.Context(), w, err)
					return
				}
			}
			httpx.HttpResponse(r.Context(), w, httpx.GetHTTPStatusCode(r.Context(), httpx.DefaultHttpSuccessCode), responseValue[0].Interface())
			return
		case 2:
			if err, ok := responseValue[1].Interface().(error); ok {
				if err != nil {
					// e.Logging(sessionLogger, err)
					httpx.HttpResponseErr(r.Context(), w, err)
					return
				}
			}
			httpx.HttpResponse(r.Context(), w, httpx.GetHTTPStatusCode(r.Context(), httpx.DefaultHttpSuccessCode), responseValue[0].Interface())
			return
		default:
			err := e.NewError(e.Error, e.ErrInternalServer, "wrong res type , first out should be response value , second out should be error ")
			// e.Logging(sessionLogger, err)
			httpx.HttpResponseErr(r.Context(), w, err)
			return
		}
	}
}

func (t *Trinity) diRouter() {
	t.routerSelfCheck()
	// register router
	for _, controller := range _bootingControllers {
		for _, requestMapping := range controller.requestMaps {
			urlPath := filepath.Join(controller.rootPath, requestMapping.subPath)
			h := http.HandlerFunc(DIHandler(t.container, controller.instanceName, requestMapping.funcName, requestMapping.isRaw))
			for i := len(requestMapping.handlers) - 1; i >= 0; i-- {
				h = requestMapping.handlers[i](h).ServeHTTP
			}
			t.mux.MethodFunc(requestMapping.method, urlPath, h)
			t.log.Infof("router   register handler: %-6s %-30s => %v.%v ", requestMapping.method, urlPath, controller.instanceName, requestMapping.funcName)
		}
	}
}

func (t *Trinity) routerSelfCheck() {
	for _, controller := range _bootingControllers {
		for _, requestMap := range controller.requestMaps {
			injectMap := injectMapPool.Get().(map[string]interface{})
			instance := t.container.GetInstance(controller.instanceName, injectMap)
			defer func() {
				for k, v := range injectMap {
					t.container.Release(k, v)
					delete(injectMap, k)
				}
				injectMapPool.Put(injectMap)
			}()
			_, ok := reflect.TypeOf(instance).MethodByName(requestMap.funcName)
			if !ok {
				t.log.Fatalf("%-8v %-10v %-7v => %v.%v , func %v not exist ", "router", "self-check", "failed", controller.instanceName, requestMap.funcName, requestMap.funcName)
				continue
			}
			t.log.Infof("%-8v %-10v %-7v => %v.%v ", "router", "self-check", "success", controller.instanceName, requestMap.funcName)
		}
	}

}

func (t *Trinity) Start(addr ...string) error {
	t.diRouter()
	address := ":http"
	if len(addr) > 0 {
		address = addr[0]
	}
	t.log.Infof("service started at %v", address)
	gErr := make(chan error)
	go func() {
		gErr <- http.ListenAndServe(address, t.mux)
	}()
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		gErr <- fmt.Errorf("receive %s", <-sigChan)
	}()

	return <-gErr
}
