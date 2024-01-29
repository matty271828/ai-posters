package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

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

// removeDataURLPrefix removes data URL prefix from base64 encoded images.
func removeDataURLPrefix(base64Data string) string {
	// Define a list of possible prefixes
	prefixes := []string{
		"data:image/png;base64,",
		"data:image/jpeg;base64,",
		// Add other image formats as needed
	}

	// Iterate over the prefixes and remove if found
	for _, prefix := range prefixes {
		if strings.HasPrefix(base64Data, prefix) {
			return strings.TrimPrefix(base64Data, prefix)
		}
	}

	// Return the original string if no prefix is found
	return base64Data
}
