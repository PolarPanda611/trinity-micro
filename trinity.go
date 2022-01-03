package trinity

import (
	"sync"

	"github.com/PolarPanda611/trinity-micro/core/ioc/container"
	"github.com/PolarPanda611/trinity-micro/core/logx"
	"github.com/PolarPanda611/trinity-micro/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/robfig/cron/v3"
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
	sync.RWMutex
	mux
	log       logrus.FieldLogger
	container *container.Container
	cron      *cron.Cron
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
	ins := &Trinity{
		mux: c[0].Mux,
		log: c[0].Logger,
		container: container.NewContainer(container.Config{
			AutoWire: true,
			Log:      c[0].Logger,
		}),
		cron: cron.New(cron.WithChain(
			cron.Recover(NewCronLogger(c[0].Logger)), // or use cron.DefaultLogger
		)),
	}

	ins.mux.Use(middleware.Recovery(ins.log))
	ins.mux.Use(logx.SessionLogger(ins.log))
	ins.cron.Start()
	ins.initInstance()
	return ins
}

func (t *Trinity) GetInstance(resourceName string) (interface{}, map[string]interface{}) {
	injectMap := injectMapPool.Get().(map[string]interface{})
	ins := t.container.GetInstance(resourceName, injectMap)
	return ins, injectMap
}

func (t *Trinity) PutInstance(insName string, injectMap map[string]interface{}, ins interface{}) {
	for k, v := range injectMap {
		t.container.Release(k, v)
		delete(injectMap, k)
	}
	injectMapPool.Put(injectMap)
}
