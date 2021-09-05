package config

import "context"

// Handler the configer loader handler
type Handler interface {
	Get(ctx context.Context, key string) (string, error)
}
