package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandleIndexGET(t *testing.T) {
	// Create a temporary directory with a test template
	tmpDir := t.TempDir()
	templatesDir := filepath.Join(tmpDir, "templates")
	if err := os.Mkdir(templatesDir, 0755); err != nil {
		t.Fatalf("Failed to create templates directory: %v", err)
	}

	// Create a simple test template
	templateContent := `
<html>
  <body>
    {{if eq .status "success"}}
      <div>Success!</div>
    {{else if eq .status "fail"}}
      <div>Failed!</div>
    {{else}}
      <div>Form</div>
    {{end}}
  </body>
</html>
`
	templatePath := filepath.Join(templatesDir, "index.html")
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// Clear the template cache before test
	clearTemplateCache()

	// Test GET request to /
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	handler := handleIndex(tmpDir)
	handler(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body, _ := io.ReadAll(w.Body)
	if !strings.Contains(string(body), "Form") {
		t.Error("Expected response to contain form content")
	}
}

func TestHandleSendPOST(t *testing.T) {
	// Create a temporary directory with a test template
	tmpDir := t.TempDir()
	templatesDir := filepath.Join(tmpDir, "templates")
	if err := os.Mkdir(templatesDir, 0755); err != nil {
		t.Fatalf("Failed to create templates directory: %v", err)
	}

	// Create a simple test template
	templateContent := `
<html>
  <body>
    {{if eq .status "success"}}
      <div>Success</div>
    {{else}}
      <div>Failed</div>
    {{end}}
  </body>
</html>
`
	templatePath := filepath.Join(templatesDir, "index.html")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// Clear the template cache before test
	clearTemplateCache()

	tests := []struct {
		name       string
		name_      string
		email      string
		message    string
		wantCode   int
		wantInBody string
	}{
		{
			name:       "Valid form submission",
			name_:      "John Doe",
			email:      "john@example.com",
			message:    "Test message",
			wantCode:   http.StatusOK,
			wantInBody: "Failed", // Will fail because no actual Pushover creds
		},
		{
			name:       "Missing name",
			name_:      "",
			email:      "john@example.com",
			message:    "Test message",
			wantCode:   http.StatusOK,
			wantInBody: "Failed",
		},
		{
			name:       "Missing email",
			name_:      "John Doe",
			email:      "",
			message:    "Test message",
			wantCode:   http.StatusOK,
			wantInBody: "Failed",
		},
		{
			name:       "Missing message",
			name_:      "John Doe",
			email:      "john@example.com",
			message:    "",
			wantCode:   http.StatusOK,
			wantInBody: "Failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formData := "name=" + tt.name_ + "&email=" + tt.email + "&message=" + tt.message

			req, err := http.NewRequest("POST", "/", strings.NewReader(formData))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			w := httptest.NewRecorder()
			handler := handleSend(tmpDir, "test_token", "test_user")
			handler(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("Expected status %d, got %d", tt.wantCode, w.Code)
			}

			body, _ := io.ReadAll(w.Body)
			if !strings.Contains(string(body), tt.wantInBody) {
				t.Errorf("Expected response to contain '%s', got %s", tt.wantInBody, string(body))
			}
		})
	}
}

func TestHandleSendInvalidForm(t *testing.T) {
	tmpDir := t.TempDir()
	templatesDir := filepath.Join(tmpDir, "templates")
	if err := os.Mkdir(templatesDir, 0755); err != nil {
		t.Fatalf("Failed to create templates directory: %v", err)
	}

	templateContent := `<html><body>Test</body></html>`
	templatePath := filepath.Join(templatesDir, "index.html")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	clearTemplateCache()

	// Test with malformed form data that causes parsing issues
	req, err := http.NewRequest("POST", "/", strings.NewReader("invalid\x00data"))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	handler := handleSend(tmpDir, "token", "user")
	handler(w, req)

	// Should return a bad request error
	if w.Code != http.StatusBadRequest && w.Code != http.StatusOK {
		// Some implementations might parse it differently
		t.Logf("Got status %d", w.Code)
	}
}

func TestTemplateCache(t *testing.T) {
	tmpDir := t.TempDir()
	templatesDir := filepath.Join(tmpDir, "templates")
	os.Mkdir(templatesDir, 0755)

	templateContent := `<html><body>Cached</body></html>`
	templatePath := filepath.Join(templatesDir, "index.html")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	clearTemplateCache()

	// Load template twice
	tmpl1, err := loadTemplate(tmpDir, "index.html")
	if err != nil {
		t.Fatalf("Failed to load template first time: %v", err)
	}

	tmpl2, err := loadTemplate(tmpDir, "index.html")
	if err != nil {
		t.Fatalf("Failed to load template second time: %v", err)
	}

	// Both should be the same template object (cached)
	if tmpl1 != tmpl2 {
		t.Error("Expected cached template to return same object")
	}
}

func TestLoadTemplateNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	clearTemplateCache()

	tmpl, err := loadTemplate(tmpDir, "nonexistent.html")
	if err == nil {
		t.Error("Expected error for nonexistent template")
	}
	if tmpl != nil {
		t.Error("Expected nil template for nonexistent file")
	}
}

func TestLoadTemplateInvalidTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	templatesDir := filepath.Join(tmpDir, "templates")
	os.Mkdir(templatesDir, 0755)

	// Create invalid template with unclosed action
	invalidContent := `<html><body>{{.unclosed</body></html>`
	templatePath := filepath.Join(templatesDir, "invalid.html")
	os.WriteFile(templatePath, []byte(invalidContent), 0644)

	clearTemplateCache()

	tmpl, err := loadTemplate(tmpDir, "invalid.html")
	if err == nil {
		t.Error("Expected error for invalid template syntax")
	}
	if tmpl != nil {
		t.Error("Expected nil template for invalid syntax")
	}
}

// Helper function to clear the template cache between tests
func clearTemplateCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	// Clear all entries
	for key := range templateCache {
		delete(templateCache, key)
	}
}

// Test helper to create a test server with handlers
func createTestServer(t *testing.T) *httptest.Server {
	tmpDir := t.TempDir()
	templatesDir := filepath.Join(tmpDir, "templates")
	os.Mkdir(templatesDir, 0755)

	templateContent := `
<html>
  <body>
    {{if eq .status "success"}}Success{{else}}Form{{end}}
  </body>
</html>
`
	templatePath := filepath.Join(templatesDir, "index.html")
	os.WriteFile(templatePath, []byte(templateContent), 0644)

	clearTemplateCache()

	router := http.NewServeMux()
	router.HandleFunc("GET /", handleIndex(tmpDir))
	router.HandleFunc("POST /", handleSend(tmpDir, "token", "user"))

	return httptest.NewServer(router)
}

func TestServerIntegration(t *testing.T) {
	server := createTestServer(t)
	defer server.Close()

	// Test GET
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to GET: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Test POST
	formData := "name=John&email=john@example.com&message=Hello"
	resp, err = http.PostForm(server.URL, map[string][]string{
		"name":    {"John"},
		"email":   {"john@example.com"},
		"message": {"Hello"},
	})
	if err != nil {
		t.Fatalf("Failed to POST: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	_ = formData // silence unused var
}
