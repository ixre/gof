package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"strings"
)

//RSA公钥私钥产生,bits: 1024 2048 4096
func GenRsaKeys(bits int) (publicKeyStr, privateKeyStr string, err error) {
	str := func(block *pem.Block) (s string, err error) {
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
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	privateKeyStr, err = str(block)
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


// 转换私钥
func ParsePrivateKey(pemKey string)(*rsa.PrivateKey,error){
	if !strings.HasPrefix(pemKey,"-----BEGIN PRIVATE KEY-----"){
		pemKey = "-----BEGIN PRIVATE KEY-----\n" + pemKey+"\n-----END PRIVATE KEY-----"
	}
	block, _ := pem.Decode([]byte(pemKey))
	if block == nil{
		return nil,errors.New("pem private key not incorrect")
	}
	e,err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err == nil{
		return e.(*rsa.PrivateKey),err
	}
	return nil,err
}

// 转换公钥
func ParsePublicKey(pemKey string)(*rsa.PublicKey,error){
	if !strings.HasPrefix(pemKey,"-----BEGIN PUBLIC KEY-----"){
		pemKey = "-----BEGIN PUBLIC KEY-----\n" + pemKey+"\n-----END PUBLIC KEY-----"
	}
	block, _ := pem.Decode([]byte(pemKey))
	if block == nil{
		return nil,errors.New("pem public key not incorrect")
	}
	e,err := x509.ParsePKIXPublicKey(block.Bytes)
	if err == nil{
		return e.(*rsa.PublicKey),err
	}
	return nil,err
}

// RSA加密
func EncryptRSA( publicKey *rsa.PublicKey,data []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader,publicKey,data)
}

// RSA加密为BASE64
func EncryptRSAToBase64(publicKey *rsa.PublicKey,data []byte) (string, error){
	bytes,err := rsa.EncryptPKCS1v15(rand.Reader,publicKey,data)
	if err == nil{
		return base64.StdEncoding.EncodeToString(bytes),nil
	}
	return "",err
}

// 根据BASE64进行RSA加密
func DecryptRSAFromBase64(privateKey *rsa.PrivateKey,encryptData string)([]byte,error){
	bytes,err := base64.StdEncoding.DecodeString(encryptData)
	if err == nil {
		return rsa.DecryptPKCS1v15(rand.Reader, privateKey, bytes)
	}
	return nil,err
}