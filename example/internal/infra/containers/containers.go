package containers

import (
	"trinity-micro/core/ioc/container"
	"trinity-micro/example/internal/infra/logx"
)

var Container *container.Container

func init() {
	Container = container.NewContainer(container.Config{
		AutoWired: true,
		Log:       logx.Logger,
	})
}
