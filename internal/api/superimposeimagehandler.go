package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/matty271828/ai-posters/internal/jobs"
)

const (
	smallFrame       = "assets/stock/smallframe.png"
	largeFrame       = "assets/stock/largeframe.png"
	largeFrameInRoom = "assets/stock/largeframeinroom.png"
)

type SuperImposeImageHandlerRequest struct {
	ImagePath string `json:"imagePath"`
}

func (a *API) SuperImposeImageHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request body
	var request SuperImposeImageHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Retrieve the frameSize parameter from URL
	frameSize := r.URL.Query().Get("frameSize")
	var framePath string

	// Select the correct frame path based on the frameSize
	switch frameSize {
	case "small":
		framePath = smallFrame
	case "large":
		framePath = largeFrame
	case "largeInRoom":
		framePath = largeFrameInRoom
	default:
		fmt.Println("Invalid frame size")
		http.Error(w, "Invalid frame size", http.StatusBadRequest)
		return
	}

	// Call the job function
	framedPath, err := jobs.SuperImposeImageJob(request.ImagePath, framePath)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the generated image as a response
	a.sendImageResponse(w, r, framedPath)
}
