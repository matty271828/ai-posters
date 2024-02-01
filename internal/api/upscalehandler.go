package api

import (
	"encoding/json"
	"net/http"

	"github.com/matty271828/ai-posters/internal/jobs"
)

type GenerateUpscaleHandlerRequest struct {
	SeedPath string `json:"seedPath"`
}

func (a *API) GenerateUpscaleHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request body
	var request GenerateUpscaleHandlerRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Retrieve the height and widths parameter from URL.
	// Only one can be used.
	height := r.URL.Query().Get("height")
	width := r.URL.Query().Get("width")

	// Call the job function
	generatedPath, err := jobs.
		GenerateUpscaleJob(request.SeedPath, height, width)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the generated image as a response
	a.sendImageResponse(w, r, generatedPath)
}
