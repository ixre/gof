package rsa

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

//RSA公钥私钥产生,bits: 1024 2048 4096
func GenRsaKeys(bits int) (publicKeyStr, privateKeyStr string, err error) {
	str := func(block *pem.Block)(s string,err error) {
		buffer := new(bytes.Buffer)
		err = pem.Encode(buffer, block)
		if err == nil {
			s = buffer.String()
		}
		return
	}
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return
	}
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes:  x509.MarshalPKCS1PrivateKey(privateKey),
	}
	privateKeyStr,err = str(block)
	if err != nil {
		return
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	pubBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err == nil {
		block = &pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: pubBytes,
		}
		publicKeyStr, err = str(block)
	}
	return
}

