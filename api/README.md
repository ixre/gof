# HTTP API

## 接口加密规则

接口采用POST请求，将参数集合按字母排序后，排除sign_type，拼接token
然后进行MD5或SHA1进行加密得到sign,并将sign添加到请求参数集合。

## 接口请求示例代码(go)
    
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
    form["sign"] = []string{sign}
    cli := http.Client{}
    rsp, err := cli.PostForm(serverUrl, form)
    if err == nil {	
        data, _ := ioutil.ReadAll(rsp.Body)
        log.Println("接口响应：", string(data))
    }
    


