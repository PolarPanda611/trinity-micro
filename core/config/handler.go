// Author: Daniel TAN
// Date: 2021-09-03 12:24:12
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-02 00:24:08
// FilePath: /trinity-micro/core/config/handler.go
// Description:
package config

import "context"

// Handler the configer loader handler
type Handler interface {
	Get(ctx context.Context, key string) (string, error)
	Exists(ctx context.Context, key string) (bool, error)
}
