package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// ResponseData holds the structure for the data received from the LLM.
type ResponseData struct {
	Date     string      `json:"date"`
	Payee    string      `json:"payee"`
	Notes    string      `json:"notes"`
	Category string      `json:"category"`
	Amount   json.Number `json:"amount"`
}

// RequestData structures the request to send to the LLM.
type RequestData struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Format string `json:"format"`
}

// QueryLLM queries the LLM model and returns the parsed response data.
// It constructs a request, executes it, and processes the response.
func QueryLLM(model, prompt string) (*ResponseData, error) {
	url := "http://localhost:11434/api/generate"
	client := &http.Client{}
	requestData := RequestData{
		Model:  model,
		Prompt: prompt,
		Stream: false,
		Format: "json",
	}

	// Serialize the request data to JSON.
	data, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	// Create a new HTTP request with a context.
	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		url,
		bytes.NewBuffer(data),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP request.
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read and decode the HTTP response body.
	var result map[string]interface{}
	if err := decodeResponse(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	// Process the total duration for logging.
	if totalDurationMs, ok := result["total_duration"].(float64); ok {
		log.Printf(
			"Total duration: %s\n",
			time.Duration(totalDurationMs).String(),
		)
	} else {
		log.Printf("Total duration data missing or in wrong format in response: %v\n", result)
	}

	// Extract and parse the main response data.
	responseData, err := parseResponseData(result)
	if err != nil {
		return nil, err
	}

	// Parse and format the amount if present.
	if err := formatAmount(responseData); err != nil {
		log.Printf("Failed to parse amount: %v", err)
		return nil, fmt.Errorf("failed to parse amount: %w", err)
	}

	log.Printf("Processed response data: %+v", responseData)
	return responseData, nil
}

// decodeResponse decodes the response body into the provided target interface.
func decodeResponse(r io.Reader, target interface{}) error {
	err := json.NewDecoder(r).Decode(target)
	if err != nil {
		log.Printf("Error decoding response body: %v", err)
		log.Printf("Response body: %s", dumpResponse(r))
	}
	return err
}

// parseResponseData extracts and parses the main response data from the result map.
func parseResponseData(result map[string]interface{}) (*ResponseData, error) {
	responseData := &ResponseData{}
	if response, ok := result["response"].(string); ok {
		if err := json.Unmarshal([]byte(response), responseData); err != nil {
			log.Printf("Error unmarshaling response data: %v", err)
			log.Printf("Response data: %s", response)
			return nil, fmt.Errorf("failed to unmarshal response data: %w", err)
		}
	} else {
		return nil, errors.New("response data is missing or not a string")
	}
	return responseData, nil
}

// formatAmount parses and formats the amount field in the ResponseData struct.
func formatAmount(responseData *ResponseData) error {
	amt, err := responseData.Amount.Float64()
	if err != nil {
		return fmt.Errorf("failed to parse amount: %w", err)
	}
	responseData.Amount = json.Number(fmt.Sprintf("%.2f", amt))
	return nil
}

// logTotalDuration logs the total duration from the result map if available.
func dumpResponse(r io.Reader) string {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
	}
	return buf.String()
}
