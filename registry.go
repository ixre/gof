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

func NewRegistry(path string, delimer string) *Registry {
	return (&Registry{
		path:    path,
		delimer: delimer,
	}).init()
}

// load from config files
func (r *Registry) init() *Registry {
	r.data = map[string]*toml.Tree{}
	r.pathMap = map[string]string{}
	filepath.Walk(r.path, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			name := info.Name()
			if strings.HasSuffix(name, ".conf") {
				t, err := toml.LoadFile(path)
				if err == nil {
					key := name[:len(name)-4]
					r.data[key] = t
					r.pathMap[key] = path
				}
				return err
			}
		}
		return nil
	})
	return r
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
	d, exist := r.data[arr[0]]
	if !exist {
		return "no such registry" + arr[0]
	}
	if len(arr) == 2 {
		return d.Get(arr[1])
	}
	tree := d.Get(arr[1]).(*toml.Tree)
	if tree == nil {
		return "no such node " + arr[1]
	}
	return tree.Get(arr[2])
}

// set registry pair
func (r *Registry) Set(key string, value interface{}) (err error) {
	arr := r.split(key)
	p := arr[0]
	d, exist := r.data[p]
	if !exist {
		d, err = r.createNode(arr, value)
		r.pathMap[arr[0]] = r.path + arr[0] + ".conf"
	} else {
		if len(arr) == 2 {
			d.Set(arr[1], value)
		} else {
			tree := d.Get(arr[1]).(*toml.Tree)
			if tree == nil {
				return errors.New("no such node " + arr[1])
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
		fi, err = os.OpenFile(r.pathMap[p], os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err == nil {
			_, err = d.WriteTo(fi)
			fi.Close()
		}
	}
	return err
}
func (r *Registry) createNode(arr []string, value interface{}) (*toml.Tree, error) {
	mp := map[string]interface{}{}
	if len(arr) == 2 {
		mp[arr[1]] = value
	} else {
		tree, _ := toml.TreeFromMap(map[string]interface{}{
			arr[2]: value,
		})
		mp[arr[1]] = tree
	}
	return toml.TreeFromMap(map[string]interface{}{
		arr[0]: mp,
	})
}
