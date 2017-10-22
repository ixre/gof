package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

// 创建接口签名
func TestGenApiSign(t *testing.T) {
	key := "test"
	secret := "123456"
	signType := "sha1"
	serverUrl := "http://localhost:7020/api"
	form := url.Values{
		"key":  []string{key},
		"api":       []string{"status.ping,status.hello"},
		"sign_type": []string{signType},
	}
	sign := Sign(signType, form, secret)
	//t.Log("-- Sign:", sign)
	form["sign"] = []string{sign}
	cli := http.Client{}
	rsp, err := cli.PostForm(serverUrl, form)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	data, _ := ioutil.ReadAll(rsp.Body)
	rsp1 := Response{}
	json.Unmarshal(data, &rsp1)
	if rsp1.Code != StatusOK {
		t.Log("请求失败：code:", rsp1.Code, "; message:", rsp1.Message)
		t.Log("接口响应：", string(data))
		t.FailNow()
	}
	t.Log("接口响应：", string(data))
}
