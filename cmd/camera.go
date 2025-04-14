package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

// OV7670 I2C address and pin assignments (BCM numbering)
const (
	OV7670Addr = 0x21 // OV7670 I2C address
	XCLKPin    = 18   // GPIO18 for camera clock (PWM)
	VSYNCPin   = 23   // GPIO23 for vertical sync
	HREFPin    = 24   // GPIO24 for horizontal sync
	PCLKPin    = 25   // GPIO25 for pixel clock
	DataPin0   = 5    // GPIO5 for D0
	DataPin7   = 12   // GPIO12 for D7 (D0-D7 for 8-bit data)
)

// Camera represents the OV7670 camera module
type Camera struct {
	i2cDev *i2c.Dev
}

// global camera instance
var cam *Camera

// InitCamera initializes the OV7670 camera
func InitCamera() error {
	// Initialize periph.io host
	if _, err := host.Init(); err != nil {
		return fmt.Errorf("failed to initialize periph host: %v", err)
	}

	// Open I2C bus
	bus, err := i2creg.Open("")
	if err != nil {
		return fmt.Errorf("failed to open I2C bus: %v", err)
	}

	// Initialize camera struct
	cam = &Camera{
		i2cDev: &i2c.Dev{Addr: OV7670Addr, Bus: bus},
	}

	// Initialize GPIO
	if err := rpio.Open(); err != nil {
		return fmt.Errorf("failed to open GPIO: %v", err)
	}

	// Configure GPIO pins
	rpio.PinMode(rpio.Pin(VSYNCPin), rpio.Input)
	rpio.PinMode(rpio.Pin(HREFPin), rpio.Input)
	rpio.PinMode(rpio.Pin(PCLKPin), rpio.Input)
	for pin := DataPin0; pin <= DataPin7; pin++ {
		rpio.PinMode(rpio.Pin(pin), rpio.Input)
	}

	// Configure XCLK (simulated clock using GPIO toggle)
	go generateXCLK()

	// Configure OV7670 registers via I2C
	if err := configureOV7670(); err != nil {
		return fmt.Errorf("failed to configure OV7670: %v", err)
	}

	return nil
}

// CloseCamera releases camera resources
func CloseCamera() {
	if cam != nil && cam.i2cDev != nil {
		if bus, ok := cam.i2cDev.Bus.(i2c.BusCloser); ok {
			bus.Close()
		}
	}
	rpio.Close()
}

// generateXCLK simulates an 8 MHz clock on XCLKPin
func generateXCLK() {
	pin := rpio.Pin(XCLKPin)
	rpio.PinMode(pin, rpio.Output)
	for {
		pin.High()
		time.Sleep(62 * time.Nanosecond) // ~8 MHz (125ns period)
		pin.Low()
		time.Sleep(62 * time.Nanosecond)
	}
}

// configureOV7670 sets up the OV7670 registers for QVGA YUV422
func configureOV7670() error {
	// Basic configuration for QVGA (320x240), YUV422
	regs := [][2]byte{
		{0x12, 0x80}, // COM7: Reset all registers
		{0x11, 0x0A}, // CLKRC: Prescaler for 8 MHz input
		{0x12, 0x04}, // COM7: QVGA
		{0x0C, 0x00}, // COM3: Disable scaling
		{0x3E, 0x00}, // COM14: No scaling
		{0x40, 0xD0}, // COM15: YUV, full output range
		{0x3A, 0x04}, // TSLB: YUYV order
		{0x8C, 0x00}, // RGB444: Disable RGB444
	}

	for _, reg := range regs {
		// Handle both return values from Write
		_, err := cam.i2cDev.Write([]byte{reg[0], reg[1]})
		if err != nil {
			return fmt.Errorf("failed to write register 0x%02X: %v", reg[0], err)
		}
		time.Sleep(10 * time.Millisecond) // Delay for register write to settle
	}
	return nil
}

// ImageCapture captures an image from the OV7670 and saves it to a file
func ImageCapture() error {
	// QVGA resolution (320x240, YUV422: 2 bytes per pixel)
	const (
		width  = 320
		height = 240
	)

	// Buffer for image data (YUV422: 2 bytes per pixel)
	imgData := make([]byte, width*height*2)

	// Wait for VSYNC to go low (start of frame)
	vsync := rpio.Pin(VSYNCPin)
	for vsync.Read() == rpio.High {
	}
	for vsync.Read() == rpio.Low {
	}

	// Read frame
	idx := 0
	for row := 0; row < height; row++ {
		// Wait for HREF to go high (start of row)
		href := rpio.Pin(HREFPin)
		for href.Read() == rpio.Low {
		}

		// Read one row
		for col := 0; col < width*2; col++ { // 2 bytes per pixel
			// Wait for PCLK to go high
			pclk := rpio.Pin(PCLKPin)
			for pclk.Read() == rpio.Low {
			}

			// Read 8-bit data (D0-D7)
			var data byte
			for i := 0; i < 8; i++ {
				if rpio.Pin(DataPin0+uint(i)).Read() == rpio.High {
					data |= 1 << i
				}
			}
			if idx < len(imgData) {
				imgData[idx] = data
				idx++
			}

			// Wait for PCLK to go low
			for pclk.Read() == rpio.High {
			}
		}

		// Wait for HREF to go low (end of row)
		for href.Read() == rpio.High {
		}
	}

	// Save image to file
	filename := fmt.Sprintf("capture_%s.raw", time.Now().Format("20060102_150405"))
	if err := os.WriteFile(filename, imgData, 0o644); err != nil {
		return fmt.Errorf("failed to save image: %v", err)
	}

	log.Printf("Image saved as %s", filename)
	return nil
}
