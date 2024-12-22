package application

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCalcHandler(t *testing.T) {
	tests := []struct {
		name        string
		requestBody string

		expectedStatus int
		expectedResult float64
	}{
		{
			name:           "Successful calculation",
			requestBody:    `{"expression": "1 + 2"}`,
			expectedStatus: http.StatusOK,
			expectedResult: 3,
		},
		{
			name:           "Invalid expression",
			requestBody:    `{"expression": "invalid"}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Malformed JSON",
			requestBody:    `{invalid}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/v1/calculate/", bytes.NewBufferString(tt.requestBody))
			w := httptest.NewRecorder()
			CalcHandler(w, req)

			res := w.Result()
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}
			if tt.expectedStatus == http.StatusOK {
				var response map[string]float64
				json.NewDecoder(res.Body).Decode(&response)
				if response["result"] != tt.expectedResult {
					t.Errorf("expected result %f, got %f", tt.expectedResult, response["result"])
				}
			}
		})
	}
}

func TestRunServer(t *testing.T) {
	go func() {
		app := New()
		app.RunServer()
	}()

	resp, err := http.Post("http://localhost:8080/api/v1/calculate/", "application/json", bytes.NewBufferString(`{"expression": "1 + 2"}`))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", resp.Status)
	}
}

// Простой тест для ConfigFromEnv
func TestConfigFromEnv(t *testing.T) {
	os.Setenv("PORT", "3000")
	config := ConfigFromEnv()
	if config.Addr != "3000" {
		t.Errorf("expected port 3000, got %s", config.Addr)
	}

	os.Unsetenv("PORT")
	config = ConfigFromEnv()
	if config.Addr != "8080" {
		t.Errorf("expected default port 8080, got %s", config.Addr)
	}
}
