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

func (this *Message) Error(err error) *Message {
	if err != nil {
		this.Result = false
		this.Message = err.Error()
	} else {
		this.Result = true
	}
	return this
}

//序列化
func (this Message) Marshal() []byte {
	json, _ := json.Marshal(this)
	return json
}
