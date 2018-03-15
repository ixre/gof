package gof

import (
	"encoding/json"
	"strings"
)

//操作Json结果
type Message struct {
	ErrCode int         `json:"errCode"`
	Result  bool        `json:"result"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func (m *Message) Error(err error) *Message {
	if err == nil {
		return m.TextError("")
	}
	return m.Error(err)
}

func (m *Message) TextError(err string) *Message {
	if err = strings.TrimSpace(err); err != "" {
		m.ErrCode = 1
		m.Result = false
		m.Message = err
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
