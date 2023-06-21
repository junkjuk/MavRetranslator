package MavSerial

import (
	"MavRetranslator/MavUtils"
	"fmt"
	"github.com/tarm/serial"
)

type MavSerial struct {
	serialPort *serial.Port
}

func NewMavSerial(name string, baud int) (*MavSerial, error) {
	ms := new(MavSerial)
	c := &serial.Config{Name: name, Baud: baud}
	s, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}
	ms.serialPort = s

	return ms, nil
}

func (ms MavSerial) Read() <-chan []byte {
	readChan := make(chan []byte, 10000)
	go func() {
		for {
			readCount := 0
			buff := make([]byte, MavUtils.MavMaxLength)
			n, _ := ms.serialPort.Read(buff[0:1])
			readCount += n
			if buff[0] != MavUtils.Mav2StartByte && buff[0] != MavUtils.Mav1StartByte {
				continue
			}
			isMavlink2 := buff[0] == MavUtils.Mav2StartByte
			var headerLenght int
			if isMavlink2 {
				headerLenght = int(MavUtils.Mav2HeaderLength)
			} else {
				headerLenght = int(MavUtils.Mav1HeaderLength)
			}
			n, _ = ms.serialPort.Read(buff[readCount : readCount+headerLenght])
			if n != headerLenght {
				continue
			}
			readCount += n
			lastLength := int(buff[1] + MavUtils.CrcLength)
			if isMavlink2 && (buff[2]&MavUtils.MavlinkIflagSigned) > 0 {
				lastLength += MavUtils.MavlinkSignatureBlockLen
			}

			n, _ = ms.serialPort.Read(buff[readCount : readCount+lastLength])
			if n != lastLength {
				continue
			}
			readCount += n

			readChan <- buff[0:readCount]
		}
	}()
	return readChan
}

func (ms MavSerial) Write(sendBuffer <-chan []byte) {
	go func() {
		for {
			_, err := ms.serialPort.Write(<-sendBuffer)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}()
}
