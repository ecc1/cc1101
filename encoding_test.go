package cc1100

import (
	"bytes"
	"testing"
)

func TestEncoding(t *testing.T) {
	cases := []struct {
		src []byte
		dst []byte
	}{
		{[]byte{},
			[]byte{}},
		{[]byte{0x00},
			[]byte{0x55, 0x50}},
		{[]byte{0x01},
			[]byte{0x57, 0x10}},
		{[]byte{0xFF},
			[]byte{0x71, 0xC0}},
		{[]byte{0x10, 0x20},
			[]byte{0xC5, 0x5C, 0x95}},
		{[]byte{0x33, 0x44, 0x55},
			[]byte{0x8E, 0x3D, 0x34, 0x96, 0x50}},
		{[]byte{0x87, 0x65, 0x43, 0x21},
			[]byte{0x69, 0x69, 0xA5, 0xD2, 0x3C, 0xB1}},
		{[]byte{0xA7, 0x12, 0x34, 0x56, 0x8D, 0x00, 0xA6},
			[]byte{0xA9, 0x6C, 0x72, 0x8F, 0x49, 0x66, 0x68, 0xD5, 0x55, 0xAA, 0x60}},
	}

	for _, c := range cases {
		result := Encode4b6b(c.src)
		if !bytes.Equal(result, c.dst) {
			t.Errorf("Encode4b6b(%X) == %X, want %X", c.src, result, c.dst)
		}
		result, err := Decode6b4b(c.dst)
		if err != nil {
			t.Errorf("Decode6b4b(%X) == %v, want %X", c.dst, err, c.src)
		} else if !bytes.Equal(result, c.src) {
			t.Errorf("Decode6b4b(%X) == %X, want %X", c.dst, result, c.src)
		}
	}
}
