package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/PolarPanda611/trinity-micro/core/config/handler"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
)

type Env string

const (
	Prod  Env = "prod"
	Stage Env = "stage"
	Dev   Env = "dev"
)

func (e Env) IsValid() bool {
	switch string(e) {
	case "prod", "stage", "dev":
		return true
	}
	return false
}

var (
	envHandler handler.EnvHandler
	ctx                = context.Background()
	configor   *Config = Default(&envHandler)
)

type mockHandler struct {
	mock.Mock
}

func (m *mockHandler) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}
func (m *mockHandler) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

type structConfig struct {
	Debug bool   `config:"default:true"`
	Host  string `config:"default:xxxxx"`
}

func (s structConfig) IsValid() bool {
	return true
}

type structConfig2 struct {
	Debug2 bool   `config:"default:true"`
	Host2  string `config:"default:xxxxx"`
}
type testSuccessCase struct {
	TestString   string       `config:"default:test;"`
	TestInt      int          `config:"default:1;"`
	TestInt16    int16        `config:"default:1;"`
	TestInt32    int32        `config:"default:1;"`
	TestInt64    int64        `config:"default:1;"`
	TestUint     uint         `config:"default:1;"`
	TestUint16   uint16       `config:"default:1;"`
	TestUint32   uint32       `config:"default:1;"`
	TestUint64   uint64       `config:"default:1;"`
	TestFloat32  float32      `config:"default:1;"`
	TestFloat64  float64      `config:"default:1;"`
	GoPath       string       `config:"default:1;"`
	GoPathAlias  string       `config:"name:GOPATHTEST;default:1;"`
	GoProxy114   string       `config:"name:GOPROXY114TEST;"`
	Mandantory   string       `config:"name:mmmmmm;optional:true"`
	Env          Env          `config:"default:prod;"`
	DebugFalse   bool         `config:"default:false;"`
	DebugTrue    bool         `config:"default:true;"`
	Sfull        structConfig `config:"name:sfull"`
	SfullDefault structConfig `config:"default:{\"Debug\":false,\"Host\":\"test struct\"}"`
	S            structConfig
	SPtr         *structConfig2
	TestSkip     int         `config:"-"`
	TestSkip2    string      `config:"-"`
	TestSkip3    interface{} `config:"-"`
}

type negativeTestCase struct {
	BrokenJson    structConfig `config:"name:broken"`
	BrokenDefault structConfig `config:"name:brokendefault;default:"`
}

