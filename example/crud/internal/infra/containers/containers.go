// Author: Daniel TAN
// Date: 2021-08-18 23:34:44
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-02 01:28:12
// FilePath: /trinity-micro/example/crud/internal/infra/containers/containers.go
// Description:
package containers

import (
	"github.com/PolarPanda611/trinity-micro/example/crud/internal/infra/logx"

	"github.com/PolarPanda611/trinity-micro/core/ioc/container"
)

type Config struct {
}

var Container *container.Container

func Init(c ...Config) {
	if len(c) > 0 {
		// init container with your config
		Container = container.NewContainer(container.Config{
			AutoWired: true,
			Log:       logx.Logger,
		})
		return
	}
	Container = container.NewContainer(container.Config{
		AutoWired: true,
		Log:       logx.Logger,
	})
}
