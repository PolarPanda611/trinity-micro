// Author: Daniel TAN
// Date: 2021-08-18 23:34:44
// LastEditors: Daniel TAN
// LastEditTime: 2021-09-27 23:11:53
// FilePath: /trinity-micro/example/internal/infra/containers/containers.go
// Description:
package containers

import (
	"github.com/PolarPanda611/trinity-micro/example/internal/infra/logx"

	"github.com/PolarPanda611/trinity-micro/core/ioc/container"
)

var Container *container.Container

func init() {
	Container = container.NewContainer(container.Config{
		AutoWired: true,
		Log:       logx.Logger,
	})
}
