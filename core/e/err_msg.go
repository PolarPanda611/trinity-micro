package e

import (
	"net/http"
	"strconv"
)

type errorCode int

func (c errorCode) GetErrorMsg() string {
	msg, ok := errTypeMap[c]
	if !ok {
		return "NoPredefinedErrorMessage"
	}
	return msg
}

func (c errorCode) GetHttpStatusCode() int {
	str := strconv.Itoa(int(c))
	code := string([]byte(str)[:3])
	r, _ := strconv.Atoi(code)
	if r < 100 || r >= 600 {
		return http.StatusInternalServerError
	}
	return r
}
func (c errorCode) Int() int {
	return int(c)
}

var (
	errTypeMap = map[errorCode]string{
		ErrInvalidRequest:                  "InvalidRequest",
		ErrHttpResponseCodeNotSuccess:      "HttpResponseCodeNotSuccess",
		ErrRecordNotFound:                  "RecordNotFound",
		ErrResourceNotFound:                "ResourcedNotFound",
		ErrInternalServer:                  "InternalServerError",
		ErrReadResponseBody:                "ReadResponseBodyError",
		ErrDecodeResponseBody:              "DecodeResponseBodyError",
		ErrVertexAccessDeniedException:     "VertexAccessDeniedException",
		ErrVertexNumberFormatException:     "VertexNumberFormatException",
		ErrVertexInvalidAddressException:   "VertexInvalidAddressException",
		ErrVertexInvalidTaxAreaIdException: "VertexInvalidTaxAreaIdException",
		ErrVertexApplicationException:      "VertexApplicationException",
		ErrVertexInvalidCountryException:   "VertexInvalidCountryException",
		ErrVertexUnknownException:          "VertexUnknownException",
		ErrExecuteSQL:                      "ExecuteSQLError",
		ErrAdvisoryLock:                    "AdvisoryLockError",
		ErrAdvisoryUnlock:                  "AdvisoryUnlockError",
		ErrDIParam:                         "DIParamError",
		ErrReadRequestBody:                 "ReadRequestBodyError",
		ErrDecodeRequestBody:               "DecodeRequestBodyError",
	}
)
