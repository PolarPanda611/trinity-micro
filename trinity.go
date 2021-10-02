// Author: Daniel TAN
// Date: 2021-09-06 00:13:18
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-02 23:22:58
// FilePath: /trinity-micro/trinity.go
// Description:
package trinity

import (
	"net/http"

	"github.com/PolarPanda611/trinity-micro/core/ioc/container"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

var (
	_defaultRouter = chi.NewRouter()
	_defaultLog    = logrus.New()
)

type Config struct {
	Mux    mux
	Logger logrus.FieldLogger
}

type Trinity struct {
	mux
	logrus.FieldLogger
	container *container.Container
}

func Default() *Trinity {
	return &Trinity{
		mux:         _defaultRouter,
		FieldLogger: _defaultLog,
		container: container.NewContainer(container.Config{
			AutoWired: true,
			Log:       _defaultLog,
		}),
	}
}

func New(c Config) *Trinity {
	return &Trinity{
		mux:         c.Mux,
		FieldLogger: c.Logger,
		container: container.NewContainer(container.Config{
			AutoWired: true,
			Log:       c.Logger,
		}),
	}
}

func (t *Trinity) Start(addr ...string) error {
	address := ":http"
	if len(addr) > 0 {
		address = addr[0]
	}
	t.Infof("service started at %v", address)
	return http.ListenAndServe(address, t.mux)
}
