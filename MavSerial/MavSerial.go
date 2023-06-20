package MavSerial

import "github.com/tarm/serial"

type MavSerial struct {
	serialPort *serial.Port
	sendBuffer chan []byte
}

func NewMavSerial(name string, baud int) (*MavSerial, error) {
	ms := new(MavSerial)
	ms.sendBuffer = make(chan []byte, 10000)
	c := &serial.Config{Name: name, Baud: baud}
	s, err := serial.OpenPort(c)
	if err == nil {
		return nil, err
	}
	ms.serialPort = s
	go ms.write()
	return ms, nil
}

func (ms MavSerial) Read() <-chan []byte {
	channel := make(chan []byte, 1000)
	return channel
}

func (ms MavSerial) write() error {
	for {
		_, err := ms.serialPort.Write(<-ms.sendBuffer)
		return err
	}
}

func (ms MavSerial) Send(packet []byte) {
	ms.sendBuffer <- packet
}
