// Author: Daniel TAN
// Date: 2021-09-05 10:24:33
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-03 14:56:13
// FilePath: /trinity-micro/core/httpx/response.go
// Description:
package httpx

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/PolarPanda611/trinity-micro/core/e"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

const (
	DefaultHttpErrorCode   int = 400
	DefaultHttpSuccessCode int = 200
)

func JsonResponse(w http.ResponseWriter, status int, res interface{}) {
	j, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(j)
}

func XMLResponse(w http.ResponseWriter, status int, res interface{}) {
	j, err := xml.Marshal(res)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	w.Write(j)
}



type Response struct {
	Status  int         `json:"status" example:"200"`
	Result  interface{} `json:"result,omitempty" `
	TraceID string      `json:"trace_id" example:"1-trace-it"`
}

type ErrorResponse struct {
	Status  int            `json:"status" example:"400"`
	Error   *ResponseError `json:"error,omitempty"`
	TraceID string         `json:"trace_id" example:"1-trace-it"`
}

type ResponseError struct {
	Code    int      `json:"code" example:"400001"`
	Message string   `json:"message" example:"ErrInvalidRequest"`
	Details []string `json:"details" example:"error detail1,error detail2"`
}

func HttpResponseErr(ctx context.Context, w http.ResponseWriter, err error) {
	httpError := e.NewAPIError(err)
	res := &ErrorResponse{
		Status: httpError.Status,
		Error: &ResponseError{
			Code:    httpError.Code,
			Message: httpError.Type,
			Details: []string{httpError.Message},
		},
	}
	JsonResponse(w, httpError.Status, res)
}

func HttpResponse(ctx context.Context, w http.ResponseWriter, status int, res interface{}) {
	result := &Response{
		Status: status,
		Result: res,
	}
	x := opentracing.SpanFromContext(ctx)
	sc, ok := x.Context().(jaeger.SpanContext)
	if ok {
		result.TraceID = sc.TraceID().String()
	}
	JsonResponse(w, status, result)
}
