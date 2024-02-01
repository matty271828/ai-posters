package imageprocessing

import (
	"image"

	"github.com/disintegration/imaging"
)

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
	frameWidth, frameHeight := 350, 350 // Adjust if necessary for your frame

	// Determine the scaling factor to maintain aspect ratio.
	scaleWidth := float64(frameWidth) / float64(poster.Bounds().Dx())
	scaleHeight := float64(frameHeight) / float64(poster.Bounds().Dy())
	scaleFactor := min(scaleWidth, scaleHeight)

	// Scale down the poster image.
	newWidth := int(float64(poster.Bounds().Dx()) * scaleFactor)
	newHeight := int(float64(poster.Bounds().Dy()) * scaleFactor)
	resizedPoster := imaging.Resize(poster, newWidth, newHeight, imaging.Lanczos)

	// Calculate the center position for the poster.
	frameCenterX, frameCenterY := 342+frameWidth/2, 325+frameHeight/2 // Adjust these base coordinates to the actual frame position
	posterX := frameCenterX - newWidth/2
	posterY := frameCenterY - newHeight/2

	// Superimpose the resized poster onto the frame at the calculated position.
	result := imaging.Overlay(frame, resizedPoster, image.Pt(posterX, posterY), 1.0)

	// Save the output image
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
