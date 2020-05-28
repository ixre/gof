package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// Save bytes to file
func BytesToFile(data []byte, path string) (err error) {
	dPath := filepath.Dir(path)
	if _, err = os.Stat(dPath); os.IsNotExist(err) {
		err = os.MkdirAll(dPath, os.ModePerm)
	}
	if err == nil {
		err = ioutil.WriteFile(path, data, os.ModePerm)
	}
	return err
}
