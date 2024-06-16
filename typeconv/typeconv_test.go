package typeconv

import "testing"

func TestMustInt(t *testing.T) {
	v := MustInt(nil)
	t.Log(v)
}
