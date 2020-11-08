package crypto

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : des.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-08 07:26
 * description :
 * history :
 */


func DesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func DesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

// 3DES加密
func TripleDesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:8])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// 3DES解密
func TripleDesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:8])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}


func padding(src []byte, blockSize int) []byte {
	padNum := blockSize -len(src)%blockSize
	pad:=bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append(src,pad...)
}

func unPadding(src []byte) []byte {
	n:=len(src)
	unPadNum :=int(src[n-1])
	return src[:n-unPadNum]
}

func Encrypt3DES(src []byte,key []byte) []byte {
	block,_:=des.NewTripleDESCipher(key)
	src=padding(src,block.BlockSize())
	blockMode :=cipher.NewCBCEncrypter(block,key[:block.BlockSize()])
	blockMode.CryptBlocks(src,src)
	return src
}

func Encrypt3DESHex(src []byte,key []byte)string{
	return hex.EncodeToString(Encrypt3DES(src,key))
}

func Decrypt3DES(src []byte,key []byte) []byte {
	block,_:=des.NewTripleDESCipher(key)
	blockMode :=cipher.NewCBCDecrypter(block,key[:block.BlockSize()])
	blockMode.CryptBlocks(src,src)
	src= unPadding(src)
	return src
}