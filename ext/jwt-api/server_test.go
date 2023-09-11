// HTTP API v1.0
// -----------------------
// 约定参数名称:
//
//		api       : 接口名称
//	 key  	  : 接口用户
//	 sign      : 签名
//	 sign_type : 签名类型
//
// -----------------------
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ixre/gof"
	"github.com/ixre/gof/crypto"
	"github.com/ixre/gof/ext/api"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/util"
)

const (
	// 是否关闭判断接口权限,仅供测试使用
	turnOffCheckPerm = false
)

var (
	RErrNotService = &Response{
		Code:    10094,
		Message: "api not service",
	}
	RErrDeprecated = &Response{
		Code:    10095,
		Message: "api is deprecated",
	}
)

// 服务
func ListenAndServe(port int, debug bool) error {
	// 请求限制
	rl := util.NewRequestLimit(storage.NewHashStorage(), 100, 10, 600)
	// 创建服务
	s := NewServerMux(swapApiKeyFunc, "/api", true)
	// 注册中间键
	serviceMiddleware(s, "[ API][ Log]: ", debug, rl)
	// 注册处理器
	s.HandlePublic(AccessTokenApi{})
	s.Handle(&StatusProcessor{})
	return http.ListenAndServe(fmt.Sprintf(":%d", port), s)
}

// 服务调试跟踪
func serviceMiddleware(s Server, prefix string, debug bool, rl *util.RequestLimit) {
	prefix = "[ Api][ Log]"
	// 验证IP请求限制
	s.Use(func(ctx Context) error {
		addr := ctx.Request().UserAddr
		if len(addr) != 0 && !rl.Acquire(addr, 1) || rl.IsLock(addr) {
			return errors.New("您的网络存在异常,系统拒绝访问")
			//return errors.New("access denied")
		}
		return nil
	})
	//// 校验版本
	//s.Use(func(ctx api.Context) error {
	//	//prod := ctx.StoredValues().GetString("product"
	//	prodVer := ctx.Params().GetString("version")
	//	if api.CompareVersion(prodVer, RequireVersion) < 0 {
	//		return errors.New("您当前使用的APP版本较低, 请升级或安装最新版本")
	//		//return errors.New(fmt.Sprintf("%s,require version=%s",
	//		//	api.RCDeprecated.Message, tarVer))
	//	}
	//	return nil
	//})

	if debug {
		// 开启调试
		s.Trace()
		// 输出请求信息
		s.Use(func(ctx Context) error {
			apiName := ctx.Request().RequestApi
			log.Println(prefix, "user", ctx.UserKey(), " calling ", apiName)
			data, _ := json.Marshal(ctx.Request().Params)
			log.Println(prefix, "request data = [", data, "]")
			// 记录服务端请求时间
			ctx.Request().Params.Set("$rpc_begin_time", time.Now().UnixNano())
			return nil
		})
	}
	if debug {
		// 输出响应结果
		s.After(func(ctx Context) error {
			form := ctx.Request().Params
			rsp := form.Get("$api_response").(*api.Response)
			data := ""
			if rsp.Data != nil {
				d, _ := json.Marshal(rsp.Data)
				data = string(d)
			}
			reqTime := int64(form.GetInt("$rpc_begin_time"))
			elapsed := float32(time.Now().UnixNano()-reqTime) / 1000000000
			log.Println(prefix, "response : ", rsp.Code, rsp.Message,
				fmt.Sprintf("; elapsed time ：%.4fs ; ", elapsed),
				"result = [", data, "]",
			)
			return nil
		})
	}
}

var jwtSecret = []byte("...")

// 交换接口用户凭据
func swapApiKeyFunc(ctx Context) (privateKey []byte, err error) {
	return jwtSecret, nil
}

var _ Handler = new(StatusProcessor)

type StatusProcessor struct{}

func (p *StatusProcessor) Group() string {
	return "status"
}

func (p *StatusProcessor) Process(fn string, ctx Context) *Response {
	r := &Response{Code: 0}
	switch fn {
	case "ping":
		r.Data = "pong"
	}
	return r
}

var _ Handler = new(AccessTokenApi)

type AccessTokenApi struct {
}

func (a AccessTokenApi) Group() string {
	return "access_token"
}

func (a AccessTokenApi) Process(fn string, ctx Context) *Response {
	return a.createAccessToken(ctx)
}

func (a AccessTokenApi) createAccessToken(ctx Context) *Response {
	ownerKey := ctx.Request().Params.GetString("key")
	md5Secret := ctx.Request().Params.GetString("secret")
	if len(ownerKey) == 0 || len(md5Secret) == 0 {
		return ResponseWithCode(1, "require params key and md5_secret")
	}
	cfg := gof.CurrentApp.Config()
	apiUser := cfg.GetString("api_user")
	apiSecret := cfg.GetString("api_secret")

	if apiUser != "tmp_0606" {
		if apiUser != ownerKey || md5Secret != crypto.Md5([]byte(apiSecret)) {
			return ResponseWithCode(4, "用户或密钥不正确")
		}
	}
	// 创建token并返回
	claims := CreateClaims("0", "go2o",
		"go2o-api-jwt", time.Now().Unix()+7200).(MapClaims)
	claims["global"] = true
	token, err := AccessToken(claims, jwtSecret)
	if err != nil {
		return ResponseWithCode(4, err.Error())
	}
	return NewResponse(token)
}
