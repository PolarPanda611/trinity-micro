// Author: Daniel TAN
// Date: 2021-09-03 12:24:12
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-02 00:24:37
// FilePath: /trinity-micro/core/config/handler/env_handler.go
// Description:
package handler

import (
	"context"
	"os"
)

// EnvHandler env handler
type EnvHandler struct {
}

// Get value from key
func (e *EnvHandler) Get(ctx context.Context, key string) (string, error) {
	return os.Getenv(key), nil
}

func (e *EnvHandler) Exists(ctx context.Context, key string) (bool, error) {
	_, exists := os.LookupEnv(key)
	return exists, nil
}
