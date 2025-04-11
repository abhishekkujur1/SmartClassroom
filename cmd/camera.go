// camera.go
package main

import (
	"fmt"

	"gocv.io/x/gocv"
)

type Camera struct {
	webcam *gocv.VideoCapture
	frame  gocv.Mat
}

func InitCamera(deviceID int) (*Camera, error) {
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		return nil, fmt.Errorf("error opening camera: %w", err)
	}

	return &Camera{
		webcam: webcam,
		frame:  gocv.NewMat(),
	}, nil
}

func (c *Camera) CaptureImage() (gocv.Mat, error) {
	if ok := c.webcam.Read(&c.frame); !ok || c.frame.Empty() {
		return gocv.NewMat(), fmt.Errorf("cannot capture image")
	}
	return c.frame.Clone(), nil
}

func (c *Camera) Close() {
	c.frame.Close()
	c.webcam.Close()
}
