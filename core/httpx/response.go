package httpx

import (
	"encoding/json"
	"net/http"
	"trinity-micro/core/e"
)

const (
	DefaultHttpErrorCode   int = 400
	DefaultHttpSuccessCode int = 200
)

type Response struct {
	Status int            `json:"status"`
	Result interface{}    `json:"result"`
	Error  *ResponseError `json:"error,omitempty"`
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
	j, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpError.Status)
	w.Write(j)
}

func HttpResponse(w http.ResponseWriter, status int, res interface{}) {
	result := &Response{
		Status: status,
		Result: res,
	}
	j, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(j)
}
