package main

import (
	"MavRetranslator/signal"
	"fmt"
	"time"
)

func main2() {
	data := startSerial()
	go func() {
		for {
			fmt.Println(string(<-data))
		}
	}()
	signal.WaitForTerminationSignal()
}

func startSerial() chan []byte {

	outCh := make(chan []byte)
	go func() {
		i := 0
		for {
			message := fmt.Sprintf("hello numb %d", i)
			outCh <- []byte(message)
			outCh <- []byte("1234")
			i++
			time.Sleep(1 * time.Second)
		}
	}()

	return outCh
}
