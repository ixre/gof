package api

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/ixre/gof/crypto"
	"testing"
	"time"
)

func TestGenRsaPair(t *testing.T) {
	pubKey, privateKey, err := crypto.GenRsaKeys(2048)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("publicKey=", pubKey)
	fmt.Println("privateKey=", privateKey)
}

func TestJwtToken(t *testing.T) {
	_, privateKey, err := crypto.GenRsaKeys(2048)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(privateKey)

	// 生成token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   "jwt",
		IssuedAt:  time.Now().Unix(),
		Audience:  "jarrysix",
		ExpiresAt: time.Now().Unix() + 3600,
		Issuer:    "go-jwt",
	})
	token, err := claims.SignedString([]byte(privateKey))
	fmt.Println(token, err)
	// 转换token
	dstClaims := jwt.StandardClaims{Issuer: "sss"} // 可以用实现了Claim接口的自定义结构
	tk, err := jwt.ParseWithClaims(token, &dstClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte(privateKey), nil
	})
	// 判断是否有效
	if !tk.Valid {
		ve, _ := err.(*jwt.ValidationError)
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			println("---token is expired")
		} else {
			println("---token is invalid")
		}
	}
	println(tk.Valid, fmt.Sprintf("%#v", dstClaims))
}
