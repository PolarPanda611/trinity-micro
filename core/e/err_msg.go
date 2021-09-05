/*
 * @Author: Daniel TAN
 * @Description:
 * @Date: 2021-04-02 19:05:41
 * @LastEditTime: 2021-09-02 17:51:39
 * @LastEditors: Daniel TAN
 * @FilePath: /fr-price-common-pkg/e/err_msg.go
 */
package e

import (
	"fmt"
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

// common
const (
	InternalError errorCode = 10000 + iota
	ValidationError
	PanicError
	StatusNotAcceptableError
)

//Batch
const (
	BatchMoveFileError errorCode = 20000 + iota
	BatchSetStatusError
	BatchGetLockError
	BatchLoadMetaError
	BatchAlreadyLockedError
	BatchLockError
	BatchUnLockError
	BatchStatusNotInQueue
	BatchExceededMaxAttempts
	BatchExceededMaxAge
	BatchUntrackedError
	BatchUnsupportedFileError
	BatchUnknownError
)

//DB
const (
	DBError errorCode = 30000 + iota
	RecordNotFoundError
	UpdateNoEffectError
	GetAdvisoryError
)

var (
	ErrIsWrongStatus        error = fmt.Errorf("wrong http response code")
	ErrIsReadResponseBody         = fmt.Errorf("read response body error ")
	ErrIsDecodeResponseBody       = fmt.Errorf("decode response body error ")
)

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
