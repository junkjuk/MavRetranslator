package main

import (
	"MavRetranslator/MavSerial"
	"MavRetranslator/signal"
	"MavRetranslator/udp"
)

func main() {
	udp, _ := udp.NewUdpTranslator(":14567")
	ser, _ := MavSerial.NewMavSerial("COM8", 1500000)
	readBuffer := ser.Read()
	send := ser.Write()
	udpBuffer := udp.ReadFromUdp()

	go func() {
		for {
			send <- <-udpBuffer
		}
	}()
	udp.WriteToUdp(readBuffer)

	signal.WaitForTerminationSignal()
}
