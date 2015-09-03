/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package mvc

import (
	"fmt"
	"github.com/jsix/gof/web"
	"net/http"
	"reflect"
	"strings"
)

func CustomHandle(c Controller, ctx *web.Context, action string, args ...interface{}) {
	w := ctx.Response
	// 拦截器
	filter, isFilter := c.(Filter)
	if isFilter {
		if !filter.Requesting(ctx) { //如果返回false，终止执行
			return
		}
	}

	t := reflect.ValueOf(c)
	method := t.MethodByName(action)

	if !method.IsValid() {
		errMsg := "No action named \"" + strings.Replace(action, "_post", "", 1) +
			"\" in " + reflect.TypeOf(c).String() + "."
		http.Error(w, errMsg, http.StatusNotFound)
		return
	} else {
		//包含基础的ResponseWriter和http.Request 2个参数
		argsLen := len(args)
		numIn := method.Type().NumIn()

		if argsLen < numIn-1 {
			errMsg := fmt.Sprintf("Can't inject to method,it's possible missing parameter!"+
				"\r\ncontroller: %s , action: %s",
				reflect.TypeOf(c).String(), action)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		} else {
			params := make([]reflect.Value, numIn)
			params[0] = reflect.ValueOf(ctx)
			for i := 1; i < numIn; i++ {
				params[i] = reflect.ValueOf(args[i-1])
			}

			method.Call(params)
		}
	}

	if isFilter {
		filter.RequestEnd(ctx)
	}
}

//控制器处理
//@controller ： 包含多种动作，URL中的文件名自动映射到控制器的函数
//				 注意，是区分大小写的,默认映射到index函数
//				 如果是POST请求将映射到控制器“函数名+_post”的函数执行
// @path    : 指定路径
// @re_post : 是否为post请求额外加上_post来区分Post和Get请求
//
func HandlePath(controller Controller, ctx *web.Context, path string, rePost bool, args ...interface{}) {
	r := ctx.Request
	if len(path) == 0 {
		path = r.URL.Path
	}
	var action string = GetAction(path, "")
	if rePost && r.Method == "POST" {
		action += "_post"
	}

	CustomHandle(controller, ctx, action, args...)
}

func Handle(controller Controller, ctx *web.Context, rePost bool, args ...interface{}) {
	HandlePath(controller, ctx, "", rePost, args...)
}

func GetAction(path string, suffix string) string {

	path = path[1:]

	// 获取Action
	var action string
	arr := strings.Split(path, "/")
	if arrLen := len(arr); arrLen == 1 {
		action = path
	} else if arrLen >= 2 {
		action = arr[1]
	}

	if len(action) == 0 {
		return "Index"
	}

	// 去扩展名
	if di := strings.Index(action, "."); di != -1 {
		// 判断后缀是否相同
		if len(suffix) != 0 && action[di:] != suffix {
			return ""
		} else {
			action = action[:di]
		}
	}

	return strings.Title(action)
}
