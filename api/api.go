// HTTP API v1.0
// -----------------------
// 约定参数名称:
//  key  	  : 接口用户
//	api       : 接口名称
//  sign      : 签名
//  sign_type : 签名类型[sha1|md5]
// -----------------------

package api

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/jsix/gof"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
)

// 接口响应
type Response struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var (
	StatusOK int64 = 10000
	RError         = &Response{
		Code:    10090,
		Message: "api happen error",
	}
	RPermissionDenied = &Response{
		Code:    10091,
		Message: "permission denied",
	}
	RMissingApiParams = &Response{
		Code:    10092,
		Message: "missing api info",
	}
	RErrApiName = &Response{
		Code:    10093,
		Message: "error api name",
	}
)

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

/* ----------- API DEFINE ------------- */

// 处理器
type Processor interface {
	Request(fn string, ctx Context) *Response
}

// API服务
type Server interface {
	// 注册客户端
	Register(name string, p Processor)
	// adds middleware
	Use(middleware ...MiddlewareFunc)
	// adds after middleware
	After(middleware ...MiddlewareFunc)
	// serve http
	ServeHTTP(w http.ResponseWriter, h *http.Request)
}

// 上下文
type Context interface {
	// 返回接口KEY
	Key() string
	// 返回对应用户编号
	User() int64
	Registry() *gof.Registry
	Request() *http.Request
	Form() Form
}

// 上下文工厂
type ContextFactory interface {
	Factory(h *http.Request, key string, userId int64) Context
}

// 中间件
type MiddlewareFunc func(ctx Context) error

// 交换信息，根据key返回用户编号、密钥和是否验证签名
type SwapFunc func(key string) (userId int64, secret string, checkSign bool)

// 数据
type Form map[string]interface{}

func (f Form) Get(key string) interface{} {
	if v, ok := f[key]; ok {
		return v
	}
	return ""
}
func (f Form) Set(key string, value interface{}) {
	f[key] = value
}

var _ Server = new(ServeMux)

// default server implement
type ServeMux struct {
	processors      map[string]Processor
	mux             sync.Mutex
	swap            SwapFunc
	factory         ContextFactory
	middleware      []MiddlewareFunc
	afterMiddleware []MiddlewareFunc
}

func NewServerMux(cf ContextFactory, swap SwapFunc) Server {
	return &ServeMux{
		swap:            swap,
		factory:         cf,
		processors:      map[string]Processor{},
		middleware:      []MiddlewareFunc{},
		afterMiddleware: []MiddlewareFunc{},
	}
}

// 注册客户端
func (a *ServeMux) Register(name string, h Processor) {
	a.mux.Lock()
	defer a.mux.Unlock()
	ls := strings.ToLower(name)
	_, b := a.processors[ls]
	if b {
		panic(errors.New("processor " + name + " has been resisted!"))
	}
	a.processors[ls] = h
}

// adds middleware
func (a *ServeMux) Use(middleware ...MiddlewareFunc) {
	a.middleware = append(a.middleware, middleware...)
}

// adds after middleware
func (a *ServeMux) After(middleware ...MiddlewareFunc) {
	a.afterMiddleware = append(a.afterMiddleware, middleware...)
}

func (a *ServeMux) ServeHTTP(w http.ResponseWriter, h *http.Request) {
	h.ParseForm()
	rsp := a.serveFunc(h)
	var data []byte
	if len(rsp) > 1 {
		data, _ = json.Marshal(rsp)
	} else {
		data, _ = json.Marshal(rsp[0])
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// 处理请求,如果同时请求多个api,那么api参数用","隔开
func (a *ServeMux) serveFunc(h *http.Request) []*Response {
	rsp, userId := a.checkApiPerm(h.Form)
	if rsp != nil {
		return []*Response{rsp}
	}
	key := h.Form.Get("key")
	name := strings.Split(h.Form.Get("api"), ",")
	arr := make([]*Response, len(name))
	// create api context
	ctx := a.factory.Factory(h, key, userId)
	// copy form data
	for i, v := range h.Form {
		ctx.Form().Set(i, v[0])
	}
	// call api
	for i, n := range name {
		arr[i] = a.call(n, ctx)
	}
	return arr
}

// call api
func (a *ServeMux) call(apiName string, ctx Context) *Response {
	data := strings.Split(apiName, ".")
	if len(data) != 2 {
		return RErrApiName
	}
	// save api name
	ctx.Form().Set("api_name", apiName)
	// process api
	entry, fn := strings.ToLower(data[0]), data[1]
	if p, ok := a.processors[entry]; ok {
		// use middleware
		for _, m := range a.middleware {
			if err := m(ctx); err != nil {
				return a.response(apiName, ctx, &Response{
					Code:    RError.Code,
					Message: err.Error(),
				})
			}
		}
		return a.response(apiName, ctx, p.Request(fn, ctx))
	}
	return RErrApiName
}

// use response middleware
func (a *ServeMux) response(apiName string, ctx Context, rsp *Response) *Response {
	if len(a.afterMiddleware) > 0 {
		ctx.Form().Set("api_response", rsp)
		for _, m := range a.afterMiddleware {
			m(ctx)
		}
	}
	return rsp
}

// 检查接口权限
func (a *ServeMux) checkApiPerm(r url.Values) (rsp *Response, userId int64) {
	key := r.Get("key")
	sign := r.Get("sign")
	signType := r.Get("sign_type")
	// 检查参数
	if key == "" || sign == "" || signType == "" {
		return RMissingApiParams, 0
	}
	if signType != "md5" && signType != "sha1" {
		return RMissingApiParams, 0
	}
	userId, userToken, checkSign := a.swap(key)
	if userId <= 0 {
		return RPermissionDenied, userId
	}
	// 检查签名
	if checkSign {
		if s := Sign(signType, r, userToken); s != sign {
			//log.Println("---", token, sign,s)
			return RPermissionDenied, userId
		}
	}
	return nil, userId
}

// 默认上下文工厂
var _ ContextFactory = new(defaultContextFactory)
var DefaultContextFactory *defaultContextFactory = &defaultContextFactory{}

var _ Context = new(defaultContext)

type defaultContext struct {
	h        *http.Request
	key      string
	userId   int64
	registry *gof.Registry
	form     Form
}

func (ctx *defaultContext) Key() string {
	return ctx.key
}

func (ctx *defaultContext) User() int64 {
	return ctx.userId
}

func (ctx *defaultContext) Registry() *gof.Registry {
	return ctx.registry
}

func (ctx *defaultContext) Request() *http.Request {
	return ctx.h
}

func (ctx *defaultContext) Form() Form {
	return ctx.form
}

type defaultContextFactory struct {
	Registry *gof.Registry
}

func (d *defaultContextFactory) Factory(h *http.Request, key string, userId int64) Context {
	return &defaultContext{
		h:        h,
		key:      key,
		userId:   userId,
		registry: d.Registry,
		form:     map[string]interface{}{},
	}
}
