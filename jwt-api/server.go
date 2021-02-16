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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	http2 "github.com/ixre/gof/net/http"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	RCInternalError = 1001 // server internal error
	RCAccessDenied  = 1002 // access denied
	RCUndefinedApi  = 1003 // api not defined
	RCNotAuthorized = 1004 // not authorized
	RCInvalidToken  = 1005 // access token is invalid
	RCTokenExpires  = 1006 // access token is expired
	RCDeprecated    = 1009 // api is deprecated
)

// 接口响应
type Response struct {
	// 响应码
	Code int
	// 响应消息
	Message string
	// 响应数据
	Data interface{} `json:"Data,omitempty"`
}

func NewResponse(data interface{}) *Response {
	return &Response{
		Code:    0,
		Message: "",
		Data:    data,
	}
}

func ResponseWithCode(code int, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}
func parseErrResp(errCode int) *Response {
	var message = ""
	switch errCode {
	case RCInternalError:
		message = "server internal error"
	case RCAccessDenied:
		message = "access denied"
	case RCUndefinedApi:
		message = "api not defined"
	case RCNotAuthorized:
		message = "not authorized"
	case RCInvalidToken:
		message = "access token is invalid"
	case RCTokenExpires:
		message = "access token is expired"
	case RCDeprecated:
		message = "api is deprecated"
	}
	return ResponseWithCode(errCode, message)
}

/* ----------- API DEFINE ------------- */

// 接口处理器
type Handler interface {
	// API Group
	Group() string
	Process(fn string, ctx Context) *Response
}

// 接口处理方法
type HandlerFunc func(ctx Context) interface{}

// 中间件
type MiddlewareFunc func(ctx Context) error

// 获取用户私钥,返回错误后将直接输出到客户端
type SwapPrivateKeyFunc func(ctx Context) (privateKey []byte, err error)

type Claims = jwt.Claims
type MapClaims = jwt.MapClaims

// 创建凭据
func CreateClaims(aud string, iss string, sub string, expires int64) Claims {
	return jwt.MapClaims{
		"aud": aud,
		"exp": expires,
		"iss": iss,
		"sub": sub,
	}
}

// 生成访问token
func AccessToken(claims Claims, privateKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(privateKey)
}

// api server
type Server interface {
	// register api handler
	Handle(p Handler)
	// public api not require authorized
	HandlePublic(h Handler)
	// adds middleware
	Use(middleware ...MiddlewareFunc)
	// adds after middleware
	After(middleware ...MiddlewareFunc)
	// trace mode
	Trace()
	// serve http
	ServeHTTP(w http.ResponseWriter, h *http.Request)
}

type RequestWrapper struct {
	Request *http.Request
	// receive posted query and form params
	Params     StoredValues
	UserAddr   string
	UserAgent  string
	RequestApi string
}

func (w RequestWrapper) Bind(e interface{}) error {
	buf, err := ioutil.ReadAll(w.Request.Body)
	//w.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	if err == nil {
		return json.Unmarshal(buf, &e)
	}
	return nil
}

type ResponseWrapper struct {
	Response *Response
	http.ResponseWriter
}

// 上下文
type Context interface {
	// Api user key
	UserKey() string
	// claims
	Claims() Claims
	// 请求
	Request() *RequestWrapper
	// 响应
	Response() *ResponseWrapper
}

// 处理接口方法
func HandleMultiFunc(fn string, ctx Context, funcMap map[string]HandlerFunc) *Response {
	if v, ok := funcMap[fn]; ok {
		d := v(ctx)
		switch d.(type) {
		case *Response:
			return d.(*Response)
		case Response:
			r := d.(Response)
			return &r
		}
		return &Response{Data: d}
	}
	return parseErrResp(RCUndefinedApi)
}

var _ Server = new(ServeMux)

type handlerWrapper struct {
	pub     bool
	handler Handler
}

// default server implement
type ServeMux struct {
	trace           bool
	cors            bool
	prefix          string
	processors      map[string]*handlerWrapper
	mux             sync.Mutex
	swapKey         SwapPrivateKeyFunc
	middleware      []MiddlewareFunc
	afterMiddleware []MiddlewareFunc
}

