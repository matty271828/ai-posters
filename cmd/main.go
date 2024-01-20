package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/matty271828/ai-posters/internal/jobs"
)

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	prompt := fmt.Sprintf("A beautiful monstera plant in the style of a 19th century scientific drawing.")
	_, _, err = jobs.GenerateImageJob(prompt, "assets/stock/blackframe.png", "assets/out/result.png")
}
