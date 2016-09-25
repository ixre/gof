package gof

import (
	"encoding/json"
)

//操作Json结果
type Message struct {
	Result  bool        `json:"result"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func (m *Message) Error(err error) *Message {
	if err != nil {
		m.Result = false
		m.Message = err.Error()
	} else {
		m.Result = true
	}
	return m
}

//序列化
func (m Message) Marshal() []byte {
	json, _ := json.Marshal(m)
	return json
}
