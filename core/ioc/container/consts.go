// Author: Daniel TAN
// Date: 2021-08-18 23:32:00
// LastEditors: Daniel TAN
// LastEditTime: 2021-12-17 23:50:47
// FilePath: /trinity-micro/core/ioc/container/consts.go
// Description:
package container

type Keyword string

const (
	_CONTAINER     Keyword = "container"
	_AUTOWIRE      Keyword = "autowire"
	_RESOURCE      Keyword = "resource"
	TAG_SPLITER            = ";"
	TAG_KV_SPLITER         = ":"
	CONTEXT                = "CONTEXT"
)
