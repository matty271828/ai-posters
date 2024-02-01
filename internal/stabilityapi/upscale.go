package stabilityapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// Upscale is used to increase the scale of an existing image
func Upscale(seedPath, height, width string) (string, error) {
	engineId := "esrgan-v1-x2plus"

	// Build REST endpoint URL
	apiHost, hasApiHost := os.LookupEnv("API_HOST")
	if !hasApiHost {
		apiHost = "https://api.stability.ai"
	}
	reqUrl := apiHost + "/v1/generation/" + engineId + "/image-to-image/upscale"

	// Acquire an API key from the environment
	apiKey, hasAPIKey := os.LookupEnv("STABILITY_API_KEY")
	if !hasAPIKey {
		return "", fmt.Errorf("missing STABILITY_API_KEY environment variable")
	}

	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)

	// Check File Existence and Permissions
	if _, err := os.Stat(seedPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist at %s", seedPath)
	}

	// Write the init image to the request
	initImageWriter, _ := writer.CreateFormField("image")
	initImageFile, initImageErr := os.Open(seedPath)
	if initImageErr != nil {
		return "", fmt.Errorf("could not open %s", seedPath)
	}
	_, _ = io.Copy(initImageWriter, initImageFile)

	// Write the options to the request
	// Only one can be set out of height and width, so we
	// give precedence to height.
	if height != "" {
		width = "width"
	}

	_ = writer.WriteField(width, height)
	writer.Close()

	// Execute the request
	payload := bytes.NewReader(data.Bytes())
	req, _ := http.NewRequest("POST", reqUrl, payload)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Accept", "image/png")
	req.Header.Add("Authorization", "Bearer "+apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var body map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return "", err
		}
		return "", fmt.Errorf("Non-200 response: %s", body)
	}

	// Write the response to a file
	outFile := fmt.Sprintf("./assets/out/v1_upscaled_image.png")
	out, err := os.Create(outFile)
	defer out.Close()
	if err != nil {
		return "", err
	}

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return "", err
	}

	return outFile, err
}
