package algorithm

import (
	"testing"
	"time"
	"strconv"
)

func TestDJBHash(t *testing.T) {
	t.Log("=",DJBHash([]byte("jarry")))
	for i := 0; i <30; i++ {
		unix := strconv.Itoa(int(time.Now().UnixNano()))
		t.Log(unix,"=",DJBHash([]byte(unix)))
		time.Sleep(time.Nanosecond * 200)
	}

}
