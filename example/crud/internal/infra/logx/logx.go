// Author: Daniel TAN
// Date: 2021-09-05 16:48:38
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-02 00:42:23
// FilePath: /trinity-micro/example/internal/infra/logx/logx.go
// Description:
package logx

import "github.com/sirupsen/logrus"

type Config struct {
}

var Logger logrus.FieldLogger

func Init(c ...Config) {
	if len(c) > 0 {
		// init the logger with configuration
		Logger = logrus.New()
		return
	}
	Logger = logrus.New()
}
