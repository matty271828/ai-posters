package stabilityapi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type TextToImageImage struct {
	Base64       string `json:"base64"`
	Seed         uint32 `json:"seed"`
	FinishReason string `json:"finishReason"`
}

type TextToImageResponse struct {
	Images []TextToImageImage `json:"artifacts"`
}

func GenerateImage(prompt string) ([]string, error) {
	// Build REST endpoint URL w/ specified engine
	engineId := "stable-diffusion-v1-6"
	apiHost, hasApiHost := os.LookupEnv("API_HOST")
	if !hasApiHost {
		apiHost = "https://api.stability.ai"
	}
	reqUrl := apiHost + "/v1/generation/" + engineId + "/text-to-image"

	// Acquire an API key from the environment
	apiKey, hasAPIKey := os.LookupEnv("STABILITY_API_KEY")
	if !hasAPIKey {
		return nil, fmt.Errorf("Missing STABILITY_API_KEY environment variable")
	}

	jsonData := fmt.Sprintf(`{
        "text_prompts": [
            {
                "text": "%s"
            }
        ],
        "cfg_scale": 7,
        "height": 1024,
        "width": 1024,
        "samples": 1,
        "steps": 30
    }`, prompt)

	data := []byte(jsonData)

	req, _ := http.NewRequest("POST", reqUrl, bytes.NewBuffer(data))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey)

	// Execute the request & read all the bytes of the body
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
		return nil, fmt.Errorf(fmt.Sprintf("Non-200 response: %s", body))
	}

	// Decode the JSON body
	var body TextToImageResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	err = os.MkdirAll("./assets/out", 0755) // 0755 is a common permission setting allowing read and execute access
	if err != nil {
		return nil, err
	}

	// Write the images to disk
	var savedFilePaths []string
	for i, image := range body.Images {
		outFile := fmt.Sprintf("./assets/out/v1_txt2img_%d.png", i)
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
