package dialect

import (
	"log"
	"testing"
)

func TestGetFieldLength(t *testing.T){
	l := getTypeLen("varchar(100)")
	l1 := getTypeLen("decimal(10,2)")
	log.Println(l,l1)
}