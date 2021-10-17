// Author: Daniel TAN
// Date: 2021-09-03 12:24:12
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-18 01:52:01
// FilePath: /trinity-micro/core/requests/requests.go
// Description:
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
)

type Interceptor func(r *http.Request) error

var _ Requests = new(HttpRequest)

type HttpRequest struct {
}

func NewRequest() Requests {
	return &HttpRequest{}
}

// Call
// will return the err when the response code is not >=200 or <= 300
// will decode the response to dest even it return error
func (r *HttpRequest) Call(ctx context.Context, method string, url string, header http.Header, body interface{}, dest interface{}, requestInterceptors ...Interceptor) error {
	var bodyTemp io.Reader
	if body != nil {
		r, ok := body.(io.Reader)
		if ok {
			bodyTemp = r
		} else {
			var bodyBytes []byte
			by, ok := body.([]byte)
			if ok {
				bodyBytes = by
			} else {
				mime := header.Get(HeaderMime)
				switch mime {
				case MimeXML, MimeTextXML:
					var err error
					bodyBytes, err = xml.Marshal(body)
					if err != nil {
						return fmt.Errorf("encode xml error, err: %v", err)
					}
				default:
					var err error
					bodyBytes, err = json.Marshal(body)
					if err != nil {
						return fmt.Errorf("encode json error, err: %v", err)
					}
				}
			}
			bodyTemp = bytes.NewReader(bodyBytes)
		}
	}
	req, err := http.NewRequest(method, url, bodyTemp)
	if err != nil {
		return fmt.Errorf("new request error, err: %v", err)
	}
	req.Header = header
	for _, interceptor := range requestInterceptors {
		if err := interceptor(req); err != nil {
			return fmt.Errorf("new request interceptor error, err: %v", err)
		}
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("do request error, err: %v", err)
	}
	defer resp.Body.Close()
	bodyRes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read request body error, err: %v", err)
	}
	resbody := ioutil.NopCloser(bytes.NewReader(bodyRes))
	mime := resp.Header.Get(HeaderMime)
	contextType := strings.Split(mime, ";")[0]
	switch contextType {
	case MimeTextXML, MimeXML:
		if err := xml.NewDecoder(resbody).Decode(dest); err != nil {
			return fmt.Errorf("decode xml error, err: %v", err)
		}
	case MimeTextHTML:
		bodyHTML, err := ioutil.ReadAll(resbody)
		if err != nil {
			return fmt.Errorf("read html error, err: %v", err)
		}
		return fmt.Errorf("html unsupported decode to destination, content: %v", string(bodyHTML))
	default:
		if err := json.NewDecoder(resbody).Decode(dest); err != nil {
			return fmt.Errorf("decode json error, err: %v", err)
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("actual http response code, actual: %v", resp.StatusCode)
	}

	return nil
}
