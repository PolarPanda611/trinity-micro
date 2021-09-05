/*
 * @Author: Daniel TAN
 * @Description:
 * @Date: 2021-04-02 17:31:18
 * @LastEditTime: 2021-08-30 19:13:14
 * @LastEditors: Daniel TAN
 * @FilePath: /fr-price-common-pkg/e/error.go
 */
/*
 * @Author: Daniel TAN
 * @Description:
 * @Date: 2020-11-23 13:10:45
 * @LastEditTime: 2020-11-24 11:12:38
 * @LastEditors: Daniel TAN
 * @FilePath: /nma/internal/job/e/e.go
 */

package e

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

// LogLevel log level
type LogLevel string

const (
	// Info info level
	Info LogLevel = "Info"
	// Error err level
	Error = "Error"
	// Warn warn level
	Warn = "Warn"
)

// http error
type httpError struct {
	Status  int
	Code    int
	Type    string
	Message string
}

// WrapError wrap err
type WrapError interface {
	LogLevel() LogLevel
	Code() errorCode
	Error() string
	Message() string
}

type wrapErrorImpl struct {
	code      errorCode
	msg       string
	rootCause error
	loglevel  LogLevel
}

// NewError  new error
/**
 * @Author: Daniel TAN
 * @Description:
 * @Date: 2020-11-23 19:25:16
 * @LastEditors: Daniel TAN
 * @param {LogLevel} loglevel
 * @param {string} msg
 * @param {error} err
 * @return {*}
 */
func NewError(loglevel LogLevel, code errorCode, msg string, rootCauses ...error) WrapError {
	errCodeMsg := code.GetErrorMsg()
	if msg == "" {
		msg = errCodeMsg
	}
	newErr := &wrapErrorImpl{
		loglevel: loglevel,
		code:     code,
		msg:      msg,
	}
	if len(rootCauses) > 0 {
		newErr.rootCause = rootCauses[0]
	} else {
		newErr.rootCause = errors.New(msg)
	}
	return newErr
}
func (e *wrapErrorImpl) Error() string {
	errCodeMsg := e.code.GetErrorMsg()
	return fmt.Sprintf("error code: %v, error type: %v, error message: %v, actual error: %v", e.code, errCodeMsg, e.msg, e.rootCause.Error())
}
func (e *wrapErrorImpl) Code() errorCode {
	return e.code
}

func (e *wrapErrorImpl) Message() string {
	return e.msg
}

func (e *wrapErrorImpl) LogLevel() LogLevel {
	return e.loglevel
}

// Logging log func
/**
 * @Author: Daniel TAN
 * @Description:
 * @Date: 2020-11-24 09:51:48
 * @LastEditors: Daniel TAN
 * @param {*logrus.Entry} logger
 * @param {error} err
 * @return {*}
 */
func Logging(logger logrus.FieldLogger, err error) {
	if logger == nil {
		return
	}
	if err == nil {
		return
	}
	wrapError, ok := err.(WrapError)
	if ok {
		switch wrapError.LogLevel() {
		case Info:
			logger.Info(err)
		case Error:
			logger.Error(err)
		case Warn:
			logger.Warn(err)
		default:
			logger.Error(err)
		}
	} else {
		logger.Error(err)
	}

}

// create a http error from the eror
func NewAPIError(err error) *httpError {
	if err == nil {
		return nil
	}
	wrapError, ok := err.(WrapError)
	if !ok {
		return &httpError{
			Status:  http.StatusInternalServerError,
			Code:    ErrInternalServer.Int(),
			Type:    ErrInternalServer.GetErrorMsg(),
			Message: err.Error(),
		}
	}
	return &httpError{
		Status:  wrapError.Code().GetHttpStatusCode(),
		Code:    wrapError.Code().Int(),
		Type:    wrapError.Code().GetErrorMsg(),
		Message: wrapError.Message(),
	}

}

func Is(err error, targetErrorCode errorCode) bool {
	wrapError, ok := err.(WrapError)
	if !ok {
		return false
	}
	return wrapError.Code() == targetErrorCode
}
