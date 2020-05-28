package web

import (
	"html/template"
)

// 获取原始的HTML
func RawHtml(html string) interface{} {
	return template.HTML(html)
}

// 获取原始的SCRIPT
func RawScript(script string) interface{} {
	return template.JS(script)
}
