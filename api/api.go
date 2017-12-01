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
	"strconv"
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
	RErrUndefinedApi = &Response{
		Code:    10093,
		Message: "api not defined",
	}
)

// 参数首字母小写后排序，排除sign和sign_type，拼接token，转换为字节
func paramsToBytes(r url.Values, token string) []byte {
	keys := keyArr{}
	for k, _ := range r {
		keys = append(keys, k)
	}
	sort.Sort(keys)
	// 拼接参数和值
	buf := bytes.NewBuffer(nil)
	for i, k := range keys {
		if k == "sign" || k == "sign_type" {
			continue
		}
		if i > 0 {
			buf.WriteString("&")
		}
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(r[k][0])
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

var _ sort.Interface = keyArr{}

type keyArr []string

func (s keyArr) Len() int {
	return len(s)
}

func (s keyArr) Less(i, j int) bool {
	return strings.ToLower(s[i]) < strings.ToLower(s[j])
}

func (s keyArr) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
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
	// trace mode
	Trace()
	// serve http
	ServeHTTP(w http.ResponseWriter, h *http.Request)
}

// 上下文
type Context interface {
	// 返回接口KEY
	Key() string
	// 返回对应用户编号
	User() int64
	// 注册表
	Registry() *gof.Registry
	// 请求
	Request() *http.Request
	// 表单数据
	Form() Form
}

// 上下文工厂
type ContextFactory interface {
	Factory(h *http.Request, key string, userId int64) Context
}

// 工厂生成器
type FactoryBuilder interface {
	// 生成下文工厂
	Build(registry *gof.Registry) ContextFactory
}

// 中间件
type MiddlewareFunc func(ctx Context) error

// 交换信息，根据key返回用户编号、密钥和是否验证签名
type SwapFunc func(key string) (userId int64, secret string, checkSign bool)

// 数据
type Form map[string]interface{}

// 获取数值
func (f Form) GetInt32(key string) int32 {
	o := f.Get(key)
	switch o.(type) {
	case int, int32, int64:
		return int32(o.(int))
	case string:
		v, _ := strconv.Atoi(o.(string))
		return int32(v)
	}
	panic("not int or string")
}

// 获取数值
func (f Form) GetInt(key string) int {
	o := f.Get(key)
	switch o.(type) {
	case int, int32, int64:
		return o.(int)
	case string:
		v, _ := strconv.Atoi(o.(string))
		return v
	}
	panic("not int or string")
}

// 获取字节
func (f Form) GetBytes(key string) []byte {
	if v, ok := f.Get(key).(string); ok {
		return []byte(v)
	}
	return []byte(nil)
}

// 获取字符串
func (f Form) GetString(key string) string {
	if v, ok := f.Get(key).(string); ok {
		return v
	}
	return ""
}

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
	trace           bool
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
func (s *ServeMux) Register(name string, h Processor) {
	s.mux.Lock()
	defer s.mux.Unlock()
	ls := strings.ToLower(name)
	_, b := s.processors[ls]
	if b {
		panic(errors.New("processor " + name + " has been resisted!"))
	}
	s.processors[ls] = h
}

// adds middleware
func (s *ServeMux) Use(middleware ...MiddlewareFunc) {
	s.middleware = append(s.middleware, middleware...)
}

// adds after middleware
func (s *ServeMux) After(middleware ...MiddlewareFunc) {
	s.afterMiddleware = append(s.afterMiddleware, middleware...)
}

func (s *ServeMux) ServeHTTP(w http.ResponseWriter, h *http.Request) {
	h.ParseForm()
	rsp := s.serveFunc(h)
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
func (s *ServeMux) serveFunc(h *http.Request) []*Response {
	rsp, userId := s.checkApiPerm(h.Form, h)
	if rsp != nil {
		return []*Response{rsp}
	}
	key := h.Form.Get("key")
	name := strings.Split(h.Form.Get("api"), ",")
	arr := make([]*Response, len(name))
	// create api context
	ctx := s.factory.Factory(h, key, userId)
	// copy form data
	for i, v := range h.Form {
		ctx.Form().Set(i, v[0])
	}
	// call api
	for i, n := range name {
		arr[i] = s.call(n, ctx)
	}
	return arr
}
func (s *ServeMux) Trace() {
	s.trace = true
	if df, ok := s.factory.(*defaultContextFactory); ok {
		df.setTrace(s.trace)
	}
}

// call api
func (s *ServeMux) call(apiName string, ctx Context) *Response {
	data := strings.Split(apiName, ".")
	if len(data) != 2 {
		return RErrUndefinedApi
	}
	// save api name
	ctx.Form().Set("$api_name", apiName)
	// process api
	entry, fn := strings.ToLower(data[0]), data[1]
	if p, ok := s.processors[entry]; ok {
		// use middleware
		for _, m := range s.middleware {
			if err := m(ctx); err != nil {
				return s.response(apiName, ctx, &Response{
					Code:    RError.Code,
					Message: err.Error(),
				})
			}
		}
		return s.response(apiName, ctx, p.Request(fn, ctx))
	}
	return RErrUndefinedApi
}

// use response middleware
func (s *ServeMux) response(apiName string, ctx Context, rsp *Response) *Response {
	if len(s.afterMiddleware) > 0 {
		ctx.Form().Set("$api_response", rsp)
		for _, m := range s.afterMiddleware {
			m(ctx)
		}
	}
	return rsp
}

// 检查接口权限
func (s *ServeMux) checkApiPerm(form url.Values, r *http.Request) (rsp *Response, userId int64) {
	key := form.Get("key")
	sign := form.Get("sign")
	signType := form.Get("sign_type")
	// 检查参数
	if key == "" || sign == "" || signType == "" {
		return RMissingApiParams, 0
	}
	if signType != "md5" && signType != "sha1" {
		return RMissingApiParams, 0
	}
	userId, userToken, checkSign := s.swap(key)
	if userId <= 0 {
		return RPermissionDenied, userId
	}
	// 检查签名
	if checkSign {
		if rs := Sign(signType, form, userToken); rs != sign {
			if !s.trace {
				return RPermissionDenied, userId
			}
			ctx := s.factory.Factory(r, key, userId)
			// copy form data
			cf := ctx.Form()
			for i, v := range form {
				cf.Set(i, v[0])
			}
			// set variables
			cf.Set("$client_sign", sign)
			cf.Set("$server_sign", rs)
			return s.response(form.Get("api"), ctx, RPermissionDenied), userId
		}
	}
	return nil, userId
}

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

// 默认工厂
var DefaultFactory FactoryBuilder = &defaultContextFactory{}

var _ ContextFactory = new(defaultContextFactory)
var _ FactoryBuilder = new(defaultContextFactory)

type defaultContextFactory struct {
	Registry *gof.Registry
	trace    bool
}

func (d *defaultContextFactory) Build(registry *gof.Registry) ContextFactory {
	return &defaultContextFactory{
		Registry: registry,
	}
}

func (d *defaultContextFactory) setTrace(trace bool) {
	d.trace = trace
}

func (d *defaultContextFactory) Factory(h *http.Request, key string, userId int64) Context {
	ctx := &defaultContext{
		h:        h,
		key:      key,
		userId:   userId,
		registry: d.Registry,
		form:     map[string]interface{}{},
	}
	if d.trace {
		if h != nil {
			ctx.form.Set("$user_ip", h.RemoteAddr)
			ctx.form.Set("$user_agent", h.UserAgent())
		}
	}
	return ctx
}
