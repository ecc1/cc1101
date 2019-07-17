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

// SPIDevice returns the pathname of the radio's SPI device.
func (hwFlavor) SPIDevice() string {
	return spiDevice
}

// Speed returns the radio's SPI speed.
func (hwFlavor) Speed() int {
	return spiSpeed
}

// CustomCS returns the GPIO pin number to use as a custom chip-select for the radio.
func (hwFlavor) CustomCS() int {
	return customCS
}

// InterruptPin returns the GPIO pin number to use for receive interrupts.
func (hwFlavor) InterruptPin() int {
	return interruptPin
}

// ReadSingleAddress returns the encoding of an address for SPI read operations.
func (hwFlavor) ReadSingleAddress(addr byte) byte {
	return READ_MODE | addr
}

// ReadBurstAddress returns the encoding of an address for SPI burst-read operations.
func (hwFlavor) ReadBurstAddress(addr byte) byte {
	reg := addr & 0x3F
	if 0x30 <= reg && reg <= 0x3D {
		log.Panicf("no burst access for CC1101 status register %02X", reg)
	}
	return READ_MODE | BURST_MODE | addr
}

// WriteSingleAddress returns the (identity) encoding of an address for SPI write operations.
func (hwFlavor) WriteSingleAddress(addr byte) byte {
	return addr
}

// WriteBurstAddress returns the encoding of an address for SPI burst-write operations.
func (hwFlavor) WriteBurstAddress(addr byte) byte {
	return BURST_MODE | addr
}

// Radio represents an open radio device.
type Radio struct {
	hw            *radio.Hardware
	receiveBuffer bytes.Buffer
	err           error
}

// Open opens the radio device.
func Open() *Radio {
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

// Close closes the radio device.
func (r *Radio) Close() {
	r.changeState(SIDLE, STATE_IDLE)
	r.hw.Close()
}

// Version returns the radio's hardware version.
func (r *Radio) Version() uint16 {
	p := r.hw.ReadRegister(PARTNUM)
	v := r.hw.ReadRegister(VERSION)
	return uint16(p)<<8 | uint16(v)
}

// Name returns the radio's name.
func (*Radio) Name() string {
	return "CC1101"
}

// Device returns the pathname of the radio's device.
func (*Radio) Device() string {
	return spiDevice
}

// Strobe writes the given command to the radio.
func (r *Radio) Strobe(cmd byte) byte {
	if verbose && cmd != SNOP {
		log.Printf("issuing %s command", strobeName(cmd))
	}
	buf := []byte{cmd}
	r.err = r.hw.SPIDevice().Transfer(buf)
	return buf[0]
}

// Reset resets the radio device.
func (r *Radio) Reset() {
	r.Strobe(SRES)
}

// Init initializes the radio device.
func (r *Radio) Init(frequency uint32) {
	r.Reset()
	r.InitRF(frequency)
}

// Error returns the error state of the radio device.
func (r *Radio) Error() error {
	err := r.hw.Error()
	if err != nil {
		return err
	}
	return r.err
}

// SetError sets the error state of the radio device.
func (r *Radio) SetError(err error) {
	r.hw.SetError(err)
	r.err = err
}

// Hardware returns the radio's hardware information.
func (r *Radio) Hardware() *radio.Hardware {
	return r.hw
}
