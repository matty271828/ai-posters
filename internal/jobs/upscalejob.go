package jobs

import (
	"github.com/matty271828/ai-posters/internal/stabilityapi"
)

// GenerateUpscaleJob is a job used to upscale a pre existing image.
func GenerateUpscaleJob(seedPath, height, width string) (string, error) {
	// generate an image
	imagePath, err := stabilityapi.Upscale(seedPath, height, width)
	if err != nil {
		return "", err
	}

	return imagePath, nil
}
