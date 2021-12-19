// Author: Daniel TAN
// Date: 2021-09-05 16:48:38
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-03 23:32:30
// FilePath: /trinity-micro/example/crud/internal/infra/logx/logx.go
// Description:
package logger

import "github.com/sirupsen/logrus"

type Config struct {
	ServiceName string
}

var Logger logrus.FieldLogger

func Init(c ...Config) {
	if len(c) > 0 {
		// init the logger with configuration
		Logger = logrus.New().WithField("service", c[0].ServiceName)
		return
	}
	Logger = logrus.New()
}
