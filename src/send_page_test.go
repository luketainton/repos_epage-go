package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestPushoverPayloadStructure(t *testing.T) {
	payload := PushoverPayload{
		Token:    "token123",
		User:     "user456",
		Title:    "Test Title",
		Message:  "Test Message",
		Priority: 1,
		Sound:    "cosmic",
	}

	// Test JSON marshaling
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	// Verify JSON contains expected fields
	jsonStr := string(data)
	expectedFields := []string{"token123", "user456", "Test Title", "Test Message", "cosmic"}
	for _, field := range expectedFields {
		if !strings.Contains(jsonStr, field) {
			t.Errorf("Expected JSON to contain '%s'", field)
		}
	}
}

func TestPushoverMessageFormat(t *testing.T) {
	tests := []struct {
		name     string
		apiToken string
		userKey  string
		name_    string
		email    string
		message  string
	}{
		{
			name:     "Valid inputs",
			apiToken: "token",
			userKey:  "user",
			name_:    "Alice",
			email:    "alice@example.com",
			message:  "Hello",
		},
		{
			name:     "Special characters",
			apiToken: "token",
			userKey:  "user",
			name_:    "John O'Brien",
			email:    "john+tag@example.com",
			message:  "Message with\nmultiple\nlines",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := PushoverPayload{
				Token:    tt.apiToken,
				User:     tt.userKey,
				Title:    "ePage from " + tt.name_,
				Message:  "Name: " + tt.name_ + "\nEmail: " + tt.email + "\n\nMessage: " + tt.message,
				Priority: 1,
				Sound:    "cosmic",
			}

			// Verify structure
			if payload.Token != tt.apiToken {
				t.Errorf("Token mismatch")
			}
			if payload.Priority != 1 {
				t.Errorf("Priority should be 1")
			}
			if payload.Sound != "cosmic" {
				t.Errorf("Sound should be 'cosmic'")
			}
		})
	}
}
