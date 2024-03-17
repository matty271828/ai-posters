package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/matty271828/ai-posters/internal/jobs"
)

type GenerateImageToImageHandlerRequest struct {
	Prompt     string `json:"prompt"`
	SeedPath   string `json:"seedPath,omitempty"`   // File path, used internally
	SeedBase64 string `json:"seedBase64,omitempty"` // Base64 encoded image
	Strength   string `json:"strength"`
	OutputPath string `json:"outputPath"`
}

func (a *API) GenerateImageToImageHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request body
	var request GenerateImageToImageHandlerRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Determine seed path
	seedPath, err := decodeImage(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the job function
	generatedPath, err := jobs.
		GenerateImageToImageJob(request.Prompt, seedPath, request.Strength)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the generated image as a response
	a.sendImageResponse(w, r, generatedPath[0])
}

func decodeImage(request GenerateImageToImageHandlerRequest) (string, error) {
	cleanBase64 := removeDataURLPrefix(request.SeedBase64)

	decodedImage, err := base64.StdEncoding.DecodeString(cleanBase64)
	if err != nil {
		return "", fmt.Errorf("invalid base64 data: %v", err)
	}

	// Ensure the "./assets/out/" directory exists
	outPath := "./assets/out"

	// Ensure the output directory exists
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		err := os.MkdirAll(outPath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	// Create a temporary file within the "./assets/out/" directory
	filePath := outPath + "/tmp.png"
	tmpFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer tmpFile.Close() // Ensure file is closed in all cases

	fmt.Printf("Length of decoded image: %d bytes\n", len(decodedImage))
	bytesWritten, err := tmpFile.Write(decodedImage)
	if err != nil {
		return "", fmt.Errorf("failed to write to temp file: %v", err)
	}
	fmt.Printf("Bytes written to temp file: %d bytes\n", bytesWritten)

	return filePath, nil
}

// removeDataURLPrefix removes the data URL prefix from a base64 encoded image string.
// This version uses a regular expression to match a broader range of image data URLs.
func removeDataURLPrefix(base64Data string) string {
	// Regular expression to match the data URL prefix and capture the base64 part
	re := regexp.MustCompile(`^data:image\/[a-zA-Z]+;base64,`)

	// Find the match and the index of the start of the base64 string
	indexes := re.FindStringSubmatchIndex(base64Data)
	if indexes != nil {
		// If a match is found, return the base64 part of the string
		return base64Data[indexes[1]:]
	}

	// If no match is found, return the original string as it might already be in base64
	return base64Data
}
