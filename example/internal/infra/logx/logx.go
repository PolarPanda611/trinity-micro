package logx

import "github.com/sirupsen/logrus"

var Logger logrus.FieldLogger

func init() {
	Logger = logrus.New()
}
