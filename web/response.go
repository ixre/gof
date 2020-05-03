/**
 * Copyright 2015 @ to2.net.
 * name : response
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var _ http.ResponseWriter = new(response)

type response struct {
	http.ResponseWriter
}

// 输出JSON
func (this *response) JsonOutput(v interface{}) {
	this.ResponseWriter.Header().Set("Content-DbType", "application/json")
	b, err := json.Marshal(v)
	if err != nil {
		str := fmt.Sprintf(`{"error":"%s"}`,
			strings.Replace(err.Error(), "\"", "\\\"", -1))
		this.ResponseWriter.Write([]byte(str))
	} else {
		this.ResponseWriter.Write(b)
	}
}
