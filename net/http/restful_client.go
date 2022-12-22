package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 错误响应处理函数
type ErrorHandler func(code int, message string) error

type RestfulClient struct {
	url           string // 	= "http://localhost:1419/openapi"
	f             FetchTokenFunc
	accessToken   string
	lastTokenUnix int64
	expires       int
	headerKey     string
	errorHandler  ErrorHandler
}

// FetchTokenFunc fetch access token function
type FetchTokenFunc func() string

// NewRestfulClient 创建新的客户端
func NewRestfulClient(url string) *RestfulClient {
	return &RestfulClient{
		url:       url,
		headerKey: "Authorization",
		errorHandler: func(code int, text string) error {
			// 如果返回接口请求错误, 响应状态码以10开头
			return errors.New(fmt.Sprintf("Error code %d : %s", code, text))
		},
	}
}

func (c *RestfulClient) UseToken(f FetchTokenFunc, expires int) {
	c.f = f
	c.expires = expires
}

func (c *RestfulClient) HandleError(errFunc ErrorHandler) {
	c.errorHandler = errFunc
}

func (c *RestfulClient) Get(path string, data interface{}) ([]byte, error) {
	return c.Request(path, "POST", data, time.Duration(0))
}

// Post 请求接口
func (c *RestfulClient) Post(path string, data interface{}) ([]byte, error) {
	return c.Request(path, "POST", data, time.Duration(0))
}

// PUT 请求接口
func (c *RestfulClient) PUT(apiPath string, data interface{}) ([]byte, error) {
	return c.Request(apiPath, "PUT", data, time.Duration(0))
}

// Delete 请求接口
func (c *RestfulClient) Delete(path string, data interface{}) ([]byte, error) {
	return c.Request(path, "DELETE", data, time.Duration(0))
}

// Patch 请求接口
func (c *RestfulClient) Patch(path string, data interface{}) ([]byte, error) {
	return c.Request(path, "PATCH", data, time.Duration(0))
}

func (c *RestfulClient) Request(apiPath string, method string, data interface{}, timeout time.Duration) ([]byte, error) {
	cli := &http.Client{}
	cli.Timeout = timeout
	reader, contentType, err := c.parseBody(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, c.path(apiPath), reader)
	if err != nil {
		return nil, err
	}
	var now = time.Now().Unix()
	if int(now-c.lastTokenUnix) > c.expires {
		c.accessToken = c.f()
		c.lastTokenUnix = now
	}
	req.Header.Add(c.headerKey, c.accessToken)
	if len(contentType) > 0 {
		req.Header.Set("Content-Type", contentType)
		//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	rsp, err := cli.Do(req)
	if err == nil {
		data, err := io.ReadAll(rsp.Body)
		if err == nil {
			return data, nil
		}
	}
	return []byte{}, c.errorHandler(rsp.StatusCode, err.Error())
}

func (c *RestfulClient) path(path string) string {
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	return c.url + path
}

func (c *RestfulClient) parseBody(data interface{}) (io.Reader, string, error) {
	if data == nil {
		return nil, "", nil
	}
	if d, ok := data.(url.Values); ok {
		return strings.NewReader(d.Encode()), "", nil
	}
	j, err := json.Marshal(data)
	if err != nil {
		return nil, "", err
	}
	return bytes.NewReader(j), "application/json", err
}
