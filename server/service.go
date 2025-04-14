package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// Response structure from the ML API
type PredictionResponse struct {
	Prediction []interface{} `json:"prediction"`
	Error      string        `json:"error,omitempty"`
}

// SendToMLAPI sends image bytes to the Flask ML API and returns prediction
func SendToMLAPI(imageBytes []byte) (*PredictionResponse, error) {
	// Create a buffer to hold multipart form data
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Create form file field 'image'
	formFile, err := writer.CreateFormFile("image", "image.raw")
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}

	// Write image bytes into the form field
	_, err = io.Copy(formFile, bytes.NewReader(imageBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to write image to form: %v", err)
	}

	// Close the writer to finalize the multipart body
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %v", err)
	}

	// Make HTTP POST request
	resp, err := http.Post("http://localhost:5000/predict", writer.FormDataContentType(), &b)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Parse response
	var result PredictionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if result.Error != "" {
		return &result, fmt.Errorf("ML API error: %s", result.Error)
	}

	return &result, nil
}
