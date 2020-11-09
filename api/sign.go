package api

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"net/url"
	"sort"
	"strings"
)

// 参数首字母小写后排序，排除sign和sign_type，secret，转换为字节
func ParamsToBytes(r url.Values, secret string, attach bool) []byte {
	keys := keyArr{}
	for k := range r {
		keys = append(keys, k)
	}
	sort.Sort(keys)
	// 拼接参数和值
	i := 0
	buf := bytes.NewBuffer(nil)
	for _, k := range keys {
		if k == "sign" || k == "sign_type" {
			continue
		}
		if i > 0 {
			buf.WriteString("&")
		}
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(r[k][0])
		i++
	}
	if attach {
		buf.WriteString(secret)
	}
	return buf.Bytes()
}

// 签名
func Sign(signType string, r url.Values, secret string) string {
	data := ParamsToBytes(r, secret, true)
	switch signType {
	case "md5":
		return byteHash(md5.New(), data)
	case "sha1":
		return byteHash(sha1.New(), data)
	}
	return ""
}

// 计算Hash值
func byteHash(h hash.Hash, data []byte) string {
	h.Write(data)
	b := h.Sum(nil)
	return hex.EncodeToString(b)
}

/*------ other support code ------*/
var _ sort.Interface = keyArr{}

type keyArr []string

func (s keyArr) Len() int {
	return len(s)
}

func (s keyArr) Less(i, j int) bool {
	return strings.ToLower(s[i]) < strings.ToLower(s[j])
}

func (s keyArr) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
