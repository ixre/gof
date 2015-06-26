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
	"io"
	"net/http"
)

// Template
type Template struct {
	Init func(*TemplateDataMap)
}

// the data map for template
type TemplateDataMap map[string]interface{}

//type FuncMap template.FuncMap

func (this TemplateDataMap) Add(key string, v interface{}) {
	this[key] = v
}

func (this TemplateDataMap) Del(key string) {
	delete(this, key)
}

// execute template
func (this *Template) ExecuteWithFunc(w io.Writer, funcMap template.FuncMap, dataMap TemplateDataMap,
	tplPath ...string) error {

	t := template.New("-")

	if funcMap != nil {
		t = t.Funcs(funcMap)
	}

	t, err := t.ParseFiles(tplPath...)
	if err != nil {
		return this.handleError(w, err)
	}

	if this.Init != nil {
		if dataMap == nil {
			dataMap = TemplateDataMap{}
		}
		this.Init(&dataMap)
	}

	err = t.Execute(w, dataMap)

	return this.handleError(w, err)
}

func (this *Template) Execute(w io.Writer, dataMap TemplateDataMap, tplPath ...string) error {
	t, err := template.ParseFiles(tplPath...)
	if err != nil {
		return this.handleError(w, err)
	}

	if this.Init != nil {
		if dataMap == nil {
			dataMap = TemplateDataMap{}
		}
		this.Init(&dataMap)
	}

	err = t.Execute(w, dataMap)

	return this.handleError(w, err)
}

func (this *Template) handleError(w io.Writer, err error) error {
	if err != nil {
		if rsp, ok := w.(http.ResponseWriter); ok {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
		}
	}
	return err
}
