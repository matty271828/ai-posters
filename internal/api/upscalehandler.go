package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/matty271828/ai-posters/internal/jobs"
)

type GenerateUpscaleHandlerRequest struct {
	SeedPath   string `json:"seedPath,omitempty"`   // Optional: For when a server path is used
	SeedBase64 string `json:"seedBase64,omitempty"` // Optional: For base64 encoded images
}

func (a *API) GenerateUpscaleHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for image upscaling")

	// Decode the JSON request body
	var request GenerateUpscaleHandlerRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Determine how to process the image based on the request content
	var imagePath string
	if request.SeedBase64 != "" {
		// Process the base64 encoded image
		imagePath, err = decodeBase64Image(request.SeedBase64)
		if err != nil {
			log.Printf("Failed to decode and save base64 image: %v", err)
			http.Error(w, "Failed to process base64 image", http.StatusInternalServerError)
			return
		}
	} else if request.SeedPath != "" {
		// Use the provided server path directly
		imagePath = request.SeedPath
	} else {
		http.Error(w, "No image data provided", http.StatusBadRequest)
		return
	}

	// Retrieve the height and width parameters from URL
	height := r.URL.Query().Get("height")
	log.Printf("Received height: %s for upscaling", height)

	// Call the job function
	log.Printf("Calling UpscaleJob with Image Path: %s, Height: %s", imagePath, height)
	generatedPath, err := jobs.GenerateUpscaleJob(imagePath, height)
	if err != nil {
		log.Printf("Failed to upscale image: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Upscaled image generated at path: %s", generatedPath)
	// Send the generated image as a response
	a.sendImageResponse(w, r, generatedPath)
}

// decodeBase64Image takes a base64 encoded string, decodes it, and saves it as a temporary file.
// It returns the path to the saved file for further processing.
func decodeBase64Image(base64Data string) (string, error) {
	cleanBase64 := removeDataURLPrefix(base64Data)

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

	// Generate a unique file name for the new image to prevent overwriting existing files
	// Here using a simple approach for demonstration; consider a more robust method for production
	filePath := fmt.Sprintf("%s/%d.png", outPath, os.Getpid())

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
