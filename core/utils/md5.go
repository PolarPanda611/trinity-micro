/*
 * @Author: Daniel TAN
 * @Description:
 * @Date: 2021-08-24 11:37:54
 * @LastEditTime: 2021-08-24 12:56:19
 * @LastEditors: Daniel TAN
 * @FilePath: /fr-price-common-pkg/core/utils/md5.go
 */
package utils

import (
	"crypto/md5"
	"encoding/base64"
	"strconv"
)

func MD5Content(content []byte) (string, error) {
	h := md5.New()
	if _, err := h.Write(content); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func GenerateNonce() (string, error) {
	// n, err := safe.RandUInt32()
	// if err != nil {
	// 	return "", err
	// }
	n := 1
	return strconv.FormatUint(uint64(n), 10), nil
}
