package api

import (
	"io"
	"net/http"
	"os"
)

type API struct {
	// TODO: Add dependencies as needed e.g., database, config, etc.
}

func NewAPI() *API {
	return &API{
		// TODO: Initialize dependencies
	}
}

// sendImageResponse reads an image file from a given path and writes it to the HTTP response.
func (a *API) sendImageResponse(w http.ResponseWriter, r *http.Request, imagePath string) {
	// Open the image file
	file, err := os.Open(imagePath)
	if err != nil {
		http.Error(w, "Unable to open image file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Detect the content type of the image
	buffer := make([]byte, 512) // You only need the first 512 bytes to sniff the content type
	if _, err := file.Read(buffer); err != nil {
		http.Error(w, "Failed to read image file", http.StatusInternalServerError)
		return
	}
	contentType := http.DetectContentType(buffer)

	// Reset the offset back to the start of the file
	if _, err := file.Seek(0, 0); err != nil {
		http.Error(w, "Failed to read image file", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header
	w.Header().Set("Content-Type", contentType)

	// Write the image file to the response
	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, "Failed to send image", http.StatusInternalServerError)
		return
	}
}
