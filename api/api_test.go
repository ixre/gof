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
		"key":          []string{key},
		"api":          []string{"status.ping,status.hello"},
		"product":      []string{"h"},
		"productType":  []string{"hello"},
		"product_kind": []string{"h"},
		"sign_type":    []string{signType},
	}
	sign := Sign(signType, form, secret)
	t.Log("-- Sign:", sign)
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

func TestParamToBytes(t *testing.T) {
	form := url.Values{
		"Key":       []string{"sdf"},
		"api":       []string{"dsfsf"},
		"sign_type": []string{"sfsf"},
		"usr":       []string{"jarrysix"},
		"Pwd":       []string{"2423424"},
		"loginType": []string{"normal"},
		"checkCode": []string{""},
	}

	t.Log("---xx = ", string(paramsToBytes(form, "123")))
	form.Set("key", form.Get("Key"))
	form.Del("Key")
	t.Log("---xx = ", string(paramsToBytes(form, "123")))

}
