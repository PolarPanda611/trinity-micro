// Author: Daniel TAN
// Date: 2021-09-06 00:13:18
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-04 01:37:32
// FilePath: /trinity-micro/trinity.go
// Description:
package trinity

import (
	"net/http"

	"github.com/PolarPanda611/trinity-micro/core/ioc/container"
	"github.com/PolarPanda611/trinity-micro/middleware"
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
	return New(Config{})
}

func New(c ...Config) *Trinity {
	if len(c) > 0 {
		if c[0].Mux == nil {
			c[0].Mux = _defaultRouter
		}
		if c[0].Logger == nil {
			c[0].Logger = _defaultLog
		}
	} else {
		c = append(c, Config{
			Mux:    _defaultRouter,
			Logger: _defaultLog,
		})
	}

	c[0].Mux.Use(middleware.Recovery(_defaultLog))
	c[0].Mux.Use(middleware.InitLogger(c[0].Logger))
	return &Trinity{
		mux:         c[0].Mux,
		FieldLogger: c[0].Logger,
		container: container.NewContainer(container.Config{
			AutoWired: true,
			Log:       c[0].Logger,
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
