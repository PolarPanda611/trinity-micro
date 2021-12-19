// Author: Daniel TAN
// Date: 2021-10-18 00:50:02
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-18 01:49:37
// FilePath: /trinity-micro/cron.go
// Description:
package trinity

import (
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type cronlog struct {
	log logrus.FieldLogger
}

func NewCronLogger(log logrus.FieldLogger) cron.Logger {
	return &cronlog{
		log: log,
	}
}

func (l *cronlog) Info(msg string, keysAndValues ...interface{}) {
	l.log.Info(msg, keysAndValues)
}

// Error logs an error condition.
func (l *cronlog) Error(err error, msg string, keysAndValues ...interface{}) {
	l.log.Error(err, msg, keysAndValues)
}

func (t *Trinity) AddCronJobs(resourceName string, spec string, cmd func(srv interface{}), jobNames ...string) {
	t.Lock()
	defer t.Unlock()
	t.cron.AddFunc(spec, func() {
		ins, injectMap := t.GetInstance(resourceName)
		defer t.PutInstance(resourceName, injectMap, ins)
		cmd(ins)
	})
	jobName := resourceName
	if len(jobNames) > 0 {
		jobName = jobNames[0]
	}
	t.log.Infof("cronjob register handler: job name: %10v, spec: %10v ", jobName, spec)
	t.cron.Stop()
	t.cron.Start()
}
