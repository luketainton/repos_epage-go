package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const pushoverAPIURL = "https://api.pushover.net/1/messages.json"

// PushoverPayload represents the JSON payload sent to Pushover API
type PushoverPayload struct {
	Token    string `json:"token"`
	User     string `json:"user"`
	Title    string `json:"title"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
	Sound    string `json:"sound"`
}

// PushoverResponse represents the JSON response from Pushover API
type PushoverResponse struct {
	Status  int    `json:"status"`
	Request string `json:"request"`
}

func sendPage(apiToken, userKey, name, email, message string) (bool, error) {
	fullMsg := fmt.Sprintf("Name: %s\nEmail: %s\n\nMessage: %s", name, email, message)

	payload := PushoverPayload{
		Token:    apiToken,
		User:     userKey,
		Title:    fmt.Sprintf("ePage from %s", name),
		Message:  fullMsg,
		Priority: 1,
		Sound:    "cosmic",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Errorf("failed to marshal payload: %w", err)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("POST", pushoverAPIURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %w", err)
	}

	var response PushoverResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if resp.StatusCode == http.StatusOK && response.Status == 1 {
		return true, nil
	}

	return false, fmt.Errorf("pushover API returned status %d", resp.StatusCode)
}
