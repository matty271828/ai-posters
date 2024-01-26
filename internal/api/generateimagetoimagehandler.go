package api

import (
	"encoding/json"
	"net/http"

	"github.com/matty271828/ai-posters/internal/jobs"
)

type GenerateImageToImageHandlerRequest struct {
	Prompt     string `json:"prompt"`
	SeedPath   string `json:"seedPath"`
	OutputPath string `json:"outputPath"`
}

func (a *API) GenerateImageToImageHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request body
	var request GenerateImageToImageHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Call the job function
	// TODO: Move the outpath out of the request to ensure is not exposed publically by the webui
	generatedPath, err := jobs.GenerateImageToImageJob(request.Prompt, request.SeedPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the generated image as a response
	a.sendImageResponse(w, r, generatedPath[0])
}
