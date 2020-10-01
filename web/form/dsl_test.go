package form

import "testing"

func TestEngine_ParseFile(t *testing.T) {
	e := &Engine{}
	form, err := e.ParseFile("test.form")
	if err != nil {
		t.Error(err)
	}

	t.Log("输出表数据")
	t.Log("ID:", form.ID, "; Label:", form.Label)
	for i, v := range form.Fields {
		t.Log("Field ", i, "- ID:", v.ID, "; Label:",
			v.Label, "; Elem:", v.Elem, "; Class:", v.Class)
		for k, v := range v.Attrs {
			t.Log("   Attrs:", k, "=", v)
		}
	}

	t.Log("将表生成DSL并另存为")
	err = e.SaveDSL(form, "test_gen.form")
	if err == nil {
		t.Log("生成成功")
	} else {
		t.Error(err)
	}

	t.Log("将表生成HTML表单并存储")
	htm, err := e.SaveHtmlForm(form, TDefaultFormHtml, "test_gen.html")
	if err == nil {
		t.Log("生成成功,结果如下：\n" + htm)
	} else {
		t.Error(err)
	}
}
