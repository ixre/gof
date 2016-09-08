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
)

type CachedTemplate struct {
	_basePath      string
	_shareFiles    []string
	_fsNotify      bool
	_set           map[string]*template.Template
	_mux           *sync.RWMutex
	_winPathRegexp *regexp.Regexp
}

// when notify is false , will not compile on file change!
func NewCachedTemplate(basePath string, notify bool, files ...string) *CachedTemplate {
	g := &CachedTemplate{
		_basePath:   basePath,
		_fsNotify:   notify,
		_set:        make(map[string]*template.Template, 0),
		_mux:        &sync.RWMutex{},
		_shareFiles: files,
	}
	return g.init()
}

func (c *CachedTemplate) init() *CachedTemplate {
	for i, v := range c._shareFiles {
		if !strings.HasPrefix(v, c._basePath) {
			c._shareFiles[i] = c._basePath + v
		}
	}
	if c._fsNotify {
		go c.fsNotify()
	}
	return c
}

// calling on file changed
func (c *CachedTemplate) fileChanged(event *fsnotify.Event) {
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

func (c *CachedTemplate) handleChange(file string) (err error) {
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
func (c *CachedTemplate) fsNotify() {
	w, err := fsnotify.NewWatcher()
	if err == nil {
		// watch event
		go func(g *CachedTemplate) {
			for {
				select {
				case event := <-w.Events:
					if event.Op&fsnotify.Write != 0 ||
						event.Op&fsnotify.Create != 0 {
						g.fileChanged(&event)
					}
				case err := <-w.Errors:
					log.Println("Error:", err)
				}
			}
		}(c)

		err = filepath.Walk(c._basePath, func(path string,
			info os.FileInfo, err error) error {
			if err == nil && info.IsDir() &&
				info.Name()[0] != '.' {
				// not hidden file
				err = w.Add(path)
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

func (c *CachedTemplate) parseTemplate(name string) (
	*template.Template, error) {
	files := append([]string{c._basePath + name},
		c._shareFiles...) //name需要第一个位置
	return template.ParseFiles(files...)
}

func (c *CachedTemplate) compileTemplate(name string) (
	*template.Template, error) {
	c._mux.Lock()
	defer c._mux.Unlock()
	tpl, err := c.parseTemplate(name)
	if err == nil {
		if T_CACHE_PARENT || (strings.Index(name, "../") == -1 &&
			strings.Index(name, "..\\") == -1) {
			c._set[name] = tpl
		}
		log.Println("[ Gof][ Template][ Compile]: ", name)
	} else {
		log.Println("[ Gof][ Template][ Error] -", err.Error())
	}
	return tpl, err
}

func (c *CachedTemplate) Execute(w io.Writer,
	name string, data interface{}) error {
	c._mux.RLock() //仅对读加锁
	tpl, ok := c._set[name]
	c._mux.RUnlock()
	if !ok {
		var err error
		if tpl, err = c.compileTemplate(name); err != nil {
			return err
		}
	}
	return tpl.Execute(w, data)
}

func (c *CachedTemplate) Render(w io.Writer,
	name string, data interface{}) error {
	return c.Execute(w, name, data)
}
