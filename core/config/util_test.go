package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeConverterBool(t *testing.T) {
	// string
	r := "true"
	var bTrue bool
	res, err := typeConverter(r, reflect.TypeOf(bTrue))
	assert.Equal(t, nil, err, "err during load config")
	assert.Equal(t, true, res, "err during load config")

	r2 := "false"
	var bFalse bool
	res, err = typeConverter(r2, reflect.TypeOf(bFalse))
	assert.Equal(t, nil, err, "err during load config")
	assert.Equal(t, false, res, "err during load config")

}

func TestTypeConverterFloat(t *testing.T) {
	// string
	r := "1"
	var aInt float32
	res, err := typeConverter(r, reflect.TypeOf(aInt))
	assert.Equal(t, nil, err, "err during load config")
	assert.Equal(t, float32(1), res, "err during load config")

	var aInt16 float64
	res, err = typeConverter(r, reflect.TypeOf(aInt16))
	assert.Equal(t, nil, err, "err during load config")
	assert.Equal(t, float64(1), res, "err during load config")

}

func TestTypeConverterUint(t *testing.T) {
	// string
	r := "1"
	var aInt uint
	res, err := typeConverter(r, reflect.TypeOf(aInt))
	assert.Equal(t, nil, err, "err during load config")
	assert.Equal(t, uint(1), res, "err during load config")

	var aInt16 uint16
	res, err = typeConverter(r, reflect.TypeOf(aInt16))
	assert.Equal(t, nil, err, "err during load config")
	assert.Equal(t, uint16(1), res, "err during load config")

	var aInt32 uint32
	res, err = typeConverter(r, reflect.TypeOf(aInt32))
	assert.Equal(t, nil, err, "err during load config")
	assert.Equal(t, uint32(1), res, "err during load config")

	var aInt64 uint64
	res, err = typeConverter(r, reflect.TypeOf(aInt64))
	assert.Equal(t, nil, err, "err during load config")
	assert.Equal(t, uint64(1), res, "err during load config")

}

func TestTypeConverterInt(t *testing.T) {
	// string
	r := "1"
	var aInt int
	res, err := typeConverter(r, reflect.TypeOf(aInt))
	assert.Equal(t, nil, err, "err during load config")
	assert.Equal(t, int(1), res, "err during load config")

	var aInt16 int16
	res, err = typeConverter(r, reflect.TypeOf(aInt16))
	assert.Equal(t, nil, err, "err during load config")
	assert.Equal(t, int16(1), res, "err during load config")

	var aInt32 int32
	res, err = typeConverter(r, reflect.TypeOf(aInt32))
	assert.Equal(t, nil, err, "err during load config")
	assert.Equal(t, int32(1), res, "err during load config")

	var aInt64 int64
	res, err = typeConverter(r, reflect.TypeOf(aInt64))
	assert.Equal(t, nil, err, "err during load config")
	assert.Equal(t, int64(1), res, "err during load config")

}

func TestTypeConverterString(t *testing.T) {
	// string
	r := "1"
	var a string
	res, err := typeConverter(r, reflect.TypeOf(a))
	assert.Equal(t, nil, err, "err during load config")
	assert.Equal(t, "1", res, "err during load config")

}
