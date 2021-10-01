// Author: Daniel TAN
// Date: 2021-09-06 10:40:48
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-02 00:32:10
// FilePath: /trinity-micro/example/internal/adapter/controller/benchmark.go
// Description:
package controller

import (
	"sync"

	"github.com/PolarPanda611/trinity-micro"
)

func init() {
	BenchmarkControllerPool := &sync.Pool{
		New: func() interface{} {
			return new(benchmarkControllerImpl)
		},
	}
	trinity.RegisterInstance("BenchmarkController", BenchmarkControllerPool)
	trinity.RegisterController("/benchmark", "BenchmarkController",
		trinity.NewRequestMapping("GET", "/simple", "Simple"),
		trinity.NewRequestMapping("GET", "path_param_raw/{id}", "PathParam"),
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
