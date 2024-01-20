package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	stabilityapi "github.com/matty271828/ai-posters/internal/app/stabilityapi"
)

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	imagePaths, err := stabilityapi.GenerateImage("A lighthouse on a cliff")
	if err != nil {
		log.Fatalf("Error generating image: %v", err)
	}

	fmt.Println("Generated images:", imagePaths)
}
