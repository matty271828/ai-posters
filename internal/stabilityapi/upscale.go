package stabilityapi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func Upscale(seedPath, height string) (string, error) {
	fmt.Printf("Attempting to upscale image: %s with Height: %s\n", seedPath, height)

	engineId := "esrgan-v1-x2plus"
	apiHost := os.Getenv("API_HOST")
	if apiHost == "" {
		apiHost = "https://api.stability.ai"
	}
	reqUrl := fmt.Sprintf("%s/v1/generation/%s/image-to-image/upscale", apiHost, engineId)

	apiKey := os.Getenv("STABILITY_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("missing STABILITY_API_KEY environment variable")
	}

	data, contentType, err := prepareMultipartFormData(seedPath, height)
	if err != nil {
		return "", err
	}

	resp, err := sendUpscaleRequest(reqUrl, data, contentType, apiKey)
	if err != nil {
		return "", err
	}

	outputPath, err := handleUpscaleResponse(resp)
	if err != nil {
		return "", err
	}

	fmt.Printf("Upscaled image saved to: %s\n", outputPath)
	return outputPath, nil
}

// prepareMultipartFormData prepares the multipart form data for the upscale request.
func prepareMultipartFormData(seedPath, height string) (*bytes.Buffer, string, error) {
	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)

	// Prepare the image file to be sent
	file, err := os.Open(seedPath)
	if err != nil {
		return nil, "", fmt.Errorf("could not open image file at %s: %v", seedPath, err)
	}
	defer file.Close()

	fileField, err := writer.CreateFormFile("image", filepath.Base(seedPath))
	if err != nil {
		return nil, "", fmt.Errorf("failed to create form file: %v", err)
	}
	if _, err := io.Copy(fileField, file); err != nil {
		return nil, "", fmt.Errorf("failed to copy image data to form file: %v", err)
	}

	// Adding height and width to the request as per provided arguments.
	if height != "" {
		if err := writer.WriteField("height", height); err != nil {
			return nil, "", fmt.Errorf("failed to add height to the form: %v", err)
		}
	}

	contentType := writer.FormDataContentType()
	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("failed to close multipart writer: %v", err)
	}

	return data, contentType, nil
}

// sendUpscaleRequest sends the upscale request to the server and returns the response.
func sendUpscaleRequest(reqUrl string, data *bytes.Buffer, contentType, apiKey string) (*http.Response, error) {
	req, err := http.NewRequest("POST", reqUrl, data)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %v", err)
	}

	return resp, nil
}

// handleUpscaleResponse handles the server's response to the upscale request.
func handleUpscaleResponse(resp *http.Response) (string, error) {
	if resp.StatusCode != http.StatusOK {
		var responseBody map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
			resp.Body.Close()
			return "", fmt.Errorf("failed to decode error response body: %v", err)
		}
		resp.Body.Close()
		return "", fmt.Errorf("failed to upscale image, server responded with status %d: %v", resp.StatusCode, responseBody)
	}

	// Define a structure to match the JSON response format for the base64 image data
	type Response struct {
		Artifacts []struct {
			Base64 string `json:"base64"`
		} `json:"artifacts"`
	}

	var decodedResponse Response
	if err := json.NewDecoder(resp.Body).Decode(&decodedResponse); err != nil {
		resp.Body.Close()
		return "", fmt.Errorf("failed to decode JSON response: %v", err)
	}
	resp.Body.Close()

	if len(decodedResponse.Artifacts) == 0 {
		return "", fmt.Errorf("no artifacts found in response")
	}

	base64Data := decodedResponse.Artifacts[0].Base64
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 image data: %v", err)
	}

	outPath := "./assets/out/v1_upscaled_image.png"
	if err := ioutil.WriteFile(outPath, imageData, 0644); err != nil {
		return "", fmt.Errorf("failed to write decoded image to file: %v", err)
	}

	fmt.Printf("Upscaled image saved to: %s\n", outPath)
	return outPath, nil
}
