// Author: Daniel TAN
// Date: 2021-10-04 01:30:26
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-04 01:35:04
// FilePath: /trinity-micro/middleware/recovery.go
// Description:
package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/PolarPanda611/trinity-micro/core/e"
	"github.com/PolarPanda611/trinity-micro/core/httpx"
	"github.com/sirupsen/logrus"
)

func Recovery(logger logrus.FieldLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// when stack finishes
					logMessage := fmt.Sprintf("Recovered from HTTP Request %v %v \n", r.Method, r.URL)
					logMessage += fmt.Sprintf("Trace: %s\n", err)
					logMessage += fmt.Sprintf("\n%s", debug.Stack())
					logger.Warn(logMessage)
					httpx.HttpResponseErr(r.Context(), w, e.NewError(e.Error, e.ErrPanic, fmt.Sprintf("panic %v", err)))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
