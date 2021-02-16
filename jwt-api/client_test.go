package api

import (
	"encoding/json"
	"errors"
	http2 "github.com/ixre/gof/util/http"
	"io/ioutil"
	"net/http"
	"testing"
)

var (
	tc *Client
)

var (
	RInternalError = &Response{
		Code:    RCInternalError,
		Message: "内部服务器出错",
	}
	RAccessDenied = &Response{
		Code:    RCAccessDenied,
		Message: "没有权限访问该接口",
	}
	RIncorrectApiParams = &Response{
		Code:    RCNotAuthorized,
		Message: "缺少接口参数，请联系技术人员解决",
	}
	RUndefinedApi = &Response{
		Code:    RCUndefinedApi,
		Message: "调用的API名称不正确",
	}
)

func init() {
	server := "http://localhost:1428/a/v2"
	tc = NewClient(server, "go2o", "123456")
	tc.UseToken(func(key, secret string) string {
		r, err1 := http.Get(server + "/access_token?key=" + key + "&secret=" + secret)
		if err1 != nil {
			println("---获取accessToken失败", err1.Error())
			return ""
		}
		bytes, _ := ioutil.ReadAll(r.Body)
		rsp := Response{}
		json.Unmarshal(bytes, &rsp)
		return rsp.Data.(string)
	}, 30000)
	tc.HandleError(func(code int, message string) error {
		switch code {
		case RCAccessDenied:
			message = RAccessDenied.Message
		case RCNotAuthorized:
			message = RIncorrectApiParams.Message
		case RCUndefinedApi:
			message = RUndefinedApi.Message
		}
		return errors.New(message)
	})
}

// 测试提交
func testPost(t *testing.T, apiName string, params map[string]string) ([]byte, error) {
	rsp, err := tc.Post(apiName, params)
	t.Log("[ Response]:", string(rsp))
	if err != nil {
		t.Error(err)
		//t.FailNow()
	}
	return rsp, err
}

// 测试提交
func testPostForm(t *testing.T, apiName string, params map[string]string) ([]byte, error) {
	rsp, err := tc.Post(apiName, params)
	t.Log("[ Response]:", string(rsp))
	if err != nil {
		t.Error(err)
		//t.FailNow()
	}
	return rsp, err
}

// 测试提交
func testGET(t *testing.T, apiName string, params map[string]string) ([]byte, error) {
	query := http2.ParseUrlValues(params).Encode()
	rsp, err := tc.Get(apiName+"?"+query, nil)
	t.Log("[ Response]:", string(rsp))
	if err != nil {
		t.Error(err)
		//t.FailNow()
	}
	return rsp, err
}

func TestReplaceSensitive(t *testing.T) {
	mp := map[string]string{
		"text":        "共产党是中华人民共和国的执政党",
		"replacement": "*",
	}
	testPost(t, "/fd/replace_sensitive", mp)
}
