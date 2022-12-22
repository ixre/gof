package util

import (
	"os"
	"path/filepath"
)

// #https://groups.google.com/g/golang-codereviews/c/PYR9LB_YY4E?pli=1

// Save bytes to file
func BytesToFile(data []byte, path string) (err error) {
	dPath := filepath.Dir(path)
	if _, err = os.Stat(dPath); os.IsNotExist(err) {
		err = os.MkdirAll(dPath, os.ModePerm)
	}
	if err == nil {
		err = os.WriteFile(path, data, os.ModePerm)
	}
	return err
}
