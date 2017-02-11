package main

import (
	"log"
	"time"

	"github.com/definitelycarter/dht-go"
	"github.com/definitelycarter/wpi-go"
)

func main() {
	var pin int = 7
	wpi.WiringPiSetup()

	reader := dht.NewReader(pin)
	for i := 0; i < 10; i++ {
		n, err := reader.ReadTempurature()
		if err != nil {
			log.Println(err)
		}
		log.Println(n)
		time.Sleep(2 * time.Second)
	}
}
