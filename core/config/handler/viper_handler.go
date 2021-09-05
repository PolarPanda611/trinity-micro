package handler

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type ViperHandler struct {
	configPath string
}

func NewViperHandler(configPath string) *ViperHandler {
	return &ViperHandler{
		configPath: configPath,
	}
}
func (v *ViperHandler) init() {
	viper.AddConfigPath(v.configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

// Get value from key
func (v *ViperHandler) Get(key string) string {
	v.init()
	return fmt.Sprintf("%v", viper.Get(key))
}
