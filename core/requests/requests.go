// Author: Daniel TAN
// Date: 2021-09-03 12:24:12
// LastEditors: Daniel TAN
// LastEditTime: 2021-09-28 01:12:13
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

	"github.com/PolarPanda611/trinity-micro/core/e"
)

type Interceptor func(r *http.Request) error

// Call
// will return the err when the response code is not >=200 or <= 300
// will decode the response to dest even it return error
func Call(ctx context.Context, method string, url string, header http.Header, body interface{}, dest interface{}, requestInterceptors ...Interceptor) error {
	var bodyTemp io.Reader
	if body != nil {
		r, ok := body.(io.Reader)
		if ok {
			bodyTemp = r
		} else {
			var bodyBytes []byte
			mime := header.Get(HeaderMime)
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
			bodyTemp = bytes.NewReader(bodyBytes)
		}
	}
	req, err := http.NewRequest(method, url, bodyTemp)
	if err != nil {
		return e.NewError(e.Info, e.ErrInvalidRequest, "new request error ", err)
	}
	req.Header = header
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

	resbody := ioutil.NopCloser(bytes.NewReader(bodyRes))
	mime := resp.Header.Get(HeaderMime)
	contextType := strings.Split(mime, ";")[0]
	switch contextType {
	case MimeTextXML, MimeXML:
		if err := xml.NewDecoder(resbody).Decode(dest); err != nil {
			return e.NewError(e.Info, e.ErrDecodeResponseBody, "decode xml error", err)
		}
	case MimeTextHTML:
		bodyHTML, err := ioutil.ReadAll(resbody)
		if err != nil {
			return e.NewError(e.Info, e.ErrDecodeResponseBody, "decode html error", err)
		}
		return e.NewError(e.Info, e.ErrDecodeResponseBody, fmt.Sprintf("decode html error, content: %v", string(bodyHTML)))
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
