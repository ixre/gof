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
	http2 "github.com/jsix/gof/net/http"
	"hash"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// 接口响应
type Response struct {
	// 错误码
	RspCode int64
	// 错误消息
	ErrMsg string
	// 数据结果
	Data interface{}
}

func NewResponse(data interface{}) *Response {
	return &Response{
		Data: data,
	}
}

func NewErrorResponse(message string) *Response {
	return &Response{
		RspCode: RErrorCode,
		ErrMsg:  message,
		Data:    nil,
	}
}

var (
	// 成功码
	RSuccessCode int64 = 0
	// 错误码
	RErrorCode int64 = 1
	// 错误响应
	RInternalError = &Response{
		RspCode: 10090,
		ErrMsg:  "server internal error",
	}
	// 无权限调用
	RAccessDenied = &Response{
		RspCode: 10091,
		ErrMsg:  "api access denied",
	}
	// 接口未定义
	RUndefinedApi = &Response{
		RspCode: 10092,
		ErrMsg:  "api not defined",
	}
	// 接口参数有误
	RIncorrectApiParams = &Response{
		RspCode: 10093,
		ErrMsg:  "incorrect api parameters",
	}
	// 接口已过期
	RDeprecated = &Response{
		RspCode: 10094,
		ErrMsg:  "api is deprecated",
	}
)

// 参数首字母小写后排序，排除sign和sign_type，secret，转换为字节
func ParamsToBytes(r url.Values, secret string) []byte {
	keys := keyArr{}
	for k := range r {
		keys = append(keys, k)
	}
	sort.Sort(keys)
	// 拼接参数和值
	i := 0
	buf := bytes.NewBuffer(nil)
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
	buf.WriteString(secret)
	return buf.Bytes()
}

// 签名
func Sign(signType string, r url.Values, secret string) string {
	data := ParamsToBytes(r, secret)
	switch signType {
	case "md5":
		return byteHash(md5.New(), data)
	case "sha1":
		return byteHash(sha1.New(), data)
	}
	return ""
}

// 计算Hash值
func byteHash(h hash.Hash, data []byte) string {
	h.Write(data)
	b := h.Sum(nil)
	return hex.EncodeToString(b)
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
	Build(registry map[string]interface{}) ContextFactory
}

// 中间件
type MiddlewareFunc func(ctx Context) error

// 交换信息，根据key返回用户编号、密钥
type SwapFunc func(key string) (userId int64, secret string)

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

func NewServerMux(cf ContextFactory, swap SwapFunc) *ServeMux {
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
	s.flushOutputWriter(w, rsp)
}

// 将响应输出
func (s *ServeMux) flushOutputWriter(w http.ResponseWriter, rsp []*Response) {
	if rsp == nil || len(rsp) == 0 {
		panic("no such response can flush to writer")
	}
	for _, r := range rsp {
		if r.RspCode > RSuccessCode {
			buf := bytes.NewBuffer(nil)
			buf.WriteString("#")
			buf.WriteString(strconv.Itoa(int(r.RspCode)))
			buf.WriteString("#")
			buf.WriteString(r.ErrMsg)
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusForbidden)
			w.Write(buf.Bytes())
			return
		}
	}
	var data []byte
	if len(rsp) > 1 {
		var arr []interface{}
		for _, v := range rsp {
			arr = append(arr, v.Data)
		}
		data, _ = json.Marshal(arr)
	} else {
		if rsp[0].Data != nil {
			data, _ = json.Marshal(rsp[0].Data)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// 处理请求,如果同时请求多个api,那么api参数用","隔开
func (s *ServeMux) serveFunc(h *http.Request) []*Response {
	rsp, userId := s.checkAccessPerm(h.Form, h)
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
		return RUndefinedApi
	}
	// save api name
	ctx.Form().Set("$api_name", apiName) // 保存接口名称
	// process api
	entry, fn := strings.ToLower(data[0]), data[1]
	if p, ok := s.processors[entry]; ok {
		// use middleware
		for _, m := range s.middleware {
			if err := m(ctx); err != nil {
				return s.response(apiName, ctx, &Response{
					RspCode: RInternalError.RspCode,
					ErrMsg:  err.Error(),
				})
			}
		}
		return s.response(apiName, ctx, p.Request(fn, ctx))
	}
	return RUndefinedApi
}

// use response middleware
func (s *ServeMux) response(apiName string, ctx Context, rsp *Response) *Response {
	if len(s.afterMiddleware) > 0 {
		ctx.Form().Set("$api_response", rsp) // 保存响应
		for _, m := range s.afterMiddleware {
			m(ctx)
		}
	}
	return rsp
}

// 检查接口权限
func (s *ServeMux) checkAccessPerm(form url.Values, r *http.Request) (rsp *Response, userId int64) {
	key := form.Get("key")
	sign := form.Get("sign")
	signType := form.Get("sign_type")
	// 检查参数
	if key == "" || sign == "" || signType == "" {
		return RIncorrectApiParams, 0
	}
	if signType != "md5" && signType != "sha1" {
		return RIncorrectApiParams, 0
	}
	userId, userSecret := s.swap(key)
	if userId <= 0 {
		return RAccessDenied, userId
	}
	// 检查签名
	if rs := Sign(signType, form, userSecret); rs != sign {
		if !s.trace {
			return RAccessDenied, userId
		}
		ctx := s.factory.Factory(r, key, userId)
		// copy form data
		cf := ctx.Form()
		for i, v := range form {
			cf.Set(i, v[0])
		}
		// set variables
		cf.Set("$user_id", userId)
		cf.Set("$user_secret", userSecret)
		cf.Set("$client_sign", sign)
		cf.Set("$server_sign", rs)
		return s.response(form.Get("api"), ctx, RAccessDenied), userId
	}
	return nil, userId
}

var _ Context = new(defaultContext)

type defaultContext struct {
	h      *http.Request
	key    string
	userId int64
	form   Form
}

func (ctx *defaultContext) Key() string {
	return ctx.key
}

func (ctx *defaultContext) User() int64 {
	return ctx.userId
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
	registry map[string]interface{}
	trace    bool
}

func (d *defaultContextFactory) Build(registry map[string]interface{}) ContextFactory {
	return &defaultContextFactory{
		registry: registry,
	}
}

func (d *defaultContextFactory) setTrace(trace bool) {
	d.trace = trace
}

func (d *defaultContextFactory) Factory(h *http.Request, key string, userId int64) Context {
	ctx := &defaultContext{
		h:      h,
		key:    key,
		userId: userId,
		form:   map[string]interface{}{},
	}
	if d.registry != nil {
		for k, v := range d.registry {
			ctx.form[k] = v
		}
	}
	if d.trace {
		if h != nil {
			ctx.form.Set("$user_ip", http2.RealIp(h))
			ctx.form.Set("$user_agent", h.UserAgent())
		}
	}
	return ctx
}

// 数据
type Form map[string]interface{}

func (f Form) Contains(key string) bool {
	_, ok := f[key]
	return ok
}

// 获取数值
func (f Form) GetInt(key string) int {
	o := f.Get(key)
	switch o.(type) {
	case int:
		return o.(int)
	case int32:
		return int(o.(int32))
	case int64:
		return int(o.(int64))
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

/*------ other support code ------*/
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
