package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func testApi(t *testing.T, apiName string, params url.Values) {
	key := "10000001"
	serverUrl := "http://localhost:1419/api"
	params["user_key"] = []string{key}
	cli := http.Client{}
	req, _ := http.NewRequest("POST", serverUrl+"/app/info", bytes.NewReader([]byte(params.Encode())))
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:82.0) Gecko/20100101 Firefox/82.0")
	rsp, err := cli.Do(req)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	data, _ := ioutil.ReadAll(rsp.Body)
	rsp1 := Response{}
	json.Unmarshal(data, &rsp1)
	if rsp1.Code != 0 {
		t.Log("请求失败：code:", rsp1.Code, "; message:", rsp1.Message)
		t.Log("接口响应：", string(data))
		t.FailNow()
	}
	t.Log("接口响应：", string(data))
}

//
//func (a *AppApi) appInfo(c api.Context, ownerId int) *api.Response {
//	claims,_ := c.Claims().(api.MapClaims)
//	println(claims["owner_code"].(string))
//	appCode := c.Params().GetString("app")
//	ownerCode := a.GetOwnerCode(c)
//	r, _ := service.OwnerService.GetApp(context.TODO(), ownerCode, appCode)
//	return api.NewResponse(r)
//}
