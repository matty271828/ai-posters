package jobs

import (
	"github.com/matty271828/ai-posters/internal/stabilityapi"
)

// GenerateImageToImageJob is a job used to generate an image based on a prompt and an image.
func GenerateImageToImageJob(prompt, seedPath, strength string) ([]string, error) {
	// generate an image
	imagePaths, err := stabilityapi.GenerateImageToImage(prompt, seedPath, strength)
	if err != nil {
		return nil, err
	}

	return imagePaths, nil
}
