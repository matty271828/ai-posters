package jobs

import (
	"github.com/matty271828/ai-posters/internal/imageprocessing"
	"github.com/matty271828/ai-posters/internal/stabilityapi"
)

// GenerateImageJob is a job used to generate an image based on a prompt.
// After generation, the image is superimposed onto a preview frame.
func GenerateImageJob(prompt, framePath, outputPath string) ([]string, string, error) {
	// generate an image
	imagePaths, err := stabilityapi.GenerateImage(prompt)
	if err != nil {
		return nil, "", err
	}

	// superimpose the image onto a poster frame
	framePreviewPath, err := imageprocessing.Frame(
		framePath, imagePaths[0], outputPath,
	)
	if err != nil {
		return nil, "", err
	}

	return imagePaths, framePreviewPath, nil
}
