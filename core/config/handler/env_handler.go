/*
 * @Author: Daniel TAN
 * @Description:
 * @Date: 2021-06-29 16:33:30
 * @LastEditTime: 2021-08-18 10:40:11
 * @LastEditors: Daniel TAN
 * @FilePath: /fr-price-common-pkg/config/handler/env_handler.go
 */
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
