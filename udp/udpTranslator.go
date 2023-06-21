package udp

import (
	"fmt"
	"net"
)

type UdpTranslator struct {
	addresses  map[string]*net.UDPAddr
	connection *net.UDPConn
}

func NewUdpTranslator(address string) (*UdpTranslator, error) {
	udpTranslator := new(UdpTranslator)
	udpTranslator.addresses = make(map[string]*net.UDPAddr)
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	udpTranslator.connection, err = net.ListenUDP("udp4", udpAddr)
	if err != nil {
		return nil, err
	}
	return udpTranslator, nil
}

func (trans *UdpTranslator) ReadFromUdp() <-chan []byte {
	readChan := make(chan []byte, 10000)
	go func() {
		for {
			var buf [280]byte
			n, addr, err := trans.connection.ReadFromUDP(buf[0:])
			if trans.addresses[addr.String()] == nil {
				trans.addresses[addr.String()] = addr
			}
			if err != nil {
				fmt.Println(err)
				return
			}
			readChan <- buf[0:n]
		}
	}()
	return readChan
}

func (trans *UdpTranslator) WriteToUdp(write <-chan []byte) {
	go func() {
		for {
			buf := <-write
			for _, addr := range trans.addresses {
				_, err := trans.connection.WriteToUDP(buf, addr)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}()
}
