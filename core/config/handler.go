/*
 * @Author: Daniel TAN
 * @Description:
 * @Date: 2021-06-29 16:33:30
 * @LastEditTime: 2021-08-18 10:34:25
 * @LastEditors: Daniel TAN
 * @FilePath: /fr-price-common-pkg/config/handler.go
 */
package config

import "context"

// Handler the configer loader handler
type Handler interface {
	Get(ctx context.Context, key string) (string, error)
}
