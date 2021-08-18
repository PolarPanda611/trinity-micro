package containers

import "trinity-micro/core/container"

var Container *container.Container

func Init() {
	c := container.Config{
		AutoFree:  true,
		AutoWired: true,
	}
	Container = container.NewContainer(c)
}
