/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-25 21:23
 * description :
 * history :
 */

package crypto

//
//cyp := NewMd5Crypto("ops", "rdm")
//i :=2
//for {
//if i = i - 1; i ==0  {
//break
//}
//
//str := cyp.Encode()
//fmt.Println("str:", str)
//
//_,unix := cyp.Decode(str)
//fmt.Println(time.Now().Unix()-unix)
//}

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"sync"
	"time"
)

var (
	d5 = md5.New()
)

const (
	unixLen = 10 //unix time长度为10
)

func getPos(token string) int {
	return len(token)/2 + 1
}
func getUnix() string {
	ux := time.Now().Unix()
	return strconv.FormatInt(ux, 10)
}

func getMd5(token, offset string) []byte {
	d5.Reset()
	d5.Write([]byte(token))
	src := d5.Sum([]byte(offset))
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst
}

type UnixCrypto struct {
	pos      int
	md5Bytes []byte
	buf      *bytes.Buffer
	mux      sync.Mutex
}

func NewUnixCrypto(token, offset string) *UnixCrypto {
	return &UnixCrypto{
		pos:      len(token)/2 + 1,
		md5Bytes: getMd5(token, offset),
		buf:      bytes.NewBufferString(""),
	}
}

// return md5 bytes
func (u *UnixCrypto) GetBytes() []byte {
	return u.md5Bytes
}

func (uc *UnixCrypto) Encode() []byte {
	unx := getUnix()
	l := uc.pos

	uc.mux.Lock()
	defer func() {
		uc.buf.Reset()
		uc.mux.Unlock()
	}()

	uc.buf.Write(uc.md5Bytes[:l])

	for i := 0; i < 10; i++ {
		uc.buf.WriteString(unx[i : i+1])
		uc.buf.Write(uc.md5Bytes[l+i : l+i+1])
	}

	uc.buf.Write(uc.md5Bytes[10+l:])
	return uc.buf.Bytes()
}

func (u *UnixCrypto) Decode(s string) ([]byte, int64) {
	smd := make([]byte, len(u.md5Bytes))
	unx := make([]byte, unixLen)

	if len(s) < len(smd) {
		return nil, 0
	}

	copy(smd, s[:u.pos])
	for i, v := range s[u.pos+unixLen*2:] {
		smd[u.pos+unixLen+i] = byte(v)
	}

	for i := 0; i < unixLen*2; i++ {
		v := s[u.pos+i]
		if i%2 == 0 {
			unx[i/2] = v
		} else {
			smd[u.pos+i/2] = v
		}
	}

	unix, err := strconv.ParseInt(string(unx), 10, 32)
	if err != nil {
		unix = 0
	}
	return smd, unix
}

func (uc *UnixCrypto) Compare(s string) (bool, []byte, int64) {
	b, u := uc.Decode(s)
	return bytes.Compare(b, uc.md5Bytes) == 0, b, u
}
