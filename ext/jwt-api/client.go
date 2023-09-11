package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// 错误响应处理函数
type ErrRspFunc func(code int, message string) error

type Client struct {
	server        string // 	= "http://localhost:1419/openapi"
	f             AccessTokenFunc
	accessToken   string
	lastTokenUnix int64
	expires       int
	headerKey     string
	key           string //     = "< replace your api user >"
	secret        string //   	= "< replace your api token >"
	errRsp        ErrRspFunc
}

// fetch access token function
type AccessTokenFunc func(key, secret string) string

// 创建新的客户端
func NewClient(server string, key string, secret string) *Client {
	return &Client{
		server:    server,
		key:       key,
		secret:    secret,
		headerKey: "Authorization",
		errRsp: func(code int, text string) error {
			// 如果返回接口请求错误, 响应状态码以10开头
			return errors.New(fmt.Sprintf("Error code %d : %s", code, text))
		},
	}
}

func (c *Client) UseToken(f AccessTokenFunc, tokenExpires int) {
	c.f = f
	c.expires = tokenExpires
}

func (c *Client) HandleError(errFunc ErrRspFunc) {
	c.errRsp = errFunc
}

func (c *Client) Get(apiPath string, data interface{}) ([]byte, error) {
	return c.Request(apiPath, "POST", data, time.Duration(0))
}

// 请求接口
func (c *Client) Post(apiPath string, data interface{}) ([]byte, error) {
	return c.Request(apiPath, "POST", data, time.Duration(0))
}

// 请求接口
func (c *Client) PUT(apiPath string, data interface{}) ([]byte, error) {
	return c.Request(apiPath, "PUT", data, time.Duration(0))
}

// 请求接口
func (c *Client) Delete(apiPath string, data interface{}) ([]byte, error) {
	return c.Request(apiPath, "DELETE", data, time.Duration(0))
}

// 请求接口
func (c *Client) Patch(apiPath string, data interface{}) ([]byte, error) {
	return c.Request(apiPath, "PATCH", data, time.Duration(0))
}

func (c *Client) Request(apiPath string, method string, data interface{}, timeout time.Duration) ([]byte, error) {
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
		c.accessToken = c.f(c.key, c.secret)
		c.lastTokenUnix = now
	}
	req.Header.Add(c.headerKey, c.accessToken)
	req.Header.Add("user-key", c.key)
	if len(contentType) > 0 {
		req.Header.Set("Content-Type", contentType)
		//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	rsp, err := cli.Do(req)
	if err == nil {
		data, err := io.ReadAll(rsp.Body)
		if err == nil {
			if len(data) >= 1 && data[0] == '#' {
				ret := string(data[1:])
				i := strings.Index(ret, "#")
				code, _ := strconv.Atoi(ret[:i])
				return data, c.errRsp(code, ret[i+1:])
			}
			return data, nil
		}
	}
	return []byte{}, err
}

func (c Client) path(path string) string {
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	return c.server + path
}

func (c *Client) parseBody(data interface{}) (io.Reader, string, error) {
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
