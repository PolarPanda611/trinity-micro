package controller

import (
	"net/http"

	"github.com/PolarPanda611/trinity-micro/core/dbx"
	"github.com/PolarPanda611/trinity-micro/core/logx"
	"github.com/PolarPanda611/trinity-micro/core/tracerx"
	"github.com/go-chi/chi/v5"
)

var (
	apiHandler = []func(http.Handler) http.Handler{
		logx.ChiLoggerRequest,
		tracerx.ChiOpenTracer(
			tracerx.OperationNameFunc(
				func(r *http.Request) string {
					chiCtx := chi.RouteContext(r.Context())
					return r.Method + "=>" + chiCtx.RoutePattern()
				},
			),
		),
		dbx.SessionHandler,
	}
)
