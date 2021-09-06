package httpx

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"trinity-micro/core/e"
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
	Code    int
	Message string
	Details []string
}

func HttpResponseErr(w http.ResponseWriter, err error) {
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

func HttpResponse(w http.ResponseWriter, status int, res interface{}) {
	result := &Response{
		Status: status,
		Result: res,
	}
	JsonResponse(w, status, result)
}