func TestLoadDefault(t *testing.T) {
	os.Setenv("GOPATHTEST", "GOPATHTESTVALUE")
	os.Setenv("GOPROXY114TEST", "GOPROXY114TESTVALUE")
	os.Setenv("Host2", "Host2xxxx")
	os.Setenv("Debug2", "false")
	os.Setenv("sfull", "{\"Debug\":false,\"Host\":\"test struct\"}")
	os.Setenv("broken", "{\"")
	os.Setenv("TestSkip1", "1235")
	os.Setenv("TestSkip2", "1235")
	os.Setenv("TestSkip3", "1235")

	var c testSuccessCase
	err := Default(&envHandler).Load(ctx, &c)
	assert.Equal(t, nil, err, "err during load config")
	assert.Equal(t, "test", c.TestString, "wrong TestString")
	assert.Equal(t, int(1), c.TestInt, "wrong TestInt")
	assert.Equal(t, int16(1), c.TestInt16, "wrong TestInt16")
	assert.Equal(t, int32(1), c.TestInt32, "wrong TestInt32")
	assert.Equal(t, int64(1), c.TestInt64, "wrong TestInt64")
	assert.Equal(t, uint(1), c.TestUint, "wrong TestUint")
	assert.Equal(t, uint16(1), c.TestUint16, "wrong TestUInt16")
	assert.Equal(t, uint32(1), c.TestUint32, "wrong TestUInt32")
	assert.Equal(t, uint64(1), c.TestUint64, "wrong TestUInt64")
	assert.Equal(t, float32(1), c.TestFloat32, "wrong TestFloat32")
	assert.Equal(t, float64(1), c.TestFloat64, "wrong TestFloat64")
	assert.Equal(t, "1", c.GoPath, "wrong GoPath")
	assert.Equal(t, "GOPATHTESTVALUE", c.GoPathAlias, "wrong GoPathAlias")
	assert.Equal(t, "GOPROXY114TESTVALUE", c.GoProxy114, "wrong GoProxy114")
	assert.Equal(t, "", c.Mandantory, "wrong Mandantory")
	assert.Equal(t, Prod, c.Env, "wrong enum")
	assert.Equal(t, false, c.DebugFalse, "wrong false bool")
	assert.Equal(t, true, c.DebugTrue, "wrong true bool")
	assert.Equal(t, true, c.S.Debug, "wrong loop true bool")
	assert.Equal(t, "xxxxx", c.S.Host, "wrong loop host bool")
	assert.Equal(t, false, c.SPtr.Debug2, "wrong loop true bool")
	assert.Equal(t, "Host2xxxx", c.SPtr.Host2, "wrong loop host bool")
	assert.Equal(t, structConfig{Debug: false, Host: "test struct"}, c.Sfull, "wrong Sfull")
	assert.Equal(t, structConfig{Debug: false, Host: "test struct"}, c.SfullDefault, "wrong Sfull")
	assert.Equal(t, 0, c.TestSkip, "wrong TestSkip")
	assert.Equal(t, "", c.TestSkip2, "wrong TestSkip2")
	assert.Equal(t, nil, c.TestSkip3, "wrong TestSkip3")

	var negativeTestConfig negativeTestCase
	configLoader := Default(&envHandler)
	assert.NotEqual(t, configLoader.Load(ctx, &negativeTestConfig), nil)
	os.Setenv("broken", "{}")
	assert.NotEqual(t, configLoader.Load(ctx, &negativeTestConfig), nil)

	mockH := &mockHandler{}
	expectErr := fmt.Errorf("test handlers")
	type testConfig struct {
		GoPath string `config:"name:GOPATHTEST;default:1"`
	}
	c2 := testConfig{}
	mockH.On("Exists", ctx, "GOPATHTEST").Once().Return(false, nil)
	Default(mockH).Load(ctx, &c2)
	assert.Equal(t, c2.GoPath, "1")
	mockH.On("Exists", ctx, "GOPATHTEST").Once().Return(false, expectErr)
	assert.Equal(t, Default(mockH).Load(ctx, &c2), expectErr)
	mockH.On("Exists", ctx, "GOPATHTEST").Once().Return(true, nil)
	mockH.On("Get", ctx, "GOPATHTEST").Once().Return("nonono", nil)
	Default(mockH).Load(ctx, &c2)
	assert.Equal(t, c2.GoPath, "nonono")
	mockH.On("Exists", ctx, "GOPATHTEST").Once().Return(true, nil)
	mockH.On("Get", ctx, "GOPATHTEST").Once().Return("", expectErr)
	assert.Equal(t, Default(mockH).Load(ctx, &c2), expectErr)
}

type TestFailedPointerCase struct {
	TestInt *int `config:"default:1;"`
}

func TestFailedPointer(t *testing.T) {
	var c TestFailedPointerCase
	err := configor.Load(ctx, &c)
	assert.Equal(t, errors.New("unsupported pointer type"), err, "wrong during loading ptr and non struct field")
}

type TestMandantoryCase struct {
	NotExistEnv string `config:"name:notexist"`
}

func TestTestMandantory(t *testing.T) {
	var c TestMandantoryCase
	err := configor.Load(ctx, &c)
	assert.Equal(t, "the param *config.TestMandantoryCase.NotExistEnv.(string) is mandantory", err.Error(), "wrong during test mandantory")
}

