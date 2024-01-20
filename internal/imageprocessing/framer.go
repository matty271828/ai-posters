package imageprocessing

import (
	"image"

	"github.com/disintegration/imaging"
)

// Frame takes two image paths (frame and poster) and combines them
func Frame(framePath, posterPath, outputPath string) (string, error) {
	// Load frame image
	frame, err := imaging.Open(framePath)
	if err != nil {
		return "", err
	}

	// Load poster image
	poster, err := imaging.Open(posterPath)
	if err != nil {
		return "", err
	}

	// Resize the poster to fit the frame.
	// Replace these with the actual dimensions of the area in the frame where the poster should fit.
	frameWidth, frameHeight := 350, 350
	resizedPoster := imaging.Resize(poster, frameWidth, frameHeight, imaging.Lanczos)

	// Superimpose the resized poster onto the frame.
	// Replace these with the actual coordinates where the poster should be placed on the frame.
	frameX, frameY := 342, 325
	result := imaging.Overlay(frame, resizedPoster, image.Pt(frameX, frameY), 1.0)

	// Save the output image to outputPath
	err = imaging.Save(result, outputPath)
	if err != nil {
		return "", err
	}

	// Return the output path and no error
	return outputPath, nil
}
