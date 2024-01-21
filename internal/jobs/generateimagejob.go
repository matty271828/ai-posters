package jobs

import (
	"github.com/matty271828/ai-posters/internal/stabilityapi"
)

// GenerateImageJob is a job used to generate an image based on a prompt.
func GenerateImageJob(prompt string) ([]string, error) {
	// generate an image
	imagePaths, err := stabilityapi.GenerateImage(prompt)
	if err != nil {
		return nil, err
	}

	return imagePaths, nil
}
