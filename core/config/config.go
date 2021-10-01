package config

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	skipFlag = "-"
)

// Option config option for config
type Option struct {
	Debug  bool
	Prefix string
}

type param struct {
	parent      interface{}
	parentValue reflect.Value
	index       int
	config      *paramConfig

	//runtime
	sourceValue string
}

func newParam(parent interface{}, parentValue reflect.Value, index int, paramConfig *paramConfig) *param {
	return &param{
		parent:      parent,
		parentValue: parentValue,
		index:       index,
		config:      paramConfig,
	}
}

func (p *param) getKeyName() string {
	if p.config.name == "" {
		return getParamName(p.parent, p.index)
	}
	return p.config.name
}

func (p *param) getFullName() string {
	return encodeObjectName(p.parent, p.index)
}

func (p *param) getDefaultValue() *string {
	return p.config.defaultValue
}

func (p *param) canSet() bool {
	return p.getValue().CanSet()
}
func (p *param) getValue() reflect.Value {
	return p.parentValue.Field(p.index)
}

func (p *param) hasConfig() bool {
	return p.config.exist
}

func (p *param) setValue(value string) error {
	p.sourceValue = value
	defaultValue, err := typeConverter(value, p.getValue().Type())
	if err != nil {
		return err
	}
	p.getValue().Set(reflect.ValueOf(defaultValue))
	return nil
}

func (p *param) isValid() error {
	v := p.getValue().Interface()
	if field, ok := v.(Field); ok {
		if !field.IsValid() {
			return fmt.Errorf("value %v invalid for type %v", v, p.getValue().Type())
		}
	}
	return nil
}

func (p *param) validate() error {
	if err := p.isValid(); err != nil {
		return err
	}
	if !p.config.optional && p.sourceValue == "" {
		return fmt.Errorf("the param %v is mandantory", p.getFullName())
	}
	return nil
}

type paramConfig struct {
	name         string  // default name
	defaultValue *string // default value
	optional     bool    // optional
	exist        bool
}

// Config instance
type Config struct {
	option   Option
	handlers []Handler
}

// Load Config instance load the setting to dest
func (c *Config) Load(ctx context.Context, dest interface{}) error {
	instanceVal := reflect.Indirect(reflect.ValueOf(dest))
	if !instanceVal.CanAddr() {
		return fmt.Errorf("Config %v should be addressable", reflect.TypeOf(dest))
	}
	for index := 0; index < instanceVal.NumField(); index++ {
		// if skip the config
		configTagValue := GetTagsValue(dest, index, "config")
		if configTagValue == skipFlag {
			continue
		}
		param := newParam(dest, instanceVal, index, decodeParamConfig(configTagValue))
		if !param.canSet() {
			continue
		}
		switch param.getValue().Type().Kind() {
		case reflect.Struct:
			v := reflect.New(param.getValue().Type()).Interface()
			if !param.hasConfig() {
				if err := c.Load(ctx, v); err != nil {
					return err
				}
				param.getValue().Set(reflect.Indirect(reflect.ValueOf(v)))
				if err := param.isValid(); err != nil {
					return err
				}
				continue
			}
		case reflect.Ptr:
			if reflect.TypeOf(param.getValue().Interface()).Elem().Kind() == reflect.Struct {
				v := reflect.New(reflect.TypeOf(param.getValue().Interface()).Elem()).Interface()
				if err := c.Load(ctx, v); err != nil {
					return err
				}
				param.getValue().Set(reflect.ValueOf(v))
				if err := param.isValid(); err != nil {
					return err
				}
				continue
			}
		}
		// Load default Value
		defaultValue := param.getDefaultValue()
		if defaultValue != nil {
			if err := param.setValue(*defaultValue); err != nil {
				return err
			}
		}

		// Load handler
		for _, v := range c.handlers {
			exists, err := v.Exists(ctx, param.getKeyName())
			if err != nil {
				return err
			}
			if !exists {
				continue
			}
			value, err := v.Get(ctx, param.getKeyName())
			if err != nil {
				return err
			}
			if err := param.setValue(value); err != nil {
				return err
			}

		}
		if err := param.validate(); err != nil {
			return err
		}

	}

	return c.isValid(dest)
}

func (c *Config) isValid(dest interface{}) error {
	if field, ok := dest.(Field); ok {
		if !field.IsValid() {
			return fmt.Errorf("value %v invalid for type %v", dest, reflect.TypeOf(dest))
		}
	}
	return nil
}

// New new config instance with customize setting
func New(option Option, handlers ...Handler) *Config {
	return &Config{
		option:   option,
		handlers: handlers,
	}
}

// Default default config instance with customize setting
func Default(handlers ...Handler) *Config {
	return &Config{
		option: Option{
			Debug:  true,
			Prefix: "",
		},
		handlers: handlers,
	}
}

func decodeParamConfig(config string) *paramConfig {
	var p paramConfig
	if config == "" {
		return &p
	}
	p.exist = true
	keyValues := strings.Split(strings.TrimSuffix(config, ";"), ";")
	for _, keyValue := range keyValues {
		kv := strings.Split(keyValue, ":")
		if len(kv) >= 2 {
			key := kv[0]
			value := strings.Join(kv[1:], ":")
			switch strings.TrimSpace(strings.ToLower(key)) {
			case "name":
				p.name = value
			case "default":
				p.defaultValue = &value
			case "optional":
				p.optional, _ = strconv.ParseBool(value)
			}
		}
	}
	return &p
}
