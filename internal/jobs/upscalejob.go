package jobs

import (
	"log"

	"github.com/matty271828/ai-posters/internal/stabilityapi"
)

// GenerateUpscaleJob is a job used to upscale a pre existing image.
func GenerateUpscaleJob(seedPath, height string) (string, error) {
	log.Printf("Starting upscale job with SeedPath: %s, Height: %s\n", seedPath, height)

	// generate an image
	imagePath, err := stabilityapi.Upscale(seedPath, height)
	if err != nil {
		log.Printf("Failed to upscale image with SeedPath: %s, Height: %s. Error: %v\n", seedPath, height, err)
		return "", err
	}

	log.Printf("Successfully upscaled image. SeedPath: %s, Height: %s. Resulting Image Path: %s\n", seedPath, height, imagePath)

	return imagePath, nil
}
