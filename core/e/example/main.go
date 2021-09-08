/*
 * @Author: Daniel TAN
 * @Date: 2021-09-03 12:24:12
 * @LastEditors: Daniel TAN
 * @LastEditTime: 2021-09-09 00:08:32
 * @FilePath: /trinity-micro/core/e/example/main.go
 * @Description:
 */
package main

import (
	"fmt"
	"trinity-micro/core/e"

	"log"

	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
)

var (
	l = logrus.New()
)

func main() {
	{
		// example 1:
		//  age should > 0 < 30
		qty := 300
		err := e.NewError(e.Info, e.ErrInvalidRequest, "wrong quantity", fmt.Errorf("expected >0 < 30 , actual %v", qty))
		// INFO[0000] loglevel: Info, error code: 400001, error type: InvalidRequest, error message: wrong quantity, actual error: expected >0 < 30 , actual 300
		e.Logging(l, err)
		// 2021/08/30 17:53:02 http response status code: 400
		// 2021/08/30 17:53:02 http response error code: 400001
		// 2021/08/30 17:53:02 http response error type: InvalidRequest
		// 2021/08/30 17:53:02 http response error message: wrong quantity
		logResponse(err)
	}

	{
		// example 2 :
		// db record not found
		id := 233
		dbError := gorm.ErrRecordNotFound
		err := e.NewError(e.Info, e.ErrResourceNotFound, "resource xxx not found", fmt.Errorf("id: %v not found , actual error: %v ", id, dbError))
		e.Logging(l, err)
		logResponse(err)
	}

	{
		// example 3 :
		// db unknown error
		dbError := gorm.Err
		err := e.NewError(e.Error, e.ErrInternalServer, "unexpected db error ", dbError)
		e.Logging(l, err)
		logResponse(err)
	}
}

func logResponse(err error) {
	httpResponse1 := e.NewAPIError(err)
	log.Printf("http response status code: %v", httpResponse1.Status)
	log.Printf("http response error code: %v", httpResponse1.Code)
	log.Printf("http response error type: %v", httpResponse1.Type)
	log.Printf("http response error message: %v", httpResponse1.Message)

}
