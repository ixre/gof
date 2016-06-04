/**
 * Copyright 2015 @ z3q.net.
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
func (this *GobFile) Unmarshal(dst interface{}) error {
	this.mux.Lock()
	fi, err := os.Open(this.File)
	if err == nil {
		enc := gob.NewDecoder(fi)
		err = enc.Decode(dst)
	}
	this.mux.Unlock()
	return err
}

// 序列化并存储到文件
func (this *GobFile) marshal(src interface{}) error {
	//检测目录是否存在,不存在则创建目录
	dir := path.Dir(this.File)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	this.mux.Lock()
	f, err := os.OpenFile(this.File,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		os.ModePerm)
	if err == nil {
		enc := gob.NewEncoder(f)
		err = enc.Encode(src)
	}
	this.mux.Unlock()
	return err
}

// 保存到文件中
func (this *GobFile) Save(src interface{}) error {
	return this.marshal(src)
}
