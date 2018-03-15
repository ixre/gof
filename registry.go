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
	data    map[string]*RegistryTree
}

type RegistryTree struct {
	key      string
	filePath string
	isRoot   bool
	tree     *toml.Tree
	glob     *Registry
}

// create new registry manager
func NewRegistry(path string, delimer string) (*Registry, error) {
	if delimer == "" {
		delimer = "."
	}
	return (&Registry{
		path:    path,
		delimer: delimer,
	}).init()
}

// load from config files
func (r *Registry) init() (*Registry, error) {
	r.data = map[string]*RegistryTree{}
	return r.dirInit()
}

func (r *Registry) dirInit() (*Registry, error) {
	d, err := os.Open(r.path)
	if err != nil {
		return r, err
	}
	files, err := d.Readdir(-1)
	if err != nil {
		return r, err
	}
	for _, f := range files {
		r.load(r.path+"/"+f.Name(), f)
	}
	return r, nil
}

func (r *Registry) load(path string, info os.FileInfo) error {
	if info.IsDir() {
		return nil
	}
	if name := info.Name(); strings.HasSuffix(name, ".conf") {
		t, err := toml.LoadFile(path)
		if err == nil {
			key := name[:len(name)-5]
			tree := &RegistryTree{
				isRoot:   true,
				key:      key,
				filePath: path,
				tree:     t,
				glob:     r,
			}
			r.data[key] = tree
		}
		return err
	}
	return nil
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
	i := strings.Index(key, r.delimer)
	return r.Use(key[:i]).Get(key[i+1:])
}

func (r *Registry) Set(key string, v interface{}) error {
	i := strings.Index(key, r.delimer)
	rt := r.Use(key[:i])
	err := rt.Set(key[i+1:], v)
	if err == nil {
		err = rt.Flush()
	}
	return err
}

func (r *Registry) Use(key string) *RegistryTree {
	t, b := r.data[key]
	if !b {
		t = &RegistryTree{
			isRoot:   true,
			key:      key,
			filePath: r.path + key + ".conf",
			tree:     nil,
			glob:     r,
		}
	}
	return t
}

// get string
func (r *Registry) GetString(key string) string {
	defer func() {
		if err := recover(); err != nil {
			panic("key=" + key)
		}
	}()
	return r.Get(key).(string)
}

func (r *Registry) createNode(arr []string, value interface{}) (*toml.Tree, error) {
	key := arr[0]
	tree, err := toml.TreeFromMap(map[string]interface{}{})
	if len(arr) == 1 {
		tree.Set(key, value)
	} else {
		tree2, _ := toml.TreeFromMap(map[string]interface{}{
			arr[1]: value,
		})
		tree.Set(key, tree2)
	}
	return tree, err
}

func (r *RegistryTree) Exists() bool {
	return r.tree != nil
}

// get registry pair
func (r *RegistryTree) Get(prop string) interface{} {
	if r.tree != nil {
		arr := r.desSplit(prop)
		if len(arr) == 1 {
			return r.tree.Get(arr[0])
		}
		tree, result := r.tree.Get(arr[0]).(*toml.Tree)
		if !result || tree == nil {
			return nil
		}
		return tree.Get(arr[1])
	}
	return nil
}

func (r *RegistryTree) GetString(prop string) string {
	return r.Get(prop).(string)
}

func (r *RegistryTree) GetInt(prop string) int {
	return r.Get(prop).(int)
}

func (r *RegistryTree) GetBool(prop string) bool {
	return r.Get(prop).(bool)
}

// get file key and config key
func (r *RegistryTree) desSplit(key string) []string {
	i := strings.LastIndex(key, r.glob.delimer)
	if i > 0 {
		return []string{key[:i], key[i+1:]}
	}
	return []string{key}
}

// set registry pair
func (r *RegistryTree) Set(prop string, value interface{}) error {
	arr := r.desSplit(prop)
	if r.tree == nil {
		tree, err := r.glob.createNode(arr, value)
		if err != nil {
			return err
		}
		r.tree = tree
		r.glob.data[r.key] = r
		return nil
	}
	if len(arr) == 1 {
		r.tree.Set(arr[0], value)
	} else {
		tree := r.tree.Get(arr[0]).(*toml.Tree)
		if tree == nil {
			return errors.New("no such node " + arr[0])
		}
		tree.Set(arr[1], value)
	}
	return nil
}

func (r *RegistryTree) Flush() error {
	dir := filepath.Dir(r.filePath)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
	}
	if err == nil {
		var fi *os.File
		fi, err = os.OpenFile(r.filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err == nil {
			_, err = r.tree.WriteTo(fi)
			fi.Close()
		}
	}
	return err
}
