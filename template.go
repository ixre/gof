/**
 * Copyright 2015 @ S1N1 Team.
 * name : template.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */

package gof

import (
	"html/template"
	"net/http"
	"io"
)

// Template
type Template struct {
	Init func(*TemplateDataMap)
}

// the data map for template
type TemplateDataMap map[string]interface{}
func (this TemplateDataMap) Add(key string ,v interface{}) {
	this[key] = v
}

func (this TemplateDataMap) Del(key string) {
	delete(this, key)
}

// execute template
func (this *Template) Execute(w io.Writer, f TemplateDataMap,
	tplPath ...string) error {

	t, err := template.ParseFiles(tplPath...)
	if err != nil {
		return this.handleError(w,err)
	}

	if this.Init != nil && f != nil {
		this.Init(&f)
	}
	err = t.Execute(w, f)

	return this.handleError(w,err)
}

func (this *Template) handleError(w io.Writer,err error)error {
	if err != nil {
		if rsp, ok := w.(http.ResponseWriter); ok {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
		}
	}
	return err
}