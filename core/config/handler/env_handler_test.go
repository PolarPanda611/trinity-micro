package handler

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	eHandler = &EnvHandler{}
	ctx      = context.Background()
)

func TestEnvHandler(t *testing.T) {
	os.Setenv("test1", "1")
	v1, _ := eHandler.Get(ctx, "test1")
	assert.Equal(t, "1", v1, "wrong env value")

	os.Setenv("test2", "2222")
	v2, _ := eHandler.Get(ctx, "test2")
	assert.Equal(t, "2222", v2, "wrong env value")

	os.Setenv("test3", "afafas")
	v3, _ := eHandler.Get(ctx, "test3")
	assert.Equal(t, "afafas", v3, "wrong env value")

	os.Setenv("test4", "阿发顺丰")
	v4, _ := eHandler.Get(ctx, "test4")
	assert.Equal(t, "阿发顺丰", v4, "wrong env value")
}
