package cc1101

// Configuration for Raspberry Pi Zero W.

const (
	spiDevice    = "/dev/spidev0.1"
	spiSpeed     = 6000000 // Hz
	interruptPin = 24      // GPIO for receive interrupts
)
