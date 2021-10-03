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
	Status  int            `json:"status"`
	Result  interface{}    `json:"result,omitempty"`
	Error   *ResponseError `json:"error,omitempty"`
	TraceID string         `json:"trace_id"`
}

type ResponseError struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

func HttpResponseErr(ctx context.Context, w http.ResponseWriter, err error) {
	httpError := e.NewAPIError(err)
	res := &Response{
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
	JsonResponse(w, status, result)
}