func TestDecodeParamConfig(t *testing.T) {
	config1 := "name:GOENV;default:test;optional:false;"
	defaultValue := "test"
	assert.Equal(t, &paramConfig{
		name:         "GOENV",
		defaultValue: &defaultValue,
		optional:     false,
		exist:        true,
	}, decodeParamConfig(config1), "wrong decode config1")

	config2 := "name:GOENV;default:test;optional:true;"
	assert.Equal(t, &paramConfig{
		name:         "GOENV",
		defaultValue: &defaultValue,
		optional:     true,
		exist:        true,
	}, decodeParamConfig(config2), "wrong decode config2")

	config3 := "   name:GOENV;  default:test;optional:true;"
	assert.Equal(t, &paramConfig{
		name:         "GOENV",
		defaultValue: &defaultValue,
		optional:     true,
		exist:        true,
	}, decodeParamConfig(config3), "wrong decode config3")

	config4 := "   name:GoEnv;  default:1;optional:true"
	defaultVal := "1"
	assert.Equal(t, &paramConfig{
		name:         "GoEnv",
		defaultValue: &defaultVal,
		optional:     true,
		exist:        true,
	}, decodeParamConfig(config4), "wrong decode config4")
}

type testWrongEnum struct {
	Env Env `config:"default:xxx;"`
}

func TestWrongEnum(t *testing.T) {
	var c testWrongEnum
	err := configor.Load(ctx, &c)
	assert.Equal(t, errors.New("value xxx invalid for type config.Env"), err, "wrong enum value")
}

type testConfigIsValid struct {
	Invalid bool `config:"default:false"`
}

func (t *testConfigIsValid) IsValid() bool {
	return false
}

func TestConfigIsValid(t *testing.T) {
	var c testConfigIsValid
	err := configor.Load(ctx, &c)
	assert.Equal(t, errors.New("value &{false} invalid for type *config.testConfigIsValid"), err, "wrong config is valid")
}

type testSubPtrConfigIsValid struct {
	SubConfig *testConfigIsValid
}

func TestSubConfigIsValid(t *testing.T) {
	var c testSubPtrConfigIsValid
	err := configor.Load(ctx, &c)
	assert.Equal(t, errors.New("value &{false} invalid for type *config.testConfigIsValid"), err, "wrong during test sub config struct pointer is valid")
}

type testStructConfigIsValid struct {
	Invalid bool `config:"default:false"`
}

func (t testStructConfigIsValid) IsValid() bool {
	return false
}

type testSubStructConfigIsValid struct {
	SubConfig testStructConfigIsValid
}

func TestSubStructConfigIsValid(t *testing.T) {
	var c testSubStructConfigIsValid
	err := configor.Load(ctx, &c)
	assert.Equal(t, errors.New("value &{false} invalid for type *config.testStructConfigIsValid"), err, "wrong during test sub config struct is valid")
}

type testUnaddressConfig struct {
	unaddress string
}

func TestUnaddressConfig(t *testing.T) {
	var c testUnaddressConfig
	err := configor.Load(ctx, c)
	assert.Equal(t, errors.New("Config config.testUnaddressConfig should be addressable"), err, "wrong during test uaddress config is valid")
}

type testUnaddressConfigParam struct {
	unaddress string
}

func TestUnaddressConfigParam(t *testing.T) {
	var c testUnaddressConfigParam
	err := configor.Load(ctx, &c)
	assert.Equal(t, nil, err, "wrong during test uaddress config is valid")
	assert.Equal(t, testUnaddressConfigParam{}, c, "wrong during test uaddress config is valid")
}

type valueStruct struct {
	Name string
}

type testValueStruct struct {
	ValueStruct valueStruct `config:"name:value_struct;"`
}

func TestValueStruct(t *testing.T) {
	var c testValueStruct
	os.Setenv("value_struct", "{\"name\":\"test\"}")
	err := configor.Load(ctx, &c)
	assert.Equal(t, nil, err, "wrong during test unaddress config is valid")
	assert.Equal(t, testValueStruct{ValueStruct: valueStruct{Name: "test"}}, c, "wrong during test unaddress config is valid")
}
