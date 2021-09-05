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
