# HTTP API
接口通过HTTP暴露给外部应用，采用签名鉴权，接口请求方需要使用KEY和SECRET方可调用>接口，SECRET用于参数签名。

## 接口测试数据

        接口地址：http://localhost:7020/api
        接口KEY：test
        接口SECRET: 123456
        
## 接口签名规则

接口采用POST请求，将参数集合按字母排序后，排除sign_type，拼接token
然后进行MD5或SHA1进行加密得到sign,并将sign添加到请求参数集合。

签名示例代码(go):
    
    // 参数排序后，排除sign和sign_type，拼接token，转换为字节
    func paramsToBytes(r url.Values, token string) []byte {
	    i := 0
	    buf := bytes.NewBuffer(nil)
	    // 键排序
	    keys := []string{}
	    for k, _ := range r {
		    keys = append(keys, k)
	    }
	    sort.Strings(keys)
	    // 拼接参数和值
	    for _, k := range keys {
		    if k == "sign" || k == "sign_type" {
			    continue
		    }
		    if i > 0 {
			    buf.WriteString("&")
		    }
		    buf.WriteString(k)
		    buf.WriteString("=")
		    buf.WriteString(r[k][0])
		    i++
	    }
	    buf.WriteString(token)
	    return buf.Bytes()
    }

    // 签名
    func Sign(signType string, r url.Values, token string) string {
	    data := paramsToBytes(r, token)
	    switch signType {
	    case "md5":
		    return md5Encode(data)
	    case "sha1":
		    return sha1Encode(data)
	    }
	    return ""
    }

    // MD5加密
    func md5Encode(data []byte) string {
	    m := md5.New()
	    m.Write(data)
	    dec := m.Sum(nil)
	    return hex.EncodeToString(dec)
    }
    // SHA1加密
    func sha1Encode(data []byte) string {
	    s := sha1.New()
	    s.Write(data)
	    d := s.Sum(nil)
	    return hex.EncodeToString(d)
    }


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
    // sign = "fe343c958b61178b3644432263cf819c153569ed"
    form["sign"] = []string{sign}
    cli := http.Client{}
    rsp, err := cli.PostForm(serverUrl, form)
    if err == nil {	
        data, _ := ioutil.ReadAll(rsp.Body)
        log.Println("接口响应：", string(data))
    }
    


