package gof

import (
	"encoding/json"
	"strings"
)

func ErrorResult(err error) *Result {
	return (&Result{}).Error(err)
}

func SuccessResult(v interface{}) *Result {
	r := &Result{}
	r.Data = v
	return r
}

type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (r *Result) Error(err error) *Result {
	if err == nil {
		return r.ErrorText("")
	}
	return r.ErrorText(err.Error())
}

func ResultWithCode(code int, msg string) *Result {
	return &Result{
		Code:    code,
		Message: msg,
		Data:    nil,
	}
}

func (r *Result) ErrorText(err string) *Result {
	if err = strings.TrimSpace(err); err != "" {
		r.Code = 1
		r.Message = err
	}
	return r
}

// 序列化
func (r Result) Marshal() []byte {
	d, _ := json.Marshal(r)
	return d
}
