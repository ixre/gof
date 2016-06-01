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
	"strings"
	"sync"
)

var (
	eventRegexp = regexp.MustCompile("\"(.+)\":\\s*(\\S+)")
)

type CachedTemplate struct {
	baseDirectory string
	set           map[string]*template.Template
	mux           *sync.RWMutex
}

func NewCachedTemplate(dir string) *CachedTemplate {
	g := &CachedTemplate{
		baseDirectory: dir,
		set:           make(map[string]*template.Template, 0),
		mux:           &sync.RWMutex{},
	}
	return g.init()
}

func (g *CachedTemplate) init() *CachedTemplate {
	go g.fsNotify()
	return g
}

// calling on file changed
func (this *CachedTemplate) fileChanged(event *fsnotify.Event) {
	if eventRegexp.MatchString(event.String()) {
		matches := eventRegexp.FindAllStringSubmatch(event.String(), 1)
		if len(matches) > 0 {
			filePath := matches[0][1]
			if i := strings.Index(filePath, this.baseDirectory); i != -1 {
				file := filePath[i+len(this.baseDirectory):]
				if strings.Index(file, "_old_") == -1 &&
					strings.Index(file, "_tmp_") == -1 &&
					strings.Index(file, "_swp_") == -1 {
					this.compileTemplate(file) // recompile template
				}
			}
		}
	}
}

// file system notify
func (this *CachedTemplate) fsNotify() {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
		os.Exit(0)
	}
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

	filepath.Walk(this.baseDirectory, func(path string,
		info os.FileInfo, err error) error {
		if info.IsDir() && info.Name()[0] != '.' {
			return w.Add(path)
		}
		return nil
	})
	var ch chan bool = make(chan bool)
	<-ch
	w.Close()
}

func (this *CachedTemplate) parseTemplate(name string) (
	*template.Template, error) {
	this.mux.Lock() //对写加锁
	tpl, err := template.ParseFiles(this.baseDirectory + name)
	this.mux.Unlock()
	return tpl, err
}

func (this *CachedTemplate) compileTemplate(name string) (
	*template.Template, error) {
	//this.mux.Lock() //仅对读加锁
	tpl, err := this.parseTemplate(name)
	if err == nil {
		this.set[name] = tpl
		log.Println("[ Gof][ Template][ Compile]: ", name)
	} else {
		log.Println("[ Gof][ Template][ Error] -", err.Error())
	}
	//this.mux.Unlock()
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
	}else{
        defer this.mux.RUnlock()
    }
	return tpl.Execute(w, data)
}

func (this *CachedTemplate) Render(w io.Writer,
	name string, data interface{}) error {
	return this.Execute(w, name, data)
}
