package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"trinity-micro/core/e"

	"github.com/sirupsen/logrus"
)

type HttpClient interface {
	Call(ctx context.Context, logger logrus.FieldLogger, method string, url string, header map[string]string, body interface{}, dest interface{}, requestInterceptors ...Interceptor) error
}

type clientImpl struct {
	logger logrus.FieldLogger
}

type Interceptor func(r *http.Request) error

// Request client
// will return the err when the response code is not >=200 or <= 300
// will decode the response to dest even it return error
func (c *clientImpl) Call(ctx context.Context, logger logrus.FieldLogger, method string, url string, header map[string]string, body interface{}, dest interface{}, requestInterceptors ...Interceptor) error {
	now := time.Now()
	logger.Infof("HttpClient.Call start , method: %v , url: %v ", method, url)
	defer logger.Infof("HttpClient.Call ended ,method: %v, url: %v, duration: %v ", method, url, time.Since(now))
	var bodyTemp io.Reader
	if body != nil {
		r, ok := body.(io.Reader)
		if ok {
			bodyTemp = r
		} else {
			var bodyBytes []byte
			mime := header[HeaderMime]
			switch mime {
			case MimeXML, MimeTextXML:
				var err error
				bodyBytes, err = xml.Marshal(body)
				if err != nil {
					return e.NewError(e.Info, e.ErrInvalidRequest, "encode xml error", err)
				}
			default:
				var err error
				bodyBytes, err = json.Marshal(body)
				if err != nil {
					return e.NewError(e.Info, e.ErrInvalidRequest, "encode json error", err)
				}
			}
			logger.Infof("HttpClient.Call request body: %v ", string(bodyBytes))
			bodyTemp = bytes.NewReader(bodyBytes)
		}
	}
	req, err := http.NewRequest(method, url, bodyTemp)
	if err != nil {
		return e.NewError(e.Info, e.ErrInvalidRequest, "new request error ", err)
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	for _, interceptor := range requestInterceptors {
		if err := interceptor(req); err != nil {
			return e.NewError(e.Info, e.ErrInvalidRequest, "new request interceptor error ", err)
		}
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return e.NewError(e.Info, e.ErrInternalServer, "new request error ", err)
	}
	defer resp.Body.Close()
	bodyRes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return e.NewError(e.Error, e.ErrReadResponseBody, err.Error())
	}
	logger.Infof("HttpClient.Call response code:%v , body: %v ", resp.StatusCode, string(bodyRes))

	resbody := ioutil.NopCloser(bytes.NewReader(bodyRes))
	mime := resp.Header.Get(HeaderMime)
	contextType := strings.Split(mime, ";")[0]
	switch contextType {
	case MimeTextXML, MimeXML:
		if err := xml.NewDecoder(resbody).Decode(dest); err != nil {
			return e.NewError(e.Info, e.ErrDecodeResponseBody, "decode xml error", err)
		}
	default:
		if err := json.NewDecoder(resbody).Decode(dest); err != nil {
			return e.NewError(e.Info, e.ErrDecodeResponseBody, "decode json error", err)
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return e.NewError(e.Info, e.ErrHttpResponseCodeNotSuccess, string(bodyRes), fmt.Errorf("actual http response code: %v", resp.StatusCode))
	}

	return nil
}

func NewClient() HttpClient {
	return &clientImpl{}
}
