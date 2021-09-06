package controller

import (
	"sync"
	"trinity-micro/core/ioc/container"
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
