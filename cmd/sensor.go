// sensor.go
package main

import (
	"log"
	"time"

	"github.com/stianeikeland/go-rpio"
)

const (
	TriggerPin         = 17 // GPIO17 (BCM pin)
	EchoPin            = 27 // GPIO27 (BCM pin)
	DistanceThreshold  = 200.0 // 2 meters in cm
)

func InitSensor() error {
	if err := rpio.Open(); err != nil {
		return err
	}
	rpio.PinMode(rpio.Pin(TriggerPin), rpio.Output)
	rpio.PinMode(rpio.Pin(EchoPin), rpio.Input)
	return nil
}

func CloseSensor() {
	rpio.Close()
}

func MeasureDistance() float64 {
	trigger := rpio.Pin(TriggerPin)
	echo := rpio.Pin(EchoPin)

	trigger.Low()
	time.Sleep(2 * time.Microsecond)
	trigger.High()
	time.Sleep(10 * time.Microsecond)
	trigger.Low()

	// Wait for echo to go high
	start := time.Now()
	for echo.Read() == rpio.Low {
		if time.Since(start) > time.Second {
			log.Println("Timeout waiting for echo high")
			return -1
		}
	}
	start = time.Now()

	// Wait for echo to go low
	for echo.Read() == rpio.High {
		if time.Since(start) > time.Second {
			log.Println("Timeout waiting for echo low")
			return -1
		}
	}
	duration := time.Since(start).Seconds()
	distance := duration * 17150 // cm

	return distance
}
