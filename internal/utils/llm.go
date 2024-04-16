package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// ResponseData holds structured data parsed from the response JSON.
type ResponseData struct {
	Date     string      `json:"date"`
	Payee    string      `json:"payee"`
	Notes    string      `json:"notes"`
	Category string      `json:"category"`
	Amount   json.Number `json:"amount"`
}

type RequestData struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Format string `json:"format"`
}

// QueryLLM queries the LLM model and returns the parsed response data.
func QueryLLM(model string, prompt string) (*ResponseData, error) {
	url := "http://localhost:11434/api/generate"

	data, err := json.Marshal(RequestData{
		Model:  model,
		Prompt: prompt,
		Stream: false,
		Format: "json",
	})
	if err != nil {
		log.Printf("Failed to marshal data to JSON: %v", err)
		return nil, fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Failed to create HTTP request: %v", err)
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to execute HTTP request: %v", err)
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read HTTP response body: %v", err)
		return nil, fmt.Errorf("failed to read HTTP response body: %w", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("Failed to unmarshal JSON response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	// Print the raw response for debugging
	log.Printf("Raw response: %+v\n", result)

	responseData := &ResponseData{}
	if response, ok := result["response"].(string); ok {
		err = json.Unmarshal([]byte(response), responseData)
		if err != nil {
			log.Printf("Failed to unmarshal response data: %v", err)
			return nil, fmt.Errorf("failed to unmarshal response data: %w", err)
		}
	} else {
		log.Println("Response data is not in expected format or missing")
		return nil, fmt.Errorf("response data is not in expected format or missing")
	}

	// Format the amount to ensure it is properly parsed
	if amt, err := responseData.Amount.Float64(); err == nil {
		responseData.Amount = json.Number(fmt.Sprintf("%.2f", amt))
	} else {
		log.Printf("Failed to parse amount from response: %v", err)
		return nil, fmt.Errorf("failed to parse amount: %w", err)
	}

	log.Printf("Response: %+v\n", responseData)

	return responseData, nil
}
