package stabilityapi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type ImageToImageImage struct {
	Base64       string `json:"base64"`
	Seed         uint32 `json:"seed"`
	FinishReason string `json:"finishReason"`
}

type ImageToImageResponse struct {
	Images []ImageToImageImage `json:"artifacts"`
}

func GenerateImageToImage(prompt, seedPath string) ([]string, error) {
	engineId := "stable-diffusion-xl-1024-v1-0"

	// Build REST endpoint URL
	apiHost, hasApiHost := os.LookupEnv("API_HOST")
	if !hasApiHost {
		apiHost = "https://api.stability.ai"
	}
	reqUrl := apiHost + "/v1/generation/" + engineId + "/image-to-image"

	// Acquire an API key from the environment
	apiKey, hasAPIKey := os.LookupEnv("STABILITY_API_KEY")
	if !hasAPIKey {
		return nil, fmt.Errorf("missing STABILITY_API_KEY environment variable")
	}

	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)

	// Write the init image to the request
	initImageWriter, _ := writer.CreateFormField("init_image")
	initImageFile, initImageErr := os.Open(seedPath)
	if initImageErr != nil {
		return nil, fmt.Errorf("could not open %s", seedPath)
	}
	_, _ = io.Copy(initImageWriter, initImageFile)

	// Write the options to the request
	_ = writer.WriteField("init_image_mode", "IMAGE_STRENGTH")
	_ = writer.WriteField("image_strength", "0.35")
	_ = writer.WriteField("text_prompts[0][text]", prompt)
	_ = writer.WriteField("cfg_scale", "7")
	_ = writer.WriteField("samples", "1")
	_ = writer.WriteField("steps", "30")
	writer.Close()

	// Execute the request
	payload := bytes.NewReader(data.Bytes())
	req, _ := http.NewRequest("POST", reqUrl, payload)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		var body map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Non-200 response: %s", body)
	}

	// Decode the JSON body
	var body ImageToImageResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	// Write the images to disk
	var savedFilePaths []string
	for i, image := range body.Images {
		outFile := fmt.Sprintf("./out/v1_img2img_%d.png", i)
		savedFilePaths = append(savedFilePaths, outFile)
		file, err := os.Create(outFile)
		if err != nil {
			return nil, err
		}

		imageBytes, err := base64.StdEncoding.DecodeString(image.Base64)
		if err != nil {
			return nil, err
		}

		if _, err := file.Write(imageBytes); err != nil {
			return nil, err
		}

		if err := file.Close(); err != nil {
			return nil, err
		}
	}
	return savedFilePaths, nil
}
