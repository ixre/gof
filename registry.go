package gof

import (
	"errors"
	"github.com/pelletier/go-toml"
	"os"
	"path/filepath"
	"strings"
)

// registry
type Registry struct {
	path    string
	delimer string
	data    map[string]*toml.Tree
	pathMap map[string]string
}

func NewRegistry(path string, delimer string) (*Registry, error) {
	return (&Registry{
		path:    path,
		delimer: delimer,
	}).init()
}

// load from config files
func (r *Registry) init() (*Registry, error) {
	var err error
	r.data = map[string]*toml.Tree{}
	r.pathMap = map[string]string{}
	err = filepath.Walk(r.path, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			name := info.Name()
			if strings.HasSuffix(name, ".conf") {
				t, err := toml.LoadFile(path)
				if err == nil {
					key := name[:len(name)-5]
					r.data[key] = t
					r.pathMap[key] = path
				}
				return err
			}
		}
		return nil
	})
	return r, err
}

// get file key and config key
func (r *Registry) split(key string) []string {
	arr := strings.Split(key, r.delimer)
	if len(arr) <= 1 {
		panic("registry key not contain delimer '" + r.delimer + "'")
	}
	return arr
}

// get registry pair
func (r *Registry) Get(key string) interface{} {
	arr := r.split(key)
	fk, tk := arr[0], arr[1]
	d, exist := r.data[fk]
	if !exist {
		return "no such registry " + fk
	}
	if len(arr) == 2 {
		return d.Get(tk)
	}
	tree, result := d.Get(tk).(*toml.Tree)
	if !result || tree == nil {
		return "no such node " + tk
	}
	return tree.Get(arr[2])
}

// set registry pair
func (r *Registry) Set(key string, value interface{}) (err error) {
	arr := r.split(key)
	fk, tk := arr[0], arr[1]
	d, exist := r.data[fk]
	if !exist {
		d, err = r.createNode(arr, value)
		r.data[fk] = d
		r.pathMap[fk] = r.path + fk + ".conf"
	} else {
		if len(arr) == 2 {
			d.Set(tk, value)
		} else {
			tree := d.Get(tk).(*toml.Tree)
			if tree == nil {
				return errors.New("no such node " + tk)
			}
			tree.Set(arr[2], value)
		}
	}
	_, err = os.Stat(r.path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(r.path, os.ModePerm)
	}
	if err == nil {
		var fi *os.File
		fi, err = os.OpenFile(r.pathMap[fk], os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err == nil {
			_, err = d.WriteTo(fi)
			fi.Close()
		}
	}
	return err
}

// get string
func (r *Registry) GetString(key string) string {
	return r.Get(key).(string)
}
func (r *Registry) createNode(arr []string, value interface{}) (*toml.Tree, error) {
	tk := arr[1]
	tree, err := toml.TreeFromMap(map[string]interface{}{})
	if len(arr) == 2 {
		tree.Set(tk, value)
	} else {
		tree2, _ := toml.TreeFromMap(map[string]interface{}{
			arr[2]: value,
		})
		tree.Set(tk, tree2)
	}
	return tree, err
}
