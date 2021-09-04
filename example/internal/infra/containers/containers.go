package containers

import "trinity-micro/core/ioc/container"

var Container *container.Container

func init() {
	Container = container.NewContainer()
}
