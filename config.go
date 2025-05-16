/**
 * config file end with enter line
https://github.com/vaughan0/go-ini/blob/master/ini.go
*/

package gof

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const confLineDelimer byte = '\n'

var (
	regex = regexp.MustCompile("^\\s*([^#\\s]+)\\s*=\\s*\"*(.+)*\"*\\s*$")
)

// Config
type Config struct {
	configDict map[string]interface{}
}

func NewConfig() *Config {
	return &Config{
		configDict: make(map[string]interface{}),
	}
}

// 从文件中加载配置
func LoadConfig(file string) (cfg *Config, err error) {
	s := &Config{}
	_err := s.load(file)
	return s, _err
}

// 从配置中读取数据
func (c *Config) GetString(key string) string {
	k, e := c.configDict[key]
	if e {
		v, _ := k.(string)
		return v
	}
	return ""
}

// 从配置中读取数据
func (c *Config) Get(key string) interface{} {
	v, e := c.configDict[key]
	if e {
		return v
	}
	return nil
}

func (c *Config) Set(key string, v interface{}) {
	if _, ok := c.configDict[key]; ok {
		panic("Key '" + key + "' is exist in config")
	}
	c.configDict[key] = v
}

func (c *Config) GetInt(key string) int {
	k, e := c.configDict[key]
	if e {
		v, ok := k.(int)
		if ok {
			return v
		}
		if sv, ok := k.(string); ok {
			if iv, err := strconv.Atoi(sv); err == nil {
				return iv
			}
		}
	}
	return 0
}

func (c *Config) GetFloat(key string) float64 {
	k, e := c.configDict[key]
	if e {
		v, ok := k.(float64)
		if ok {
			return v
		}
		if sv, ok := k.(string); ok {
			if iv, err := strconv.ParseFloat(sv, 64); err == nil {
				return iv
			}
		}
	}
	return 0
}

// 从文件中加载配置
func (c *Config) load(file string) (err error) {
	c.configDict = make(map[string]interface{})
	//var allContent string = ""
	f, _err := os.Open(file)
	if _err != nil {
		return _err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	prefix := ""
	for {
		line, _err := reader.ReadString(confLineDelimer)
		if _err == io.EOF {
			break
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section := line[1 : len(line)-1]
			if len(section) > 0 {
				prefix = section + "."
				continue
			}
		}
		if regex.Match([]byte(line)) {
			matches := regex.FindStringSubmatch(line)
			k := strings.TrimSpace(matches[1])
			v := strings.TrimSpace(matches[2])
			if len(prefix) > 0 {
				k = prefix + k
			}
			c.configDict[k] = v
		}
	}
	return nil
}
