package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

func Md5(b []byte) string {
	c := md5.New()
	c.Write(b)
	return hex.EncodeToString(c.Sum(nil))
}
func Sha1(b []byte) string {
	c := sha1.New()
	c.Write(b)
	return hex.EncodeToString(c.Sum(nil))
}
