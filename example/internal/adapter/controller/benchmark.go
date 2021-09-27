// Author: Daniel TAN
// Date: 2021-09-06 10:40:48
// LastEditors: Daniel TAN
// LastEditTime: 2021-09-27 23:10:59
// FilePath: /trinity-micro/example/internal/adapter/controller/benchmark.go
// Description:
package controller

import (
	"sync"

	"github.com/PolarPanda611/trinity-micro/core/ioc/container"
)

func init() {
	BenchmarkControllerPool := &sync.Pool{
		New: func() interface{} {
			return new(benchmarkControllerImpl)
		},
	}
	container.RegisterInstance("BenchmarkController", BenchmarkControllerPool)
	container.RegisterController("/benchmark", "BenchmarkController",
		container.NewRequestMapping("GET", "/simple", "Simple"),
		container.NewRequestMapping("GET", "path_param_raw/{id}", "PathParam"),
	)
}

type benchmarkControllerImpl struct {
}

func (c *benchmarkControllerImpl) Simple() string {
	return "ok"
}

func (c *benchmarkControllerImpl) PathParam(Args struct {
	ID int `path_param:"id"`
}) int {
	return Args.ID
}
