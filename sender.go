package main

import (
	"MavRetranslator/MavSerial"
	"MavRetranslator/signal"
	"fmt"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", ":14567")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ser, _ := MavSerial.NewMavSerial("COM8", 1500000)
	readBuffer := ser.Read()
	send := ser.Write()
	addresses := make(map[string]*net.UDPAddr)
	go func() {
		for {
			var buf [280]byte
			n, addr, err := conn.ReadFromUDP(buf[0:])
			if addresses[addr.String()] == nil {
				addresses[addr.String()] = addr
			}
			if err != nil {
				fmt.Println(err)
				return
			}
			send <- buf[0:n]
		}
	}()

	go func() {
		for {
			buf := <-readBuffer
			for _, addr := range addresses {
				_, err := conn.WriteToUDP(buf, addr)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}()
	signal.WaitForTerminationSignal()
}
