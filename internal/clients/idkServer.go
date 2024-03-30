package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func CreateIDKToken(accessToken string, refreshToken string, idkBackendBaseUrl string) (string, error) {
	requestBodyMap := map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}

	requestBodyBytes, err := json.Marshal(requestBodyMap)
	if err != nil {
		return "", err
	}

	requestUrl := fmt.Sprintf("%s/token", idkBackendBaseUrl)
	response, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned non-OK status: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return "", err
	}

	// Navigating through the nested JSON response to extract the desired value
	token, ok := responseData["jwtToken"].(string)
	if !ok || len(token) == 0 {
		return "", fmt.Errorf("token not found")
	}

	return token, nil
}

func ProcessPrompt(prompt string, os string, readmeData string, existingScript string, pwd string, jwtToken string, idkBackendBaseUrl string) (*PromptResponse, error, int) {
	requestBodyMap := map[string]interface{}{
		"prompt":         prompt,
		"os":             os,
		"existingScript": existingScript,
		"readmeData":     readmeData,
		"pwd":            pwd,
	}

	requestBodyBytes, err := json.Marshal(requestBodyMap)
	if err != nil {
		return nil, err, 0
	}

	requestUrl := fmt.Sprintf("%s/prompt", idkBackendBaseUrl)
	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return nil, err, 0
	}
	req.Header.Set("Authorization", jwtToken)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err, 0
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned non-OK status"), response.StatusCode
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err, response.StatusCode
	}

	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, err, response.StatusCode
	}

	return &PromptResponse{
		Response:   responseData["response"].(string),
		ActionType: responseData["actionType"].(string),
	}, nil, response.StatusCode
}

type PromptResponse struct {
	Response   string
	ActionType string
}
