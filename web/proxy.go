/**
 * Copyright 2015 @ S1N1 Team.
 * name : proxy.go
 * author : newmin
 * date : -- :
 * description :
 * history :
 */
package web

import (
	"net/http"
)

type ResponseProxyWriter struct {
	writer http.ResponseWriter
	Output []byte
}

func (this *ResponseProxyWriter) Header() http.Header {
	return this.writer.Header()
}
func (this *ResponseProxyWriter) Write(bytes []byte) (int, error) {
	this.Output = append(this.Output, bytes[0:len(bytes)]...)
	return this.writer.Write(bytes)
}
func (this *ResponseProxyWriter) WriteHeader(i int) {
	this.writer.WriteHeader(i)
}

//创建一个新的HttpWriter
func NewRespProxyWriter(w http.ResponseWriter) *ResponseProxyWriter {
	return &ResponseProxyWriter{
		writer: w,
		Output: []byte{},
	}
}
