/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2014-12-26 23:50
 * description :
 * history :
 */

package crypto

import (
	"fmt"
	"testing"
)

func TestEncode(t *testing.T){
	c := NewUnixCrypto("jarryix","12345")
	r:=c.Encode()
	s:= string(r)
	t.Log("encode:",string(s))
	ds,unix,err := c.Decode([]byte(s))
	t.Log("--dst:",string(ds),"unix:",unix,err)
}

func TestDecode(t *testing.T){
	//s := "25245e2640746f322e6e657480011596d091f58321984223c284faea4f14167725"

}

func Test_A(t *testing.T) {
	cyp := NewUnixCrypto("sonven", "3dsdgfdfgdfg")
	i := 2
	for {
		if i = i - 1; i == 0 {
			break
		}

		s := string(cyp.Encode())
		fmt.Println("str:",s)
		//cyp.Compare(s)

		//r,bytes,unix := cyp.Compare(s)
		//fmt.Println("dst:",string(bytes),time.Unix(unix,0).String())
		//fmt.Println("src:",string(cyp.GetBytes()))
		//fmt.Println("result:",r)
	}
}
