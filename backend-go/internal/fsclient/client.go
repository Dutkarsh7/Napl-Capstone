package fsclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"smart-pharmacy/backend-go/internal/models"
)

type FieldIssue struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

type ValidationResponse struct {
	Status  string       `json:"status"` // "valid", "incomplete", "invalid"
	Missing []string     `json:"missing"`
	Issues  []FieldIssue `json:"issues"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Validate sends the parsed prescription data to the F# Giraffe validation microservice
func (c *Client) Validate(ctx context.Context, parsed models.ParsedRawPrescription) (*ValidationResponse, error) {
	bodyBytes, err := json.Marshal(parsed)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize validation payload: %w", err)
	}

	url := fmt.Sprintf("%s/validate/prescription", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create validation request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("F# validation service request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("F# validation service returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var res ValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode F# validation response: %w", err)
	}

	return &res, nil
}
