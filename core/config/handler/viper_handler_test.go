package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	vHandler *ViperHandler = NewViperHandler(".")
)

func TestViperHandler(t *testing.T) {
	assert.Equal(t, "true", vHandler.Get("Hacker"), "wrong env value")
	assert.Equal(t, "[skateboarding snowboarding go]", vHandler.Get("hobbies"), "wrong env value")
	assert.Equal(t, "leather", vHandler.Get("clothing.jacket"), "wrong env value")
}
