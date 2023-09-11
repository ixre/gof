package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// 错误响应处理函数
type ErrRspFunc func(code int, message string) error

type Client struct {
	server   string // 	= "http://localhost:1419/openapi"
	key      string //     = "< replace your api user >"
	secret   string //   	= "< replace your api token >"
	signType string //		= "sha1" // [sha1|md5]
	errRsp   ErrRspFunc
}

// 创建新的客户端
func NewClient(server string, key string, secret string, signType string, errFunc ErrRspFunc) *Client {
	if errFunc == nil {
		// 如果返回接口请求错误, 响应状态码以10开头
		errFunc = func(code int, text string) error {
			return errors.New(fmt.Sprintf("Error code %d : %s", code, text))
		}
	}
	return &Client{
		server:   server,
		key:      key,
		secret:   secret,
		signType: signType,
		errRsp:   errFunc,
	}
}

// 请求接口
func (c Client) Post(apiName string, params map[string]string) ([]byte, error) {
	cli := &http.Client{}
	form := c.copy(params)
	form["api"] = []string{apiName}
	form["key"] = []string{c.key}
	form["sign_type"] = []string{c.signType}
	sign := Sign(c.signType, form, c.secret)
	form["sign"] = []string{sign}
	rsp, err := cli.PostForm(c.server, form)
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

func (c Client) copy(data map[string]string) map[string][]string {
	values := url.Values{}
	if data != nil {
		for k, v := range data {
			values[k] = []string{v}
		}
	}
	return values
}
