package cc1101

import (
	"testing"
	"unsafe"
)

func TestRFConfiguration(t *testing.T) {
	have := int(unsafe.Sizeof(RFConfiguration{}))
	want := TEST0 - IOCFG2 + 1
	if have != want {
		t.Errorf("Sizeof(RFConfiguration) == %d, want %d", have, want)
	}
}
