/*
 * @Author: Daniel TAN
 * @Description:
 * @Date: 2021-08-18 09:27:53
 * @LastEditTime: 2021-08-18 09:28:00
 * @LastEditors: Daniel TAN
 * @FilePath: /fr-price-common-pkg/core/utils/type_converter.go
 */
package utils

import (
	"encoding/json"
	"encoding/xml"
	"strconv"
	"time"
)

func MustXMLMarshal(obj interface{}) string {
	content, err := xml.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func MustJSONMarshal(obj interface{}) string {
	content, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func MustFormatTime(layout string, timeString string) *time.Time {
	t, err := time.Parse(layout, timeString)
	if err != nil {
		panic(err)
	}
	return &t
}

func MustParseBool(boolString string) bool {
	b, err := strconv.ParseBool(boolString)
	if err != nil {
		panic(err)
	}
	return b
}

func MustParseFloat64(floatString string) float64 {
	if floatString == "" {
		return 0
	}
	f, err := strconv.ParseFloat(floatString, 64)
	if err != nil {
		panic(err)
	}
	return f
}
