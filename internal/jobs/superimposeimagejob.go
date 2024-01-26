package jobs

import (
	"github.com/matty271828/ai-posters/internal/imageprocessing"
)

// SuperImposeImageJob is used to superimpose an image of a poster onto a frame.
func SuperImposeImageJob(imagePath, framePath string) (string, error) {
	// superimpose the image onto a poster frame
	resultPath, err := imageprocessing.Frame(
		framePath, imagePath, "./assets/generated/result.png",
	)
	if err != nil {
		return "", err
	}

	return resultPath, nil
}
