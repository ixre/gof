/**
 * Copyright 2015 @ to2.net.
 * name : gob_file.go
 * author : jarryliu
 * date : 2016-06-04 15:23
 * description :
 * history :
 */
package util

import (
	"encoding/gob"
	"os"
	"path"
	"sync"
)

type GobFile struct {
	mux  sync.Mutex
	File string
}

func NewGobFile(file string) *GobFile {
	return &GobFile{
		File: file,
	}
}

// 从文件中反序列化
func (g *GobFile) Unmarshal(dst interface{}) error {
	g.mux.Lock()
	fi, err := os.Open(g.File)
	if err == nil {
		enc := gob.NewDecoder(fi)
		err = enc.Decode(dst)
	}
	g.mux.Unlock()
	return err
}

// 序列化并存储到文件
func (g *GobFile) marshal(src interface{}) error {
	//检测目录是否存在,不存在则创建目录
	dir := path.Dir(g.File)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	g.mux.Lock()
	f, err := os.OpenFile(g.File,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		os.ModePerm)
	if err == nil {
		enc := gob.NewEncoder(f)
		err = enc.Encode(src)
	}
	g.mux.Unlock()
	return err
}

// 保存到文件中
func (g *GobFile) Save(src interface{}) error {
	return g.marshal(src)
}
