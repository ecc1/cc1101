package cc1101

import (
	"bytes"
	"testing"
)

func TestFrequency(t *testing.T) {
	cases := []struct {
		f       uint32
		b       []byte
		fApprox uint32 // 0 => equal to f
	}{
		{315000000, []byte{0x0D, 0x20, 0x00}, 0},
		{915000000, []byte{0x26, 0x20, 0x00}, 0},
		// some that can't be represented exactly:
		{434000000, []byte{0x12, 0x15, 0x55}, 433999877},
		{868000000, []byte{0x24, 0x2A, 0xAB}, 868000122},
		{916300000, []byte{0x26, 0x2D, 0xDE}, 916300048},
		{916600000, []byte{0x26, 0x31, 0x11}, 916599975},
	}
	for _, c := range cases {
		b := frequencyToRegisters(c.f)
		if !bytes.Equal(b, c.b) {
			t.Errorf("frequencyToRegisters(%d) == % X, want % X", c.f, b, c.b)
		}
		f := registersToFrequency(c.b)
		if c.fApprox == 0 {
			if f != c.f {
				t.Errorf("registersToFrequency(% X) == %d, want %d", c.b, f, c.f)
			}
		} else {
			if f != c.fApprox {
				t.Errorf("registersToFrequency(% X) == %d, want %d", c.b, f, c.fApprox)
			}
		}
	}
}
