/**
 * Copyright 2015 @ z3q.net.
 * name : template
 * author : jarryliu
 * date : 2016-06-01 18:10
 * description :
 * history :
 */
package gof

import (
	"bytes"
	"errors"
	"github.com/fsnotify/fsnotify"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

var (
	eventRegexp = regexp.MustCompile("\"(.+)\":\\s*(\\S+)")
	// cache template file of parent directory
	T_CACHE_PARENT = true
	T_SHOW_LOG     = false
)

type CacheTemplate struct {
	_basePath         string
	_shareFiles       []string
	_funcMap          template.FuncMap
	_fsNotify         bool
	_set              map[string]*template.Template
	_mux              *sync.RWMutex
	_winPathRegexp    *regexp.Regexp
	_includeMux       *sync.RWMutex
	_includeCache     map[string]string      // 包含子模板缓存
	_includeCheckFunc func(path string) bool // 检查子模板是否更改
}

// when notify is false , will not compile on file change!
func NewCacheTemplate(basePath string, notify bool, files ...string) *CacheTemplate {
	g := &CacheTemplate{
		_basePath:     basePath,
		_fsNotify:     notify,
		_set:          make(map[string]*template.Template, 0),
		_mux:          &sync.RWMutex{},
		_shareFiles:   files,
		_includeMux:   &sync.RWMutex{},
		_includeCache: make(map[string]string),
	}
	return g.init()
}

func (c *CacheTemplate) init() *CacheTemplate {
	// 初始化模板函数
	c._funcMap = map[string]interface{}{
		"include": c.include,
	}
	// 设置共享的子模板路径
	for i, v := range c._shareFiles {
		if !strings.HasPrefix(v, c._basePath) {
			c._shareFiles[i] = c._basePath + v
		}
	}
	// 开启监视
	if c._fsNotify {
		go c.fsNotify()
	}
	return c
}

func (c *CacheTemplate) println(err bool, v ...interface{}) {
	if err || T_SHOW_LOG {
		v = append([]interface{}{"[ Gof][ Template]"}, v...)
		log.Println(v...)
	}
}

// calling on file changed
func (c *CacheTemplate) fileChanged(event *fsnotify.Event) {
	c.resetIncludeCache()
	eventStr := event.String()
	if runtime.GOOS == "windows" {
		if c._winPathRegexp == nil {
			c._winPathRegexp = regexp.MustCompile("\\\\+")
		}
		eventStr = c._winPathRegexp.ReplaceAllString(eventStr, "/")
	}
	if eventRegexp.MatchString(eventStr) {
		matches := eventRegexp.FindAllStringSubmatch(eventStr, 1)
		if len(matches) > 0 {
			filePath := matches[0][1]
			if i := strings.Index(filePath, c._basePath); i != -1 {
				file := filePath[i+len(c._basePath):]
				if strings.Index(file, "_old_") == -1 &&
					strings.Index(file, "_tmp_") == -1 &&
					strings.Index(file, "_swp_") == -1 {
					c.handleChange(file) //do some things on file changed.
				}
			}
		}
	}
}

func (c *CacheTemplate) handleChange(file string) (err error) {
	filePath := c._basePath + file
	fi, err := os.Stat(filePath)
	if err == nil {
		if fi.IsDir() {
			c._mux.Lock()
			c._set = map[string]*template.Template{}
			c._mux.Unlock()
			return nil
		}
		_, err = c.compileTemplate(file) // recompile template
	}
	return err

	fullName := c._basePath + file
	for _, v := range c._shareFiles {
		if v == fullName {
			c._set = map[string]*template.Template{}
			break
		}
	}
	//todo: bug
	//if f, err := os.Stat(file); err == nil && !f.IsDir() {
	_, err = c.compileTemplate(file) // recompile template
	//}

	return err
}

// file system notify
func (c *CacheTemplate) fsNotify() {
	w, err := fsnotify.NewWatcher()
	if err == nil {
		// watch event
		go func(g *CacheTemplate) {
			for {
				select {
				case event := <-w.Events:
					if event.Op&fsnotify.Write != 0 ||
						event.Op&fsnotify.Create != 0 {
						g.fileChanged(&event)
					}
				case err := <-w.Errors:
					c.println(true, "[ Notify][ Error]:", err)
				}
			}
		}(c)

		err = filepath.Walk(c._basePath, func(path string,
			info os.FileInfo, err error) error {
			if err == nil && info.IsDir() {
				if n := info.Name(); n[0] != '.' && n[0] != ' ' {
					err = w.Add(path)
					//log.Println(info.Name())
				}
			}
			return err
		})
	}
	if err != nil {
		w.Close()
		panic(err)
		os.Exit(0)
		return
	}
	<-make(chan bool)
}

func (c *CacheTemplate) parseTemplate(name string) (
	*template.Template, error) {
	// 主要的模板文件,需要第一个位置
	files := append([]string{c._basePath + name}, c._shareFiles...)
	// 新建模板并设置模板函数
	t := template.New(c.tplName(name))
	if c._funcMap != nil {
		t = t.Funcs(c._funcMap)
	}
	return t.ParseFiles(files...)
}

func (c *CacheTemplate) compileTemplate(name string) (
	*template.Template, error) {
	c._mux.Lock()
	defer c._mux.Unlock()
	tpl, err := c.parseTemplate(name)
	if err == nil {
		// 上一级选择性缓存
		if T_CACHE_PARENT || (strings.Index(name, "../") == -1 &&
			strings.Index(name, "..\\") == -1) {
			c._set[name] = tpl
		}
		c.println(false, "[ Compile]: ", name)
	} else {
		c.println(true, "[ Compile][ Error]: ", err.Error())
	}

	return tpl, err
}

func (c *CacheTemplate) Funcs(funcMap template.FuncMap) {
	for k, v := range funcMap {
		if _, ok := c._funcMap[k]; ok {
			panic(errors.New("exists func " + k))
		}
		c._funcMap[k] = v
	}
}

// 获取模板的名称
func (c *CacheTemplate) tplName(path string) string {
	if li := strings.LastIndex(path, "/"); li != -1 {
		return path[li+1:]
	}
	return path
}

// 清除子模板的缓存
func (c *CacheTemplate) resetIncludeCache() {
	c._includeMux.Lock()
	defer c._includeMux.Unlock()
	c._includeCache = make(map[string]string)
}

// 子模板中间件，如果需要更新缓存，则返回false
func (c *CacheTemplate) IncludeMiddle(f func(path string) bool) {
	c._includeCheckFunc = f
}

// 读取子模板的
func (c *CacheTemplate) include(path string) template.HTML {
	c._includeMux.RLock()
	str, ok := c._includeCache[path]
	c._includeMux.RUnlock()
	if ok {
		if c._includeCheckFunc == nil || c._includeCheckFunc(path) {
			return template.HTML(str)
		}
	}
	c._includeMux.Lock()
	defer c._includeMux.Unlock()
	data, err := c.read(path)
	if err != nil {
		return template.HTML(err.Error())
	}
	result := string(data)
	c._includeCache[path] = result
	return template.HTML(result)
}

func (c *CacheTemplate) read(path string) ([]byte, error) {
	var err error
	buf := bytes.NewBuffer(nil)
	tpl := template.New(path)
	tpl.Funcs(c._funcMap)
	tpl, err = tpl.ParseFiles(c._basePath + path)
	if err == nil {
		err = tpl.ExecuteTemplate(buf, c.tplName(path), nil)
	}
	return buf.Bytes(), err
}

func (c *CacheTemplate) Execute(w io.Writer,
	name string, data interface{}) (err error) {
	c._mux.RLock() //仅对读加锁
	tpl, ok := c._set[name]
	c._mux.RUnlock()
	if !ok {
		if tpl, err = c.compileTemplate(name); err != nil {
			return err
		}
	}
	return tpl.Execute(w, data)
}
