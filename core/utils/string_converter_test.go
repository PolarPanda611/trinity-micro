/*
 * @Author: Daniel TAN
 * @Description:
 * @Date: 2021-08-06 10:12:19
 * @LastEditTime: 2021-08-26 10:52:08
 * @LastEditors: Daniel TAN
 * @FilePath: /fr-price-common-pkg/core/utils/string_converter_test.go
 */
package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StringConverter(t *testing.T) {
	type Test struct {
		S1 string
		S2 int64
		S3 int32
		S4 int16
		S5 int
		S6 bool
		S7 interface{}
		S8 struct {
			Test string
			Code int
		}
		S9  *string
		S10 float32
		S11 float64
	}
	dest := &Test{}
	{
		val := getStructFieldValue(dest, 0)
		err := StringConverter("123", &val)
		assert.Equal(t, nil, err, "wrong err ")
		assert.Equal(t, "123", dest.S1, "wrong dest ")
	}
	{
		val := getStructFieldValue(dest, 1)
		err := StringConverter("123", &val)
		assert.Equal(t, nil, err, "wrong err ")
		assert.Equal(t, int64(123), dest.S2, "wrong dest ")
	}
	{
		val := getStructFieldValue(dest, 2)
		err := StringConverter("123", &val)
		assert.Equal(t, nil, err, "wrong err ")
		assert.Equal(t, int32(123), dest.S3, "wrong dest ")
	}
	{
		val := getStructFieldValue(dest, 3)
		err := StringConverter("123", &val)
		assert.Equal(t, nil, err, "wrong err ")
		assert.Equal(t, int16(123), dest.S4, "wrong dest ")
	}
	{
		val := getStructFieldValue(dest, 4)
		err := StringConverter("123", &val)
		assert.Equal(t, nil, err, "wrong err ")
		assert.Equal(t, int(123), dest.S5, "wrong dest ")
	}
	{
		val := getStructFieldValue(dest, 5)
		err := StringConverter("true", &val)
		assert.Equal(t, nil, err, "wrong err ")
		assert.Equal(t, true, dest.S6, "wrong dest ")
	}
	{
		val := getStructFieldValue(dest, 1)
		err := StringConverter("aaa", &val)
		assert.Equal(t, "strconv.ParseInt: parsing \"aaa\": invalid syntax", err.Error(), "wrong err ")
	}
	{
		val := getStructFieldValue(dest, 2)
		err := StringConverter("aaa", &val)
		assert.Equal(t, "strconv.ParseInt: parsing \"aaa\": invalid syntax", err.Error(), "wrong err ")
	}
	{
		val := getStructFieldValue(dest, 3)
		err := StringConverter("aaa", &val)
		assert.Equal(t, "strconv.ParseInt: parsing \"aaa\": invalid syntax", err.Error(), "wrong err ")
	}
	{
		val := getStructFieldValue(dest, 4)
		err := StringConverter("aaa", &val)
		assert.Equal(t, "strconv.Atoi: parsing \"aaa\": invalid syntax", err.Error(), "wrong err ")
	}
	{
		val := getStructFieldValue(dest, 6)
		err := StringConverter("true", &val)
		assert.Equal(t, nil, err, "wrong err ")
		assert.Equal(t, (interface{})("true"), dest.S7, "wrong dest ")
	}
	{
		val := getStructFieldValue(dest, 7)
		err := StringConverter(`{"Test":"123","Code":1}`, &val)
		assert.Equal(t, nil, err, "wrong err ")
		assert.Equal(t, struct {
			Test string
			Code int
		}{Test: "123", Code: 1}, dest.S8, "wrong dest ")
	}
	{
		val := getStructFieldValue(dest, 7)
		err := StringConverter(`{"Test":"123","Code":1}123`, &val)
		assert.Equal(t, "invalid character '1' after top-level value", err.Error(), "wrong err ")
	}
	{
		val := getStructFieldValue(dest, 8)
		err := StringConverter("true", &val)
		assert.Equal(t, "unsupported type", err.Error(), "wrong err ")
	}
	{
		val := getStructFieldValue(dest, 9)
		err := StringConverter("3.4", &val)
		assert.Equal(t, nil, err, "wrong err ")
		assert.Equal(t, float32(3.4), dest.S10, "wrong dest ")
	}
	{
		val := getStructFieldValue(dest, 10)
		err := StringConverter("3.4", &val)
		assert.Equal(t, nil, err, "wrong err ")
		assert.Equal(t, float64(3.4), dest.S11, "wrong dest ")
	}
}

func getStructFieldValue(dest interface{}, index int) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(dest)).Field(index)
}
