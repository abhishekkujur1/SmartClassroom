package main

import (
	"fmt"
	"time"

	"gocv.io/x/gocv"
)

func main() {
	// Initialize ultrasonic sensor
	if err := InitSensor(); err != nil {
		fmt.Println("❌ Failed to initialize sensor:", err)
		return
	}
	defer CloseSensor()

	// Initialize fan and light control (relays)
	InitDevices()

	// Initialize camera
	cam, err := InitCamera(0)
	if err != nil {
		fmt.Println("❌ Failed to open camera:", err)
		return
	}
	defer cam.Close()

	window := gocv.NewWindow("Webcam")
	defer window.Close()

	fmt.Println("📷 Smart Classroom System Running...")
	fmt.Println("Press ESC to exit")

	for {
		// Measure distance using ultrasonic sensor
		distance := MeasureDistance()
		fmt.Printf("📏 Distance: %.2f cm\n", distance)

		// If a person is detected within threshold (2 meters)
		if distance > 0 && distance < DistanceThreshold {
			TurnOnFan()
			TurnOnLight()
			fmt.Println("✅ Person Detected: Fan & Light ON")

			// Capture image
			img, err := cam.CaptureImage()
			if err != nil {
				fmt.Println("❌ Error capturing image:", err)
				break
			}
			window.IMShow(img)

			// Static name for now - replace with face recognition later
			err = MarkAttendance("Student1")
			if err != nil {
				fmt.Println("❌ Attendance error:", err)
			}
		} else {
			TurnOffFan()
			TurnOffLight()
			fmt.Println("❌ No one detected: Fan & Light OFF")
		}

		// Break on ESC key
		if window.WaitKey(1) == 27 {
			break
		}

		time.Sleep(2 * time.Second) // Delay between checks
	}
}
