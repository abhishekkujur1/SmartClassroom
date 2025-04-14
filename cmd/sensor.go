package cmd

import (
	"log"
	"time"

	"github.com/stianeikeland/go-rpio"
)

// GPIO pin numbers for Raspberry Pi (BCM)
const (
	TriggerPin        = 17    // GPIO17 (BCM)
	EchoPin           = 27    // GPIO27 (BCM)
	DistanceThreshold = 200.0 // 2 meters in cm
)

// InitSensor initializes the GPIO pins
func InitSensor() error {
	if err := rpio.Open(); err != nil {
		return err
	}
	rpio.PinMode(rpio.Pin(TriggerPin), rpio.Output)
	rpio.PinMode(rpio.Pin(EchoPin), rpio.Input)
	return nil
}

// CloseSensor releases the GPIO resources
func CloseSensor() {
	rpio.Close()
}

// MeasureDistance measures the distance using an ultrasonic sensor
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

func StartCapture() {
	// Initialize the sensor
	if err := InitSensor(); err != nil {
		log.Fatalf("Failed to initialize sensor: %v", err)
	}
	defer CloseSensor()

	// Initialize the camera
	if err := InitCamera(); err != nil {
		log.Fatalf("Failed to initialize camera: %v", err)
	}
	defer CloseCamera()

	// Continuously measure distance
	for {
		distance := MeasureDistance()
		if distance < 0 {
			log.Println("Invalid distance measurement")
			continue
		}

		log.Printf("Distance: %.2f cm", distance)

		// Call imageCapture when distance is less than 2 meters
		if distance < DistanceThreshold {
			if err := ImageCapture(); err != nil {
				log.Printf("Failed to capture image: %v", err)
			} else {
				log.Println("Image captured successfully")
			}
		}

		// Small delay to avoid overwhelming the sensor
		time.Sleep(1 * time.Second)
	}
}