func NewServerMux(swap SwapPrivateKeyFunc, prefix string, cors bool) Server {
	return &ServeMux{
		cors:            cors,
		prefix:          prefix,
		swapKey:         swap,
		processors:      map[string]*handlerWrapper{},
		middleware:      []MiddlewareFunc{},
		afterMiddleware: []MiddlewareFunc{},
	}
}

// 注册客户端
func (s *ServeMux) handle(h Handler, pub bool) {
	s.mux.Lock()
	defer s.mux.Unlock()
	ls := strings.ToLower(h.Group())
	_, b := s.processors[ls]
	if b {
		panic(errors.New("processor " + h.Group() + " has been resisted!"))
	}
	s.processors[ls] = &handlerWrapper{pub: pub, handler: h}
}

func (s *ServeMux) Handle(h Handler) {
	s.handle(h, false)
}

func (s *ServeMux) HandlePublic(h Handler) {
	s.handle(h, true)
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
	if rsp := s.serve(w, h); rsp != nil {
		s.flushOutputWriter(w, rsp)
	}
}

func (s *ServeMux) serve(w http.ResponseWriter, h *http.Request) *Response {
	if s.cors {
		var origin = h.Header.Get("ORIGIN")
		if h.Method == "OPTIONS" {
			s.preFlight(w, origin)
			return nil
		}
		w.Header().Add("Access-Control-Allow-Origin", origin)
	}
	ctx := s.factoryContext(h, w)
	// use middleware
	for _, m := range s.middleware {
		if err := m(ctx); err != nil {
			return ResponseWithCode(RCAccessDenied, err.Error())
		}
	}
	entry, fn := s.getEntry(ctx)
	// 获取处理器
	proc, ok := s.processors[entry]

	// require authorized
	if proc == nil || !proc.pub {
		// check headers
		accessToken := h.Header.Get("Authorization")
		if len(accessToken) == 0 {
			return parseErrResp(RCAccessDenied)
		}
		// swapKey private key
		privateKey, err := s.swapKey(ctx)
		if err != nil {
			return ResponseWithCode(RCNotAuthorized, err.Error())
		}
		// valid jwt token
		claims, code := s.jwtVerify(accessToken, privateKey)
		if code > 0 {
			return parseErrResp(code)
		}
		ctx.SetClaims(claims)
	}
	if !ok {
		return parseErrResp(RCUndefinedApi)
	}
	return s.serveFunc(ctx, proc.handler, fn)
}

