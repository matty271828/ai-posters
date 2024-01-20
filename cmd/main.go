package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/matty271828/ai-posters/internal/imageprocessing"
	stabilityapi "github.com/matty271828/ai-posters/internal/stabilityapi"
)

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// generate an image
	imagePaths, err := stabilityapi.GenerateImage("A lighthouse on a cliff")
	if err != nil {
		log.Fatalf("Error generating image: %v", err)
	}
	fmt.Println("Generated images:", imagePaths)

	// superimpose the image onto a poster frame
	outputPath, err := imageprocessing.Frame("assets/stock/blackframe.png", imagePaths[0], "assets/out/result.png")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Println("Image saved to: ", outputPath)
}
