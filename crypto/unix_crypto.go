/**
 * Copyright 2014 @ to2.net.
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
	"errors"
	"strconv"
	"time"
)

var (
	d5 = md5.New()
)

const (
	unixLen = 10 //unix time长度为10
)

// Unix时间戳加密
type UnixCrypto struct {
	pos      int
	md5Bytes []byte
	//buf      *bytes.Buffer
	//mux      sync.Mutex
}

func NewUnixCrypto(token, offset string) *UnixCrypto {
	u := &UnixCrypto{}
	u.pos = int(token[2])%9+1
	u.md5Bytes = u.getMd5(token, offset)
	return u
}

func (u *UnixCrypto) getMd5(token, offset string) []byte {
	d5.Reset()
	d5.Write([]byte(token))
	src := d5.Sum([]byte(offset))
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst[8:24]
}

// 获取UNIX时间戳字符串
func (u *UnixCrypto) getUnix() string {
	ux := time.Now().Unix()
	return strconv.FormatInt(ux, 10)
}

// return md5 bytes
func (u *UnixCrypto) GetBytes() []byte {
	return u.md5Bytes
}

// 编码
func (u *UnixCrypto) Encode() []byte {
	unixStr := u.getUnix()
	l := u.pos
	buf := bytes.NewBuffer(nil)
	buf.Write(u.md5Bytes[:l])
	//前10位，unix反序和md5交叉
	for i := 0; i < 10; i++ {
		buf.WriteString(unixStr[9-i : 9-i+1])
		buf.Write(u.md5Bytes[l+i : l+i+1])
	}
	//拼接md5，10位以后的字符
	buf.Write(u.md5Bytes[10+l:])
	return buf.Bytes()
}

// 解码，返回Token及Unix时间
func (u *UnixCrypto) Decode(result []byte) (token []byte, unix int64,err error) {
	//解码得到的token
	l := len(u.md5Bytes)
	token = make([]byte,l)
	if len(result) < l{
		return nil, 0,errors.New("decode bytes invalid length")
	}
	unixArr := make([]byte, unixLen)
	// 解码第一部分
	copy(token, result[:u.pos])
	// 解码的二部分,如果长度不匹配
	p2 := result[u.pos+unixLen*2:]
	if u.pos+unixLen + len(p2) != l {
		return nil, 0,errors.New("decode bytes sign not match")
	}
	for i, v := range p2 {
		token[u.pos+unixLen+i] = byte(v)
	}
	// 解码第三部分
	for i := 0; i < unixLen*2; i++ {
		v := result[u.pos+i]
		if i%2 == 0 {
			unixArr[9-i/2] = v
		} else {
			token[u.pos+i/2] = v
		}
	}
	unix, err = strconv.ParseInt(string(unixArr), 10, 32)
	if err != nil {
		unix = 0
	}
	return token, unix,err
}

func (uc *UnixCrypto) Compare(result []byte) (match bool, token []byte, unix int64) {
	token, unix, err := uc.Decode(result)
	if err == nil {
		match = bytes.Compare(token, uc.md5Bytes) == 0
	}
	return match, token, unix
}
