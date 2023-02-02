package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func testApi(t *testing.T, apiName string, params url.Values) {
	key := "test"
	secret := "123456"
	signType := "sha1"
	serverUrl := "http://localhost:7020/api"
	params["key"] = []string{key}
	params["api"] = []string{apiName}
	params["key"] = []string{key}
	params["sign_type"] = []string{signType}

	sign := Sign(signType, params, secret)
	//t.Log("-- Sign:", sign)
	params["sign"] = []string{sign}
	cli := http.Client{}
	rsp, err := cli.PostForm(serverUrl, params)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	data, _ := io.ReadAll(rsp.Body)
	rsp1 := Response{}
	json.Unmarshal(data, &rsp1)
	if rsp1.Code != RSuccessCode {
		t.Log("请求失败：code:", rsp1.Code, "; message:", rsp1.Message)
		t.Log("接口响应：", string(data))
		t.FailNow()
	}
	t.Log("接口响应：", string(data))
}

// 创建接口签名
func TestGenApiSign(t *testing.T) {
	key := "test"
	secret := "123456"
	signType := "sha1"
	serverUrl := "http://localhost:7020/api"
	form := url.Values{
		"key":          {key},
		"api":          {"status.ping,status.hello"},
		"product":      {"h"},
		"productType":  {"hello@$g"},
		"product_kind": {"h"},
		"sign_type":    {signType},
	}
	sign := Sign(signType, form, secret)
	t.Log("-- Sort params:", string(ParamsToBytes(form, secret, true)))
	t.Log("-- Sign:", sign)
	form["sign"] = []string{sign}
	cli := http.Client{}
	rsp, err := cli.PostForm(serverUrl, form)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	data, _ := io.ReadAll(rsp.Body)
	rsp1 := Response{}
	json.Unmarshal(data, &rsp1)
	if rsp1.Code != RSuccessCode {
		t.Log("请求失败：code:", rsp1.Code, "; message:", rsp1.Message)
		t.Log("接口响应：", string(data))
		t.FailNow()
	}
	t.Log("接口响应：", string(data))
}

func TestParamToBytes(t *testing.T) {
	form := url.Values{
		"Key":       {"sdf"},
		"api":       {"dsfsf"},
		"sign_type": {"sfsf"},
		"user":      {"jarrysix"},
		"Password":  {"2423424"},
		"loginType": {"normal"},
		"checkCode": {""},
	}

	t.Log("---xx = ", string(ParamsToBytes(form, "123", true)))
	form.Set("key", form.Get("Key"))
	form.Del("Key")
	t.Log("---xx = ", string(ParamsToBytes(form, "123", true)))
}

func TestSign(t *testing.T) {
	params := "api=member.login&key=go2o&product=app&pwd=c4ca4238a0b923820dcc509a6f75849b&user=18666398028&version=1.0.0&sign_type=sha1&sign=2933eaffccf9fe49a0ad9a97fe311a41afb6e3b2"
	values, _ := url.ParseQuery(params)
	sign := Sign("sha1", values, "131409")
	if sign2 := values.Get("sign"); sign2 != sign {
		println(sign, "/", sign2)
		t.Failed()
	}
}
