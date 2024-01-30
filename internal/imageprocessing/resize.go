package imageprocessing

import (
	"fmt"
	"math"
	"os"

	"github.com/disintegration/imaging"
)

const (
	r1024x1024 = 1024.0 / 1024.0
	r1152x896  = 1152.0 / 896.0
	r1216x832  = 1216.0 / 832.0
	r1344x768  = 1344.0 / 768.0
	r1536x640  = 1536.0 / 640.0
	r640x1536  = 640.0 / 1536.0
	r768x1344  = 768.0 / 1344.0
	r832x1216  = 832.0 / 1216.0
	r896x1152  = 896.0 / 1152.0
)

// Valid image sizes allowed by Stability API
var validDimensions = map[float64][2]int{
	r1024x1024: {1024, 1024},
	r1152x896:  {1152, 896},
	r1216x832:  {1216, 832},
	r1344x768:  {1344, 768},
	r1536x640:  {1536, 640},
	r640x1536:  {640, 1536},
	r768x1344:  {768, 1344},
	r832x1216:  {832, 1216},
	r896x1152:  {896, 1152},
}

// Resize is used to ensure an uploaded image fits into dimensions
// allowed for stability API requests.
func Resize(imagePath string) (string, error) {
	// Load frame image
	image, err := imaging.Open(imagePath)
	if err != nil {
		return "", err
	}

	// Calculate the aspect ratio of the uploaded image.
	uploadedAspectRatio := float64(image.Bounds().Dx()) / float64(image.Bounds().Dy())

	required, err := CheckResizeRequired(imagePath)
	if err != nil {
		return "", err
	}

	if !required {
		return imagePath, nil
	}

	var selectedDimension float64
	closestAspectRatioDiff := float64(999999) // Initialize with a large value.
	for dimension := range validDimensions {
		aspectRatioDiff := math.Abs(uploadedAspectRatio - dimension)
		if aspectRatioDiff < closestAspectRatioDiff {
			closestAspectRatioDiff = aspectRatioDiff
			selectedDimension = dimension
		}
	}

	// Get the width and height from the selected dimension.
	width, height := validDimensions[selectedDimension][0], validDimensions[selectedDimension][1]

	// Resize the uploaded image to match the selected dimension while maintaining the aspect ratio.
	resizedImage := imaging.Resize(image, width, height, imaging.Lanczos)

	// Ensure the "./assets/out/" directory exists
	outPath := "./assets/out"

	// Ensure the output directory exists
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		err := os.MkdirAll(outPath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	// Create a temporary file within the "./assets/out/" directory
	filePath := outPath + "/resized_img.png"
	tmpFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create resized file: %v", err)
	}
	defer tmpFile.Close() // Ensure file is closed in all cases

	// Save the output image to outputPath
	err = imaging.Save(resizedImage, filePath)
	if err != nil {
		return "", err
	}

	// Return the output path and no error
	return filePath, nil
}

// CheckResizeRequired is used to validate whether an image needs to be
// resized for the stability API.
func CheckResizeRequired(imagePath string) (bool, error) {
	// Load the image
	image, err := imaging.Open(imagePath)
	if err != nil {
		return false, err
	}

	// Get the dimensions of the uploaded image
	uploadedWidth := image.Bounds().Dx()
	uploadedHeight := image.Bounds().Dy()

	// Check if the dimensions match any of the valid dimensions
	for _, dims := range validDimensions {
		if uploadedWidth == dims[0] && uploadedHeight == dims[1] {
			return false, nil // No resize required if the dimensions match exactly
		}
	}

	// If no match is found, resizing is required
	return true, nil
}
