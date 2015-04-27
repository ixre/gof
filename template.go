/**
 * Copyright 2015 @ S1N1 Team.
 * name : template.go
 * author : newmin
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

type Template struct {
	Init func(m *map[string]interface{})
}

// the data map for template
type TemplateMapFunc func(m *map[string]interface{})

// execute single template file
func (this *Template) Render(w io.Writer, tplPath string, f TemplateMapFunc,
) error {
	return this.Execute(w, f, tplPath)
}

// execute template
func (this *Template) Execute(w io.Writer, f TemplateMapFunc,
	tplPath ...string) error {

	t, err := template.ParseFiles(tplPath...)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	if this.Init != nil {
		this.Init(&data)
	}
	if f != nil {
		f(&data)
	}

	return t.Execute(w, &data)
}

// execute template,when happen error return a http error.
func (this *Template) ExecuteIncludeErr(w http.ResponseWriter, f TemplateMapFunc, tplPath ...string) {
	err := this.Execute(w, f, tplPath...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
