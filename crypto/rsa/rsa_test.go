package rsa

import (
	"fmt"
	"testing"
)


func TestRsaToken(t *testing.T) {
	_,privateKey,err := GenRsaKeys(2048)
	if err != nil{
		t.Error(err)
	}
	fmt.Println(privateKey)
}
