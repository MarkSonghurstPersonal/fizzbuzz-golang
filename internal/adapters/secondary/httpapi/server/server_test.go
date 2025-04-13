package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func TestServer(t *testing.T) {
	server := New()
	defer server.Close()

	tests := []struct {
		name           string
		path           string
		queryParams    url.Values
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid division",
			path:           "/divide",
			queryParams:    url.Values{"a": {"10"}, "b": {"3"}},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"remainder":1}`,
		},
		{
			name:           "Invalid query parameter 'a'",
			path:           "/divide",
			queryParams:    url.Values{"a": {"invalid"}, "b": {"3"}},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Invalid query parameter: 'a'"}`,
		},
		{
			name:           "Invalid query parameter 'b'",
			path:           "/divide",
			queryParams:    url.Values{"a": {"10"}, "b": {"invalid"}},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Invalid query parameter: 'b'"}`,
		},
		{
			name:           "Division by zero",
			path:           "/divide",
			queryParams:    url.Values{"a": {"10"}, "b": {"0"}},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Division by zero is not allowed"}`,
		},
		{
			name:           "Unsupported path",
			path:           "/unsupported",
			queryParams:    url.Values{},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"Unsupported path"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Construct the request URL with query parameters
			reqURL := server.URL + tt.path + "?" + tt.queryParams.Encode()

			// Make the HTTP GET request
			resp, err := http.Get(reqURL)
			if err != nil {
				t.Fatalf("Failed to make GET request: %v", err)
			}
			defer resp.Body.Close()

			// Read the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			// Check the status code
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			// Check the response body
			var actualBody map[string]any
			var expectedBody map[string]any

			if err := json.Unmarshal(body, &actualBody); err != nil {
				t.Fatalf("Failed to unmarshal response body: %v", err)
			}
			if err := json.Unmarshal([]byte(tt.expectedBody), &expectedBody); err != nil {
				t.Fatalf("Failed to unmarshal expected body: %v", err)
			}

			if len(actualBody) != len(expectedBody) {
				t.Errorf("Expected body %v, got %v", expectedBody, actualBody)
			}
			for key, expectedValue := range expectedBody {
				if actualValue, ok := actualBody[key]; !ok || actualValue != expectedValue {
					t.Errorf("Expected body %v, got %v", expectedBody, actualBody)
				}
			}
		})
	}
}