// 将响应输出
func (s *ServeMux) flushOutputWriter(w http.ResponseWriter, rsp *Response) {
	if rsp == nil {
		panic("no such response can flush to writer")
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if rsp.Code > RCInternalError && RCDeprecated > rsp.Code {
		ret := fmt.Sprintf("{\"err_code\":%d,\"err_msg\":\"%s\"}", rsp.Code, rsp.Message)
		_, _ = w.Write([]byte(ret))
		return
	}
	//note: 直接返回包含状态码的接口
	var data []byte
	data = s.marshal(rsp)
	_, _ = w.Write(data)
	return

	// 如果包含数据, 直接返回数据, 否则返回Response
	if rsp.Data != nil {
		switch rsp.Data.(type) {
		case string:
			data = []byte(rsp.Data.(string))
		case int:
			data = []byte(strconv.Itoa(rsp.Data.(int)))
		case bool:
			if rsp.Data.(bool) {
				data = []byte("true")
			} else {
				data = []byte("false")
			}
		default:
			data = s.marshal(rsp.Data)
		}
	} else {
		data = s.marshal(rsp)
	}
	_, _ = w.Write(data)
}

// 处理请求,如果同时请求多个api,那么api参数用","隔开
func (s *ServeMux) serveFunc(ctx Context, entry Handler, fn string) *Response {
	// call api
	rsp := entry.Process(fn, ctx)
	if rsp == nil {
		return parseErrResp(RCUndefinedApi)
	}
	return s.responseMiddleware(ctx, rsp)
}

func (s *ServeMux) Trace() {
	s.trace = true
}

// use response middleware
func (s *ServeMux) responseMiddleware(ctx Context, rsp *Response) *Response {
	if len(s.afterMiddleware) > 0 {
		ctx.Response().Response = rsp // 保存响应
		for _, m := range s.afterMiddleware {
			_ = m(ctx)
		}
	}
	return rsp
}

func (s *ServeMux) marshal(d interface{}) []byte {
	b, _ := json.Marshal(d)
	return b
}

func (s *ServeMux) preFlight(w http.ResponseWriter, origin string) {
	header := w.Header()
	header.Add("Access-Control-Allow-Origin", origin)
	header.Add("Access-Control-Allow-Methods", "PUT, GET, POST, DELETE, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type,Credentials, Accept, Authorization, Access-Control-Allow-Credentials")
	header.Add("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(""))
}

// valid jwt token, if not right return error responseMiddleware
func (s *ServeMux) jwtVerify(token string, privateKey []byte) (Claims, int) {
	// 转换token
	dstClaims := jwt.MapClaims{} // 可以用实现了Claim接口的自定义结构
	tk, err := jwt.ParseWithClaims(token, &dstClaims, func(t *jwt.Token) (interface{}, error) {
		return privateKey, nil
	})
	if tk == nil {
		return nil, RCNotAuthorized
	}
	// 判断是否有效
	if !tk.Valid {
		ve, _ := err.(*jwt.ValidationError)
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			return nil, RCTokenExpires
		} else {
			//println("--", ve.Errors)
			return nil, RCInvalidToken
		}
	}
	return dstClaims, 0
}

func (s *ServeMux) getEntry(ctx Context) (entry, action string) {
	if v := ctx.Request().Params["$api"]; v != nil {
		s := strings.Replace(v.(string), "/", ".", -1)
		a := strings.Split(s, ".")
		if len(a) == 0 {
			panic("api should named like 'user.login' or 'user/login'")
		}
		return a[0], a[1]
	}
	path := ctx.Request().Request.URL.Path
	if s.prefix != "" {
		path = path[len(s.prefix):]
	}
	if path[0] == '/' {
		path = path[1:]
	}
	arr := strings.Split(path, "/")
	entry = arr[0]
	if len(arr) >= 2 {
		action = arr[1]
	}
	// save api name
	if action == "" {
		ctx.Request().RequestApi = entry
	} else {
		ctx.Request().RequestApi = entry + "." + action
	}
	return entry, action
}

func (s *ServeMux) factoryContext(h *http.Request, w http.ResponseWriter) *defaultContext {
	_ = h.ParseForm()
	userKey := h.Header.Get("user-key")
	if userKey == "" {
		userKey = h.FormValue("user-key")
	}
	return createContext(h, w, userKey)
}

var _ Context = new(defaultContext)

type defaultContext struct {
	r       *RequestWrapper
	w       *ResponseWrapper
	userKey string
	claims  Claims
}

func (c *defaultContext) Request() *RequestWrapper {
	return c.r
}

func createContext(h *http.Request, w http.ResponseWriter, userKey string) *defaultContext {
	r := &RequestWrapper{
		Request: h,
		Params:  map[string]interface{}{},
	}
	ctx := &defaultContext{
		w:       &ResponseWrapper{ResponseWriter: w},
		userKey: userKey,
		r:       r,
	}
	if h != nil {
		r.UserAddr = http2.RealIp(h)
		r.UserAgent = h.UserAgent()
		// parseForm query params
		for i, v := range h.URL.Query() {
			r.Params.Set(i, v[0])
		}
		// parseForm form data
		for i, v := range h.Form {
			r.Params.Set(i, v[0])
		}
	}
	return ctx
}

func (c *defaultContext) UserKey() string {
	return c.userKey
}

func (c *defaultContext) SetClaims(claims Claims) {
	c.claims = claims
}

func (c *defaultContext) Claims() Claims {
	return c.claims
}

func (c *defaultContext) Response() *ResponseWrapper {
	return c.w
}

// 数据
type StoredValues map[string]interface{}

func (f StoredValues) Contains(key string) bool {
	_, ok := f[key]
	return ok
}

// 获取数值
func (f StoredValues) GetInt(key string) int {
	o := f.Get(key)
	if o == nil {
		return 0
	}
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
func (f StoredValues) GetBytes(key string) []byte {
	if v, ok := f.Get(key).(string); ok {
		return []byte(v)
	}
	return []byte(nil)
}

// 获取字符串
func (f StoredValues) GetString(key string) string {
	if v, ok := f.Get(key).(string); ok {
		return v
	}
	return ""
}

func (f StoredValues) Get(key string) interface{} {
	if v, ok := f[key]; ok {
		return v
	}
	return nil
}
func (f StoredValues) Set(key string, value interface{}) {
	f[key] = value
}
