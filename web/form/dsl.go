package form

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ixre/gof/db/orm"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	tpl "text/template"
)

type (
	Form struct {
		ID     string
		Label  string
		Fields []*Field
	}
	Field struct {
		ID        string
		Label     string
		Elem      string
		Class     string
		Attrs     map[string]string
		AttrsHtml string
	}
	Engine struct {
	}
)

var (
	attrRegex  = regexp.MustCompile("\\s*([^#\\s]+)\\s*=\\s*\"*([^#\\s\"]*)\"*\\s*")
	fieldRegex = regexp.MustCompile("field\\{([\\s\\S]+?)\\}")
	// 默认表单HTML模板
	TDefaultFormHtml string = `<!-- FORM ID:{{.ID}} Name:{{.Label}} -->
            {{range $i,$f := .Fields}}
            <div class="form-row">
                <label>{{$f.Label}}</label>
                <div class="f">
                    <{{$f.Elem}} field="{{$f.ID}}" class="{{$f.Class}}"{{$f.AttrsHtml}}></{{$f.Elem}}>
                </div>
            </div>
            {{end}}`
)

func init() {
	TDefaultFormHtml = strings.Replace(
		TDefaultFormHtml, "            ", "", -1)
}

// 将DSL转为表单对象
func (e *Engine) Parse(dsl string) (*Form, error) {
	p1i := strings.Index(dsl, "field")
	if p1i == -1 {
		return nil, errors.New("表单未包含任何域")
	}
	f := &Form{}
	p1Str := dsl[:p1i]
	p2Str := dsl[p1i:]
	mcs := attrRegex.FindAllStringSubmatch(p1Str, -1)
	for _, v := range mcs {
		switch v[1] {
		case "id":
			f.ID = v[2]
		case "label":
			f.Label = v[2]
		}
	}
	fMcs := fieldRegex.FindAllString(p2Str, -1)
	f.Fields = make([]*Field, len(fMcs))
	for i, mc := range fMcs {
		fd := &Field{Attrs: make(map[string]string)}
		mcs = attrRegex.FindAllStringSubmatch(mc, -1)
		for _, v := range mcs {
			switch v[1] {
			case "id":
				fd.ID = v[2]
			case "label":
				fd.Label = v[2]
			case "element":
				fd.Elem = v[2]
			case "class":
				fd.Class = v[2]
			default:
				fd.Attrs[v[1]] = v[2]
			}
		}
		f.Fields[i] = fd
	}
	return f, nil
}

// 将DSL文件转为表单对象
func (e *Engine) ParseFile(dslFile string) (*Form, error) {
	data, err := ioutil.ReadFile(dslFile)
	if err == nil {
		return e.Parse(string(data))
	}
	return nil, err
}

func (e *Engine) SaveDSL(f *Form, path string) error {
	fi, err := os.Open(path)
	if os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
		fi, err = os.Create(path)
	}
	if err == nil {
		defer fi.Close()
		fi.WriteString("id=" + f.ID)
		fi.WriteString("\nlabel=" + f.Label)
		for _, f := range f.Fields {
			fi.WriteString("\nfield{\n")
			fi.WriteString(fmt.Sprintf("    id=%s\n", f.ID))
			fi.WriteString(fmt.Sprintf("    label=%s\n", f.Label))
			fi.WriteString(fmt.Sprintf("    element=%s\n", f.Elem))
			fi.WriteString(fmt.Sprintf("    class=%s\n", f.Class))
			if f.Attrs != nil {
				for k, a := range f.Attrs {
					fi.WriteString(fmt.Sprintf("    %s=%s\n", k, a))
				}
			}
			fi.WriteString("}")
		}
	}
	return err
}

func (t *Engine) title(str string) string {
	arr := strings.Split(str, "_")
	for i, v := range arr {
		arr[i] = strings.Title(v)
	}
	return strings.Join(arr, "")
}

// 将数据库表格转为表单对象
func (e *Engine) TableToForm(tb *orm.Table) *Form {
	f := &Form{
		ID:     e.title(tb.Name),
		Label:  e.title(tb.Comment),
		Fields: []*Field{},
	}
	for _, v := range tb.Columns {
		if !v.Auto {
			fd := &Field{
				ID:    e.title(v.Name),
				Label: e.title(v.Comment),
				Elem:  "input",
			}
			if v.IsPk {
				fd.Attrs = map[string]string{
					"type":    "hidden",
					"primary": "true",
				}
			}
			f.Fields = append(f.Fields, fd)
		}
	}
	return f
}

func (e *Engine) htmlPrepare(f *Form) {
	for _, v := range f.Fields {
		if v.AttrsHtml != "" || v.Attrs == nil || len(v.Attrs) == 0 {
			continue
		}
		for k, v2 := range v.Attrs {
			v.AttrsHtml += fmt.Sprintf(" %s=\"%s\"", k, v2)
		}
	}
}

// 为表单对象生成HTML
func (e *Engine) HtmlForm(f *Form, template string) (string, error) {
	e.htmlPrepare(f)
	t, err := (&tpl.Template{}).Parse(template)
	if err == nil {
		buf := bytes.NewBuffer(nil)
		err = t.Execute(buf, f)
		return buf.String(), err
	}
	return "", err
}

func (e *Engine) HtmlDefaultForm(f *Form) (string, error) {
	return e.HtmlForm(f, TDefaultFormHtml)
}

// 保存表单到文件中
func (e *Engine) SaveHtmlForm(f *Form, template string,
	path string) (htm string, err error) {
	htm, err = e.HtmlForm(f, template)
	if err == nil {
		if _, err = os.Stat(path); os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(path), os.ModePerm)
		}
		err = ioutil.WriteFile(path, []byte(htm), os.ModePerm)
	}
	return htm, err
}
