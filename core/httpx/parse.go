/*
 * @Author: Daniel TAN
 * @Description:
 * @Date: 2021-08-06 09:15:58
 * @LastEditTime: 2021-09-02 17:58:29
 * @LastEditors: Daniel TAN
 * @FilePath: /fr-price-common-pkg/core/httpx/parse.go
 */
package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"trinity-micro/core/e"
	"trinity-micro/core/utils"

	"github.com/go-chi/chi/v5"
)

var (
	contextType    = reflect.ValueOf(context.Background()).Type()
	httpWriterType = reflect.ValueOf(NewWriter()).Type()
	requestType    = reflect.ValueOf(&http.Request{}).Type()
)

type w struct {
}

func (w w) Header() http.Header {
	return nil
}
func (w w) Write([]byte) (int, error) {
	return 0, nil
}
func (w w) WriteHeader(statusCode int) {
}
func NewWriter() http.ResponseWriter {
	return w{}
}

type HTTPContextKey string

const (
	HTTPStatus  HTTPContextKey = "HTTP_STATUS_CONTEXT_KEY"
	HTTPRequest HTTPContextKey = "HTTP_REQUEST_CONTEXT_KEY"
)

func GetHTTPStatusCode(ctx context.Context, defaultStatus int) int {
	if val, ok := ctx.Value(HTTPStatus).(*int); ok {
		if val != nil && *val != 0 {
			return *val
		}
	}
	return defaultStatus
}

func SetHttpStatusCode(ctx context.Context, status int) {
	if val, ok := ctx.Value(HTTPStatus).(*int); ok {
		*val = status
	}
}

func GetRawRequest(ctx context.Context) *http.Request {
	return ctx.Value(HTTPRequest).(*http.Request)
}

