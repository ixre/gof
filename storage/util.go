package storage

import (
	"bytes"
	"encoding/gob"
	"strings"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : util.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-09-30 10:07
 * description :
 * history :
 */

// encode bytes from interface
func EncodeBytes(v interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(v)
	checkTypeErr(err)
	return buf.Bytes(), err
}

// decode to interface
func DecodeBytes(b []byte, dst interface{}) error {
	buf := bytes.NewBuffer(b)
	err := gob.NewDecoder(buf).Decode(dst)
	checkTypeErr(err)
	return err
}

func checkTypeErr(err error) {
	if err != nil && strings.Index(err.Error(), "type not registered") != -1 {
		panic(err)
	}
}
