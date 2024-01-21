package api

import (
	"encoding/json"
	"net/http"

	"github.com/matty271828/ai-posters/internal/jobs"
)

type GenerateImageHandlerRequest struct {
	Prompt     string `json:"prompt"`
	OutputPath string `json:"outputPath"`
}

func (a *API) GenerateImageHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request body
	var request GenerateImageHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Call the job function
	// TODO: Move the outpath out of the request to ensure is not exposed publically by the webui
	generatedPath, err := jobs.GenerateImageJob(request.Prompt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the generated image as a response
	a.sendImageResponse(w, r, generatedPath[0])
}
