package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ecc1/cc1100"
)

const (
	PumpDevice = 0xA7
	Wakeup     = 0x5D
	Ack        = 0x06

	numWakeups  = 100
	xmitDelay   = 40 * time.Millisecond
	recvTimeout = 3 * time.Second
)

func main() {
	dev, err := cc1100.Open()
	if err != nil {
		log.Fatal(err)
	}
	err = dev.Reset()
	if err != nil {
		log.Fatal(err)
	}
	err = dev.InitRF()
	if err != nil {
		log.Fatal(err)
	}

	dev.StartRadio()
	command := []byte{
		PumpDevice,
		cc1100.PumpID[0]<<4 | cc1100.PumpID[1],
		cc1100.PumpID[2]<<4 | cc1100.PumpID[3],
		cc1100.PumpID[4]<<4 | cc1100.PumpID[5],
		Wakeup,
		0,
	}
	packet := cc1100.EncodePacket(command)
	for i := 0; i < numWakeups; i++ {
		dev.OutgoingPackets() <- packet
		time.Sleep(xmitDelay)
	}
	tries := numWakeups
	for {
		dev.OutgoingPackets() <- packet
		tries++
		timeout := time.After(recvTimeout)
		var response cc1100.Packet
		select {
		case response = <-dev.IncomingPackets():
			fmt.Print("\n")
			break
		case <-timeout:
			fmt.Print(".")
			continue
		}
		data, err := dev.DecodePacket(response)
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}
		if len(data) == 7 && data[4] == Ack {
			rssi, err := dev.ReadRSSI()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("ACK after %d wakeups; RSSI = %d\n", tries, rssi)
			break
		}
		fmt.Printf("Unexpected response: % X\n", data)
	}
}
