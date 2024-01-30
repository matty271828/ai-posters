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

	// Dimensions where the poster will fit.
	frameWidth, frameHeight := 350, 350

	// Determine the scaling factor.
	scaleWidth := float64(frameWidth) / float64(poster.Bounds().Dx())
	scaleHeight := float64(frameHeight) / float64(poster.Bounds().Dy())
	scaleFactor := min(scaleWidth, scaleHeight)

	// Scale down the poster image while maintaining its aspect ratio.
	newWidth := int(float64(poster.Bounds().Dx()) * scaleFactor)
	newHeight := int(float64(poster.Bounds().Dy()) * scaleFactor)
	resizedPoster := imaging.Resize(poster, newWidth, newHeight, imaging.Lanczos)

	// Coordinates where the poster should be placed on the frame.
	frameX, frameY := 342, 325

	// Superimpose the resized poster onto the frame.
	result := imaging.Overlay(frame, resizedPoster, image.Pt(frameX, frameY), 1.0)

	// Save the output image to outputPath
	err = imaging.Save(result, outputPath)
	if err != nil {
		return "", err
	}

	// Return the output path and no error
	return outputPath, nil
}

// min returns the smaller of x or y.
func min(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}
