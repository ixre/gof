package crypto

import (
	"fmt"
	"testing"
)

func TestAesCBC(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef") // 16字节、24字节或32字节的密钥
	iv := []byte("1234567890abcdef")                  // 初始向量，长度与块大小相同（AES通常为16字节）
	plaintext := []byte("这是一段要加密的明文")

	ciphertext, err := AESEncrypt(plaintext, key, iv)
	if err != nil {
		fmt.Println("加密出错:", err)
		return
	}
	fmt.Printf("加密后的密文: %x\n", ciphertext)

	decryptedText, err := AESDecrypt(ciphertext, key, iv)
	if err != nil {
		fmt.Println("解密出错:", err)
		return
	}
	fmt.Printf("解密后的明文: %s\n", decryptedText)
}
