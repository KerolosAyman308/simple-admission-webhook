package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

type TestResult struct {
	Allowed bool `json:"allowed"`
}

func TestAdmissionWebhook(t *testing.T) {
	// Root directory of the project (assuming tests run from /cmd)
	wd, _ := os.Getwd()
	projectRoot := filepath.Join(wd, "..")

	tests := []struct {
		name     string
		handler  http.HandlerFunc
		jsonFile string
		expected bool
	}{
		{"Validation Allowed Case", Validate, "tests/validate-allowed.json", true},
		{"Validation Disallowed Case", Validate, "tests/validate-disallowed.json", false},
		{"Validation Developer Case", Validate, "tests/validate-another-allowed.json", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join(projectRoot, tt.jsonFile)
			content, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("Failed to read test file %s: %v", tt.jsonFile, err)
			}

			req := httptest.NewRequest("POST", "/", bytes.NewBuffer(content))
			defer req.Body.Close()
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			tt.handler(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status OK; got %v", resp.Status)
			}

			var result TestResult
			if err := json.Unmarshal(body, &result); err != nil {
				t.Fatalf("Failed to decode response: %v. Body: %s", err, string(body))
			}

			if result.Allowed != tt.expected {
				t.Errorf("%s: expected allowed=%v; got %v. Response: %s", tt.name, tt.expected, result.Allowed, string(body))
			}
		})
	}
}
