package cc1101

import (
	"bytes"
	"log"

	"github.com/ecc1/radio"
)

const (
	hwVersion = 0x0014
)

type hwFlavor struct{}

func (f hwFlavor) Name() string {
	return "CC1101"
}

func (f hwFlavor) SPIDevice() string {
	return spiDevice
}

func (f hwFlavor) Speed() int {
	return spiSpeed
}

func (f hwFlavor) CustomCS() int {
	return customCS
}

func (f hwFlavor) InterruptPin() int {
	return interruptPin
}

func (f hwFlavor) ReadSingleAddress(addr byte) byte {
	return READ_MODE | addr
}

func (f hwFlavor) ReadBurstAddress(addr byte) byte {
	reg := addr & 0x3F
	if 0x30 <= reg && reg <= 0x3D {
		log.Panicf("no burst access for CC1101 status register %02X", reg)
	}
	return READ_MODE | BURST_MODE | addr
}

func (f hwFlavor) WriteSingleAddress(addr byte) byte {
	return addr
}

func (f hwFlavor) WriteBurstAddress(addr byte) byte {
	return BURST_MODE | addr
}

type Radio struct {
	hw            *radio.Hardware
	receiveBuffer bytes.Buffer
	stats         radio.Statistics
	err           error
}

func Open() radio.Interface {
	r := &Radio{hw: radio.Open(hwFlavor{})}
	v := r.Version()
	if r.Error() != nil {
		return r
	}
	if v != hwVersion {
		r.hw.Close()
		r.SetError(radio.HardwareVersionError{Actual: v, Expected: hwVersion})
		return r
	}
	return r
}

func (r *Radio) Close() {
	r.changeState(SIDLE, STATE_IDLE)
	r.hw.Close()
}

func (r *Radio) Version() uint16 {
	p := r.hw.ReadRegister(PARTNUM)
	v := r.hw.ReadRegister(VERSION)
	return uint16(p)<<8 | uint16(v)
}

func (r *Radio) Strobe(cmd byte) byte {
	if verbose && cmd != SNOP {
		log.Printf("issuing %s command", strobeName(cmd))
	}
	buf := []byte{cmd}
	r.err = r.hw.SPIDevice().Transfer(buf)
	return buf[0]
}

func (r *Radio) Reset() {
	r.Strobe(SRES)
}

func (r *Radio) Init(frequency uint32) {
	r.Reset()
	r.InitRF(frequency)
}

func (r *Radio) Statistics() radio.Statistics {
	return r.stats
}

func (r *Radio) Error() error {
	err := r.hw.Error()
	if err != nil {
		return err
	}
	return r.err
}

func (r *Radio) SetError(err error) {
	r.hw.SetError(err)
	r.err = err
}

func (r *Radio) Hardware() *radio.Hardware {
	return r.hw
}
