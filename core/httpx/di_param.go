package httpx

import (
	"context"
	"net/http"
	"reflect"

	"github.com/PolarPanda611/trinity-micro/core/e"
)

func DIParamHandler(handler interface{}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), HttpxContext, NewContext(r, 200)))
		handlerType := reflect.TypeOf(handler)
		// sessionLogger := log.ForContext(r.Context())
		inParams, err := InvokeHandler(handlerType, r)
		if err != nil {
			// e.Logging(sessionLogger, err)
			HttpResponseErr(r.Context(), w, err)
			return
		}
		responseValue := reflect.ValueOf(handler).Call(inParams)
		switch len(responseValue) {
		case 0:
			HttpResponse(r.Context(), w, GetHTTPStatusCode(r.Context(), DefaultHttpSuccessCode), nil)
			return
		case 1:
			if err, ok := responseValue[0].Interface().(error); ok {
				if err != nil {
					// e.Logging(sessionLogger, err)
					HttpResponseErr(r.Context(), w, err)
					return
				}
			}
			HttpResponse(r.Context(), w, GetHTTPStatusCode(r.Context(), DefaultHttpSuccessCode), responseValue[0].Interface())
			return
		case 2:
			if err, ok := responseValue[1].Interface().(error); ok {
				if err != nil {
					HttpResponseErr(r.Context(), w, err)
					return
				}
			}
			HttpResponse(r.Context(), w, GetHTTPStatusCode(r.Context(), DefaultHttpSuccessCode), responseValue[0].Interface())
			return
		default:
			err := e.NewError(e.Error, e.ErrInternalServer, "wrong res type , first out should be response value , second out should be error ")
			HttpResponseErr(r.Context(), w, err)
			return
		}

	}
}

func InvokeHandler(handlerType reflect.Type, r *http.Request) ([]reflect.Value, error) {
	if !IsHandler(handlerType) {
		return nil, e.NewError(e.Error, e.ErrDIParam, "wrong handler type , must be func ")
	}
	numsIn := HandlerNumsIn(handlerType)
	InParams := make([]reflect.Value, numsIn)
	var i = 0
	for i < numsIn {
		inType := handlerType.In(i)
		inKind := inType.Kind()
		switch inKind {
		case reflect.Interface:
			if !contextType.Implements(inType) {
				return nil, e.NewError(e.Error, e.ErrDIParam, "wrong handler , interface only support context")
			}
			InParams[i] = reflect.ValueOf(r.Context())
		case reflect.Struct:
			targetValue := reflect.New(inType).Interface()
			if err := Parse(r, targetValue); err != nil {
				return nil, err
			}
			InParams[i] = reflect.ValueOf(targetValue).Elem()
		default:
			return nil, e.NewError(e.Error, e.ErrDIParam, "wrong handler , unsupported type ")
		}
		i++
	}
	return InParams, nil
}

func InvokeMethod(handlerType reflect.Type, r *http.Request, instance interface{}, w http.ResponseWriter) ([]reflect.Value, error) {
	if !IsHandler(handlerType) {
		return nil, e.NewError(e.Error, e.ErrDIParam, "wrong handler type , must be func ")
	}
	numsIn := HandlerNumsIn(handlerType)
	InParams := make([]reflect.Value, numsIn)
	var i = 0
	for i < numsIn {
		inType := handlerType.In(i)
		inKind := inType.Kind()
		switch inKind {
		case reflect.Interface:
			if contextType.Implements(inType) {
				InParams[i] = reflect.ValueOf(r.Context())
				break
			}
			if httpWriterType.Implements(inType) {
				InParams[i] = reflect.ValueOf(w)
				break
			}
			return nil, e.NewError(e.Error, e.ErrDIParam, "wrong handler , interface only support context and httpResponseWriter")
		case reflect.Struct:
			targetValue := reflect.New(inType).Interface()
			if err := Parse(r, targetValue); err != nil {
				return nil, err
			}
			InParams[i] = reflect.ValueOf(targetValue).Elem()
		case reflect.Ptr:
			if inType == requestType {
				InParams[i] = reflect.ValueOf(r)
				break
			}
			if inType == reflect.TypeOf(instance) {
				InParams[i] = reflect.ValueOf(instance)
				break
			}
			targetValue := reflect.New(inType.Elem()).Interface()
			if err := Parse(r, targetValue); err != nil {
				return nil, err
			}
			InParams[i] = reflect.ValueOf(targetValue)
			// return nil, e.NewError(e.Error, e.ErrDIParam, "wrong handler , unsupported ptr ")
		default:
			return nil, e.NewError(e.Error, e.ErrDIParam, "wrong handler , unsupported type ")
		}
		i++
	}
	return InParams, nil
}

func IsHandler(handlerType reflect.Type) bool {
	return handlerType.Kind() == reflect.Func
}

func HandlerNumsIn(handlerType reflect.Type) int {
	return handlerType.NumIn()
}
