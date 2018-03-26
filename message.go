package gof

import (
	"encoding/json"
	"strings"
)

//操作Json结果
type Message struct {
	// 错误码
	ErrCode int `json:"errCode"`
	// 错误信息
	ErrMsg string `json:"errMsg"`
	//todo:删除，用ErrCode代替
	Result bool `json:"result"`
	//Code    int         `json:"code"`
	Data interface{} `json:"data"`
}

func (m *Message) Error(err error) *Message {
	if err == nil {
		return m.ErrorText("")
	}
	return m.ErrorText(err.Error())
}

func (m *Message) ErrorText(err string) *Message {
	if err = strings.TrimSpace(err); err != "" {
		m.ErrCode = 1
		m.Result = false
		m.ErrMsg = err
	} else {
		m.Result = true
	}
	return m
}

//序列化
func (m Message) Marshal() []byte {
	d, _ := json.Marshal(m)
	return d
}
