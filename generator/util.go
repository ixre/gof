package generator

import (
	"github.com/jsix/gof/util"
	"regexp"
)

var (
	revertRegexp = regexp.MustCompile("\\$\\{([^\\}]+)\\}")
)

// 保存到文件
func SaveFile(s string, path string) error {
	return util.BytesToFile([]byte(s), path)
}

// 还原模板的标签: ${...} -> {{...}}
func RevertTPVariable(str string) string {
	return revertRegexp.ReplaceAllString(str, "{{$1}}")
}
