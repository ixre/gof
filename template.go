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
	"gopkg.in/fsnotify.v1"
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
)

type CachedTemplate struct {
	basePath      string
	shareFiles    []string
	set           map[string]*template.Template
	mux           *sync.RWMutex
	winPathRegexp *regexp.Regexp
}

func NewCachedTemplate(basePath string, files ...string) *CachedTemplate {
	g := &CachedTemplate{
		basePath:   basePath,
		set:        make(map[string]*template.Template, 0),
		mux:        &sync.RWMutex{},
		shareFiles: files,
	}
	return g.init()
}

func (this *CachedTemplate) init() *CachedTemplate {
	for i, v := range this.shareFiles {
		if !strings.HasPrefix(v, this.basePath) {
			this.shareFiles[i] = this.basePath + v
		}
	}
	go this.fsNotify()
	return this
}

// calling on file changed
func (this *CachedTemplate) fileChanged(event *fsnotify.Event) {
	eventStr := event.String()
	if runtime.GOOS == "windows" {
		if this.winPathRegexp == nil {
			this.winPathRegexp = regexp.MustCompile("\\\\+")
		}
		eventStr = this.winPathRegexp.ReplaceAllString(eventStr, "/")
	}
	if eventRegexp.MatchString(eventStr) {
		matches := eventRegexp.FindAllStringSubmatch(eventStr, 1)
		if len(matches) > 0 {
			filePath := matches[0][1]
			if i := strings.Index(filePath, this.basePath); i != -1 {
				file := filePath[i+len(this.basePath):]
				if strings.Index(file, "_old_") == -1 &&
					strings.Index(file, "_tmp_") == -1 &&
					strings.Index(file, "_swp_") == -1 {
					//todo: bug
					//if f, err := os.Stat(file); err == nil && !f.IsDir() {
					this.compileTemplate(file) // recompile template
					//}
				}
			}
		}
	}
}

// file system notify
func (this *CachedTemplate) fsNotify() {
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
		}(this)

		err = filepath.Walk(this.basePath, func(path string,
			info os.FileInfo, err error) error {
			if err == nil && info.IsDir() &&
				info.Name()[0] != '.' { // not hidden file
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

func (this *CachedTemplate) parseTemplate(name string) (
	*template.Template, error) {
	this.mux.Lock() //对写加锁
	files := append([]string{this.basePath + name},
		this.shareFiles...) //name需要第一个位置
	tpl, err := template.ParseFiles(files...)
	this.mux.Unlock()
	return tpl, err
}

func (this *CachedTemplate) compileTemplate(name string) (
	*template.Template, error) {
	tpl, err := this.parseTemplate(name)
	if err == nil {
		this.set[name] = tpl
		log.Println("[ Gof][ Template][ Compile]: ", name)
	} else {
		log.Println("[ Gof][ Template][ Error] -", err.Error())
	}
	return tpl, err
}

func (this *CachedTemplate) Execute(w io.Writer,
	name string, data interface{}) error {
	this.mux.RLock() //仅对读加锁
	tpl, ok := this.set[name]
	if !ok {
		this.mux.RUnlock()
		var err error
		if tpl, err = this.compileTemplate(name); err != nil {
			return err
		}
		this.set[name] = tpl
	} else {
		defer this.mux.RUnlock()
	}
	return tpl.Execute(w, data)
}

func (this *CachedTemplate) Render(w io.Writer,
	name string, data interface{}) error {
	return this.Execute(w, name, data)
}
