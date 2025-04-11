package cmd

import (
	"github.com/stianeikeland/go-rpio"
)

const (
	FanPin   = 22 // Example GPIO pin for fan
	LightPin = 23 // Example GPIO pin for light
)

// InitDevices sets up relay pins
func InitDevices() {
	rpio.PinMode(rpio.Pin(FanPin), rpio.Output)
	rpio.PinMode(rpio.Pin(LightPin), rpio.Output)
}

// TurnOnFan turns on the fan
func TurnOnFan() {
	rpio.Pin(FanPin).High()
}

// TurnOffFan turns off the fan
func TurnOffFan() {
	rpio.Pin(FanPin).Low()
}

// TurnOnLight turns on the light
func TurnOnLight() {
	rpio.Pin(LightPin).High()
}

// TurnOffLight turns off the light
func TurnOffLight() {
	rpio.Pin(LightPin).Low()
}
