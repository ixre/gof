package rsa

import (
	"fmt"
	"github.com/ixre/gof/crypto/rsa"
	"testing"
)

func TestRsaToken(t *testing.T) {
	publicKey, privateKey, err := rsa.GenRsaKeys(2048)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("publicKey=", publicKey)
	fmt.Println("privateKey=", privateKey)
}
