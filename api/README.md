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

## 接口返回
接口返回错误格式为:"#错误码#错误消息"，如：
```
#10091#api access denied
```


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
    
## 创建服务
创建处理程序
```
var _ api.Handler = new(MemberApi)
type MemberApi struct {}

func (m MemberApi) Process(fn string, ctx api.Context) *api.Response {
	return api.HandleMultiFunc(fn,ctx,map[string]api.HandlerFunc{
		"login":m.login,
	})
}

func (m MemberApi) login(ctx api.Context)interface{}{
	return 1
}
```
创建服务
```
// 服务
func NewServe(debug bool,version string) http.Handler {
	// 初始化变量
	registry := map[string]interface{}{}
	// 创建上下文工厂
	factory := api.DefaultFactory.Build(registry)
	serve := NewService(factory, version, debug)
	// 创建http处理器
	hs := http.NewServeMux()
	hs.Handle("/api", serve)
	return hs
}


// 服务
func NewService(factory api.ContextFactory, ver string, debug bool) *api.ServeMux {
	// 创建服务
	s := api.NewServerMux(factory, swapApiKeyFunc)
	// 注册处理器
	s.Register("member", &MemberApi{})
	//s.Register("dept", &DeptApi{})
	//s.Register("role", &RoleApi{})
	//s.Register("res", &ResApi{})
	//s.Register("user", &UserApi{})
	// 注册中间键
	serviceMiddleware(s, "[ Go2o][ API][ Log]: ", ver, debug)
	return s
}

// 服务调试跟踪
func serviceMiddleware(s api.Server, prefix string, tarVer string, debug bool) {
	prefix = "[ Api][ Log]"
	if debug {
		// 开启调试
		s.Trace()
		// 输出请求信息
		s.Use(func(ctx api.Context) error {
			apiName := ctx.Form().Get("$api_name").(string)
			log.Println(prefix, "user", ctx.Key(), " calling ", apiName)
			data, _ := url.QueryUnescape(ctx.Request().Form.Encode())
			log.Println(prefix, "request data = [", data, "]")
			// 记录服务端请求时间
			ctx.Form().Set("$rpc_begin_time", time.Now().UnixNano())
			return nil
		})
	}
	// 校验版本
	s.Use(func(ctx api.Context) error {
		//prod := ctx.FormData().GetString("product")
		prodVer := ctx.Form().GetString("version")
		if api.CompareVersion(prodVer, tarVer) < 0 {
			return errors.New(fmt.Sprintf("%s,require version=%s",
				api.RDeprecated.Message, tarVer))
		}
		return nil
	})

	if debug {
		// 输出响应结果
		s.After(func(ctx api.Context) error {
			form := ctx.Form()
			rsp := form.Get("$api_response").(*api.Response)
			data := ""
			if rsp.Data != nil {
				d, _ := json.Marshal(rsp.Data)
				data = string(d)
			}
			reqTime := int64(ctx.Form().GetInt("$rpc_begin_time"))
			elapsed := float32(time.Now().UnixNano()-reqTime) / 1000000000
			log.Println(prefix, "response : ", rsp.Code, rsp.Message,
				fmt.Sprintf("; elapsed time ：%.4fs ; ", elapsed),
				"result = [", data, "]",
			)
			if rsp.Code == api.RAccessDenied.Code {
				data, _ := url.QueryUnescape(ctx.Request().Form.Encode())
				sortData := api.ParamsToBytes(ctx.Request().Form, form.GetString("$user_secret"))
				log.Println(prefix, "request data = [", data, "]")
				log.Println(" sign not match ! key =", form.Get("key"),
					"\r\n   server_sign=", form.GetString("$server_sign"),
					"\r\n   client_sign=", form.GetString("$client_sign"),
					"\r\n   sort_params=", string(sortData))
			}
			return nil
		})
	}
}

// 交换接口用户凭据
func swapApiKeyFunc(ctx api.Context, key string) (userId int, userSecret string) {
	if key == "go2o"{
		return 1,"131409"
	}
	//log.Println(fmt.Sprintf("[ UAMS][ API]: 接口用户[%s]交换凭据失败： %s", key, r.ErrMsg))
	return 0, ""
}
```


