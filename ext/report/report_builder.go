package report

import (
	"bytes"
	"strings"
	"time"
)

func getFileName(p *ExportItem) string {
	dtStr := time.Now().Format("200601020405")
	i := strings.LastIndex(p.PortalKey, "/")
	if i != -1 {
		return p.PortalKey[i+1:] + "-" + dtStr
	}
	return p.PortalKey + "-" + dtStr
}

// 生成WEB导出勾选项及脚本
func buildWebExportCheckOptions(p IDataExportPortal, token string) string {
	portal := p.(*ExportItem)
	buf := bytes.NewBufferString("")
	fileName := getFileName(portal)
	// 输出Javascript支持库
	buf.WriteString(`<script type="text/javascript">
        var $_expo={
            //处理请求的Url地址
            urlHandler:'processExport',
            form:null,
            chkFileName:function(ele){
                var ckFile = document.getElementById("ck-file-txt");
                if(ele.checked){
                    ckFile.setAttribute("_value",ckFile.value);
                    ckFile.value = "";
                    ckFile.style.visibility = "hidden";
                }else{
                    ckFile.value = ckFile.getAttribute("_value")||"";
                    ckFile.style.visibility = "visible";
                }
            },
            chkInit:function(){
               if(this.form != null)return false;
               this.form = document.forms["export_form"];
               this.form.setAttribute("action",this.urlHandler);
               document.getElementById("params").value = this.getParams();
            },
            getParams:function(){
               var regMatch=/(\?|&)params=(.+)&*/i.exec(location.search);
               return regMatch?regMatch[2]:'';
            },
            submit:function(e){
                this.chkInit();
                this.form.submit();
            }
        };
        </script>`)

	// 输出Wrapper
	buf.WriteString(`<div class="expo-wrapper" id="expo-wrapper">
		<form name="export_form" method="POST" target="export_frame">`)
	// portal
	buf.WriteString("\n<input type=\"hidden\" name=\"portal\" value=\"")
	buf.WriteString(portal.PortalKey)
	buf.WriteString("\"/>\n")
	// params
	buf.WriteString(`<input type="hidden" name="params" value="" id="params"/>`)
	// token
	buf.WriteString("\n<input type=\"hidden\" name=\"token\" value=\"")
	buf.WriteString(token)
	buf.WriteString("\"/>\n")

	// 输出导出格式
	buf.WriteString(`
        <div class="expo-grp">
          <div class="expo-tit"><strong>选择导出格式</strong></div>
          <ul class="columnList">
            <li class="wbExpo_format_csv"><input type="radio" name="export_format" style="border:none"
              value="csv" checked="checked" id="wbExpo_format_csv"/>
                <label for="wbExpo_format_csv">CSV数据文件(推荐)</label>
            </li>
            <li class="wbExpo_format_excel"><input type="radio" name="export_format" style="border:none"
               value="excel" id="wbExpo_format_excel"/>
                <label for="wbExpo_format_excel">Excel</label>
            </li>
            <li class="wbExpo_format_txt"><input type="radio" name="export_format" style="border:none"
              value="txt" id="wbExpo_format_txt"/>
                <label for="wbExpo_format_txt">文本</label>
            </li>
          </ul>
        </div>
        <div style="clear:both"></div><br />`)

	colNames := portal.GetColumnMapping()
	hasCols := len(colNames) != 0
	if hasCols {
		// 输出勾选框
		buf.WriteString(`<div class="expo-grp expo-cols">
            <div class="expo-tit"><strong>请选择要导出的列:</strong></div>
            <ul class="columnList">`)
		for _, col := range colNames {
			buf.WriteString("<li><input type=\"checkbox\" style=\"border:none\" checked=\"checked\"")
			buf.WriteString(" name=\"export_field\" id=\"export_field_")
			buf.WriteString(col.Field)
			buf.WriteString("\" value=\"")
			buf.WriteString(col.Field)
			buf.WriteString("\"/><label for=\"export_field_")
			buf.WriteString(col.Field)
			buf.WriteString("\">")
			buf.WriteString(col.Name)
			buf.WriteString("</label></li>")
		}
		buf.WriteString("</ul><div style=\"clear:both\"></div></div>")
		// 输出勾选框
		buf.WriteString(`<div class="expo-grp">
            <div class="expo-tit"><strong>导出文件名:</strong></div>
            <input type="checkbox" id="ck-auto-name" checked="checked" onclick="$_expo.chkFileName(this)"/>
            <label for="ck-auto-name">自动生成文件名</label>&nbsp;&nbsp;
            <input type="text" name="file_name" class="expo-file-name" id="ck-file-txt" style="visibility:hidden"`)
		// 文件名
		buf.WriteString(" _value=\"")
		buf.WriteString(fileName)
		buf.WriteString("\" value=\"")
		buf.WriteString(fileName)
		buf.WriteString("\"/></div>")
		// 输出按钮
		buf.WriteString(`<iframe id="export_frame" name="export_frame" style="display:none"></iframe>
        <div style="clear:both"></div><input type="button" class="gra-btn expo-btn btn-export" onclick="$_expo.submit()"
         value=" 导出 "/>`)
	} else {
		buf.WriteString("<div class=\"expo-no-field\">该导出方案不包含可选择的导出列</div>")
	}
	// 输出表单结束标签
	buf.WriteString("</form></div>")

	return buf.String()

}
