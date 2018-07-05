// HTTP API v1.0
// -----------------------
// 约定参数名称:
//	api       : 接口名称
//  key  	  : 接口用户
//  sign      : 签名
//  sign_type : 签名类型
// -----------------------
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	// 是否关闭判断接口权限,仅供测试使用
	turnOffCheckPerm = false
)

var (
	RErrNotService = &Response{
		Code:   10094,
		ErrMsg: "api not service",
	}
	RErrDeprecated = &Response{
		Code:   10095,
		ErrMsg: "api is deprecated",
	}
)

// 服务
func ListenAndServe(port int, debug bool) error {
	// 创建上下文工厂
	factory := DefaultFactory.Build(nil)
	// 创建服务
	s := NewServerMux(factory, apiSwapFunc)
	hs := http.NewServeMux()
	hs.Handle("/api", s)
	hs.Handle("/api_v1", s)
	// 中间件
	tarVer := "1.0"
	// 校验版本
	s.Use(func(ctx Context) error {
		prod := ctx.Form().GetString("product")
		prodVer := ctx.Form().GetString("version")
		if prod == "mzl" && CompareVersion(prodVer, tarVer) < 0 {
			return errors.New(fmt.Sprintf("%d:%s,require version=%s",
				RErrDeprecated.Code, RErrDeprecated.ErrMsg, tarVer))
		}
		return nil
	})
	if debug {
		// 开启调试
		s.Trace()
		// 输出请求信息
		s.Use(func(ctx Context) error {
			apiName := ctx.Form().Get("$api_name").(string)
			log.Println("[ Api][ Log]: user", ctx.Key(), " calling ", apiName)
			data, _ := url.QueryUnescape(ctx.Request().Form.Encode())
			log.Println("[ Api][ Log]: request data = [", data, "]")
			// 记录服务端请求时间
			ctx.Form().Set("$rpc_begin_time", time.Now().UnixNano())
			return nil
		})
		// 输出响应结果
		s.After(func(ctx Context) error {
			form := ctx.Form()
			rsp := form.Get("$api_response").(*Response)
			data := ""
			if rsp.Data != nil {
				d, _ := json.Marshal(rsp.Data)
				data = string(d)
			}
			reqTime := int64(ctx.Form().GetInt("$rpc_begin_time"))
			elapsed := float32(time.Now().UnixNano()-reqTime) / 1000000000
			log.Println("[ Api][ Log]: response : ", rsp.Code, rsp.ErrMsg,
				fmt.Sprintf("; elapsed time ：%.4fs ; ", elapsed),
				"result = [", data, "]",
			)
			if rsp.Code == RAccessDenied.Code {
				data, _ := url.QueryUnescape(ctx.Request().Form.Encode())
				sortData := ParamsToBytes(ctx.Request().Form, form.GetString("$user_secret"))
				log.Println("[ Api][ Log]: request data = [", data, "]")
				log.Println(" sign not match ! key =", form.Get("key"),
					"\r\n   server_sign=", form.GetString("$server_sign"),
					"\r\n   client_sign=", form.GetString("$client_sign"),
					"\r\n   sort_params=", string(sortData))
			}
			return nil
		})
	}

	// 注册处理器
	s.Register("status", &StatusProcessor{})
	return http.ListenAndServe(fmt.Sprintf(":%d", port), s)
}

// 检查接口参数
func apiSwapFunc(key string) (userId int64, userToken string, checkSign bool) {
	//todo: return user info
	if turnOffCheckPerm {
		return 1, "123456", false
	}
	if key == "test" {
		return 1, "123456", true
	}
	if key == "80line365mzl00" {
		return 1, "239d2d9fb16dbe18d81c54d1764bd33b", true
	}
	return 0, "", false
}

func CompareVersion(v, v1 string) int {
	return intVer(v) - intVer(v1)
}
func intVer(s string) int {
	arr := strings.Split(s, ".")
	for i, v := range arr {
		if l := len(v); l < 3 {
			arr[i] = strings.Repeat("0", 3-l) + v
		}
	}
	intVer, err := strconv.Atoi(strings.Join(arr, ""))
	if err != nil {
		panic(err)
	}
	return intVer
}

var _ Processor = new(StatusProcessor)

type StatusProcessor struct{}

func (p *StatusProcessor) Request(fn string, ctx Context) *Response {
	r := &Response{Code: CodeOK}
	switch fn {
	case "ping":
		r.Data = "pong"
	}
	return r
}