func Parse(r *http.Request, v interface{}) error {
	if v == nil {
		return fmt.Errorf("parsing error , empty value to parse")
	}
	destVal := reflect.Indirect(reflect.ValueOf(v))
	inType := destVal.Type()
	for index := 0; index < destVal.NumField(); index++ {
		val := destVal.Field(index)
		if !val.CanSet() {
			return e.NewError(e.Error, e.ErrDIParam, fmt.Sprintf("di param : %v is not exported , cannot set", inType.Field(index).Name))
		}
		// check if header param
		if headerParam, isExist := inType.Field(index).Tag.Lookup("header_param"); isExist {
			headerValString := r.Header.Get(headerParam)
			if err := utils.StringConverter(headerValString, &val); err != nil {
				return e.NewError(e.Error, e.ErrDIParam, fmt.Sprintf("header param %v converted error, cannot set ,err:%v ,  val : %v  ", inType.Field(index).Name, err, headerValString))
			}
		}
		// check if path param
		if pathParam, isExist := inType.Field(index).Tag.Lookup("path_param"); isExist {
			paramValString := chi.URLParam(r, pathParam)
			if err := utils.StringConverter(paramValString, &val); err != nil {
				return e.NewError(e.Error, e.ErrDIParam, fmt.Sprintf("path param %v converted error, cannot set , err:%v  val : %v  ", inType.Field(index).Name, err, paramValString))
			}
		}
		// check if query param
		if queryParam, isExist := inType.Field(index).Tag.Lookup("query_param"); isExist {
			if queryParam == "" {
				switch val.Type().Kind() {
				case reflect.String:
					val.Set(reflect.ValueOf(r.URL.RawQuery))
				case reflect.Map:
					switch inType.Field(index).Type.String() {
					case "url.Values":
						val.Set(reflect.ValueOf(r.URL.Query()))
					case "map[string][]string":
						res := make(map[string][]string)
						for k := range r.URL.Query() {
							res[k] = r.URL.Query()[k]
						}
						val.Set(reflect.ValueOf(res))
					case "map[string]string":
						res := make(map[string]string)
						for k := range r.URL.Query() {
							res[k] = r.URL.Query().Get(k)
						}
						val.Set(reflect.ValueOf(res))
					case "map[string]interface {}":
						res := make(map[string]interface{})
						for k := range r.URL.Query() {
							res[k] = r.URL.Query().Get(k)
						}
						val.Set(reflect.ValueOf(res))
					default:
						return e.NewError(e.Error, e.ErrDIParam, fmt.Sprintf("unsupported map type to decode query param , actual:%v", inType.Field(index).Type.String()))
					}
				default:
					return e.NewError(e.Error, e.ErrDIParam, fmt.Sprintf("param %v get all query param converted error, only support string , val : %v ", inType.Field(index).Name, r.URL.RawQuery))
				}
			} else {
				queryValString := r.URL.Query().Get(queryParam)
				if err := utils.StringConverter(queryValString, &val); err != nil {
					return e.NewError(e.Error, e.ErrDIParam, fmt.Sprintf("param %v converted error, err :%v , val : %v ", inType.Field(index).Name, err, queryValString))
				}
			}
		}
		// check if body param
		if bodyParam, isExist := inType.Field(index).Tag.Lookup("body_param"); isExist {
			respBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return e.NewError(e.Error, e.ErrReadRequestBody, fmt.Sprintf("read request body error  , err : %v ", err))
			}
			r.Body = ioutil.NopCloser(bytes.NewBuffer(respBytes))
			if bodyParam == "" {
				switch val.Type().Kind() {
				case reflect.String:
					val.Set(reflect.ValueOf(string(respBytes)))
				case reflect.Struct, reflect.Slice:
					// if is []byte
					if fmt.Sprintf("%v", inType.Field(index).Type) == "[]uint8" {
						val.Set(reflect.Indirect(reflect.ValueOf(respBytes)))
					} else {
						targetVal := reflect.New(inType.Field(index).Type).Interface()
						if err := json.Unmarshal(respBytes, targetVal); err != nil {
							return e.NewError(e.Error, e.ErrDecodeRequestBody, fmt.Sprintf("param %v converted error, err :%v , val : %v ", inType.Field(index).Name, err, string(respBytes)))
						}
						val.Set(reflect.Indirect(reflect.ValueOf(targetVal)))
					}
				case reflect.Map:
					if fmt.Sprintf("%v", inType.Field(index).Type) != "map[string]interface {}" {
						return e.NewError(e.Error, e.ErrDIParam, fmt.Sprintf("param %v converted error, map only support map[string]interface{}, val : %v ", inType.Field(index).Name, string(respBytes)))
					}
					bodyVal := make(map[string]interface{})
					if len(respBytes) > 0 {
						d := json.NewDecoder(bytes.NewReader(respBytes))
						d.UseNumber()
						if err := d.Decode(&bodyVal); err != nil {
							return e.NewError(e.Error, e.ErrDecodeRequestBody, fmt.Sprintf("param %v converted error,err :%v , val : %v ", inType.Field(index).Name, err, string(respBytes)))
						}
					}
					val.Set(reflect.ValueOf(bodyVal))
				case reflect.Interface:
					var bodyVal interface{}
					if len(respBytes) > 0 {
						if err := json.Unmarshal(respBytes, &bodyVal); err != nil {
							return e.NewError(e.Error, e.ErrDecodeRequestBody, fmt.Sprintf("param %v converted error,err :%v , val : %v ", inType.Field(index).Name, err, string(respBytes)))
						}
					}
					val.Set(reflect.ValueOf(bodyVal))
				default:
					return e.NewError(e.Error, e.ErrDIParam, fmt.Sprintf("unsupported type , only support string , struct ,Slice ,  map[string]interface{} , interface{} , []byte, actual: %v", val.Type().Kind()))
				}
			} else {
				bodyVal := make(map[string]interface{})
				if len(respBytes) > 0 {
					if err := json.Unmarshal(respBytes, &bodyVal); err != nil {
						return e.NewError(e.Error, e.ErrDecodeRequestBody, fmt.Sprintf("param %v converted error,err :%v , val : %v ", inType.Field(index).Name, err, string(respBytes)))
					}
				}
				value, err := bodyParamConverter(bodyVal, bodyParam, inType.Field(index).Type)
				if err != nil {
					return e.NewError(e.Error, e.ErrDecodeRequestBody, fmt.Sprintf("param %v converted error,err :%v , val : %v ", inType.Field(index).Name, err, bodyVal))
				}
				val.Set(reflect.ValueOf(value))
			}
		}
	}
	return nil
}

// bodyParamConverter
/*
 @bodyVal the father val of
 @key the key name of the parentVal
 @destType the type of dest type
 bodyParamConverter will get the key from the body value
 and convert the value to the dest type value
*/
func bodyParamConverter(bodyVal map[string]interface{}, key string, destType reflect.Type) (interface{}, error) {
	value, ok := bodyVal[key]
	if !ok {
		return nil, fmt.Errorf("key %v not exist", key)
	}
	switch destType.Kind() {
	case reflect.Int64:
		convertedValue, ok := value.(int64)
		if !ok {
			return nil, fmt.Errorf("key %v convert to int64 error ", key)
		}
		return convertedValue, nil
	case reflect.Int32:
		convertedValue, ok := value.(int32)
		if !ok {
			return nil, fmt.Errorf("key %v convert to int32 error ", key)
		}
		return convertedValue, nil
	case reflect.Int:
		convertedValue, ok := value.(int)
		if !ok {
			return nil, fmt.Errorf("key %v convert to int error ", key)
		}
		return convertedValue, nil
	case reflect.String:
		convertedValue, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("key %v convert to string error ", key)
		}
		return convertedValue, nil
	case reflect.Struct:
		c, _ := json.Marshal(value)
		targetVal := reflect.New(destType).Interface()
		decoder := json.NewDecoder(bytes.NewReader(c))
		if err := decoder.Decode(targetVal); err != nil {
			return nil, err
		}
		return reflect.Indirect(reflect.ValueOf(targetVal)).Interface(), nil
	default:
		return nil, fmt.Errorf("type %v not support", destType)
	}

}
