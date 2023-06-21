package main

import (
	"MavRetranslator/MavSerial"
	"MavRetranslator/signal"
	"MavRetranslator/udp"
)

func main() {
	udpTranslator, _ := udp.NewUdpTranslator(":14567")
	serial, _ := MavSerial.NewMavSerial("COM8", 1500000)

	MessagesFromSerialPort := serial.Read()
	udpTranslator.Write(MessagesFromSerialPort)

	messagesFromUdp := udpTranslator.Read()
	serial.Write(messagesFromUdp)

	signal.WaitForTerminationSignal()
}
