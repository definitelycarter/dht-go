package dht

import (
	"errors"
	"time"

	"github.com/definitelycarter/wpi-go"
)

type DHT11Reader struct {
	pin      int
	temp     int
	humidity int
	last     time.Time
}

func NewReader(pin int) DHT11Reader {
	return DHT11Reader{pin: pin}
}

func (r DHT11Reader) ReadTempurature() (int, error) {
	if time.Now().Sub(r.last) < time.Duration(2*time.Second) {
		return r.temp, nil
	}
	defer func() { r.last = time.Now() }()

	if err := r.readData(); err != nil {
		return 0, err
	}

	return r.temp, nil
}

func (r *DHT11Reader) readData() error {
	cycles := [80]uint32{}
	data := [5]uint8{}

	wpi.PinMode(r.pin, wpi.OUTPUT)
	time.Sleep(20 * time.Millisecond)

	wpi.DigitalWrite(r.pin, wpi.LOW)
	time.Sleep(20 * time.Millisecond)

	wpi.DigitalWrite(r.pin, wpi.HIGH)
	time.Sleep(40 * time.Microsecond)

	wpi.PinMode(r.pin, wpi.INPUT)
	time.Sleep(10 * time.Microsecond)

	if pulse := expectPulse(r.pin, wpi.LOW); pulse == 0 {
		return errors.New("Timeout waiting for low pulse")
	}
	if pulse := expectPulse(r.pin, wpi.HIGH); pulse == 0 {
		return errors.New("Timeout waiting for high pulse")
	}

	for i := 0; i < 80; i += 2 {
		// always near 50 microseconds
		cycles[i] = expectPulse(r.pin, wpi.LOW)
		// either 70 microseconds for 1
		// or 28 microseconds for 0
		cycles[i+1] = expectPulse(r.pin, wpi.HIGH)
	}

	var threshold uint32
	for i := 2; i < 80; i += 2 {
		threshold += cycles[i]
	}
	threshold /= 40

	for i := 0; i < 40; i++ {
		var (
			low = cycles[2*i]
			//high = cycles[2*i+1]
		)
		data[i/8] <<= 1
		if low > threshold {
			data[i/8] |= 1
		}
	}

	if data[4] == ((data[0] + data[1] + data[2] + data[3]) & 0xFF) {
		// var f float64 = float64(data[2])*1.8 + 32
		r.temp = int(data[2])
		r.humidity = int(data[0])
		return nil
	}
	return errors.New("checksum failed")
}

func expectPulse(pin int, expected int) uint32 {
	var count uint32
	n := 0
	for ; count < 3200; count++ {
		if n = wpi.DigitalRead(pin); n == expected {
			return count + 1
		}
	}
	return 0
}
