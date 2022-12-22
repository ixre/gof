package gof

import (
	"encoding/json"
	"strings"
)

type Result struct {
	ErrCode int               `json:"ErrCode"`
	ErrMsg  string            `json:"ErrMsg,omitempty"`
	Data    map[string]string `json:"Data,omitempty"`
}

func (r *Result) Error(err error) *Result {
	if err == nil {
		return r.ErrorText("")
	}
	return r.ErrorText(err.Error())
}

func ResultWithCode(code int, message string) *Result {
	return &Result{
		ErrCode: code,
		ErrMsg:  message,
		Data:    nil,
	}
}

func (r *Result) ErrorText(err string) *Result {
	if err = strings.TrimSpace(err); err != "" {
		r.ErrCode = 1
		r.ErrMsg = err
	}
	return r
}

// 序列化
func (r Result) Marshal() []byte {
	d, _ := json.Marshal(r)
	return d
}
