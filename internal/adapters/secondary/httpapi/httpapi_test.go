package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/MarkSonghurstPersonal/fizzbuzz-golang/internal/adapters/secondary/httpapi/server"
)

const (
	ignoredValue = 0 // For test cases where the value being divided doesn't matter.
)

func TestDivide(t *testing.T) {

	tests := []struct {
		name             string
		padResponseBytes int // set to > 0 to include a padded additional response of this many bytes.
		apiResponse      any
		apiStatusCode    int
		expectedResult   int
		expectedError    string
	}{
		{
			name:           "Valid division",
			apiResponse:    server.DivisionResult{Remainder: 1},
			apiStatusCode:  http.StatusOK,
			expectedResult: 1,
			expectedError:  "",
		},
		{
			name:          "Division by zero",
			apiResponse:   server.ErrorResult{Message: "Division by zero is not allowed"},
			apiStatusCode: http.StatusBadRequest,
			expectedError: "400 Bad Request: Division by zero is not allowed",
		},
		{
			name:          "Some other client-side error",
			apiResponse:   server.ErrorResult{Message: "Some other client-side error"},
			apiStatusCode: http.StatusConflict,
			expectedError: "409 Conflict: Some other client-side error",
		},
		{
			name:          "Internal server error",
			apiResponse:   server.ErrorResult{Message: "Internal server error"},
			apiStatusCode: http.StatusInternalServerError,
			expectedError: "500 Internal Server Error: Internal server error",
		},
		{
			name:          "Bogus JSON response on a 200",
			apiResponse:   "this is not JSON",
			apiStatusCode: http.StatusOK,
			expectedError: "failed to decode result response: json: cannot unmarshal string into Go value of type server.DivisionResult",
		},
		{
			name:          "Bogus JSON response on a 400",
			apiResponse:   "this is not JSON",
			apiStatusCode: http.StatusBadRequest,
			expectedError: "failed to decode error response for status code 400: json: cannot unmarshal string into Go value of type server.ErrorResult",
		},
		{
			name:          "Unexpected Status Code",
			apiResponse:   server.ErrorResult{Message: "These are not the droids you're looking for"},
			apiStatusCode: http.StatusSeeOther,
			expectedError: "unexpected status code: 303 See Other: These are not the droids you're looking for",
		},
		{
			name:             "Content Length too large",
			padResponseBytes: maxResponseSize + 1,
			apiResponse:      server.DivisionResult{Remainder: 1},
			apiStatusCode:    http.StatusOK,
			expectedError:    "response too large: 1041 bytes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock server
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.apiStatusCode)

				if err := json.NewEncoder(w).Encode(tt.apiResponse); err != nil {
					t.Fatalf("Failed to encode mock response: %v", err)
				}
				if tt.padResponseBytes > 0 {
					w.Write(make([]byte, tt.padResponseBytes))
				}
			}))
			defer mockServer.Close()

			// Create an API instance using the mock server
			api := API{
				server: mockServer,
			}

			// Call the divide function. The values of a and b don't matter, as we're using canned responses.

			result, err := api.divide(ignoredValue, ignoredValue)

			// Check the result
			if result != tt.expectedResult {
				t.Errorf("Expected result %d, got %d", tt.expectedResult, result)
			}

			// Check the error
			if err != nil {
				if tt.expectedError == "" {
					t.Errorf("Unexpected error: %v", err)
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error %q, got %q", tt.expectedError, err.Error())
				}
			} else if tt.expectedError != "" {
				t.Errorf("Expected error %q, got nil", tt.expectedError)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		ctx context.Context
		wg  *sync.WaitGroup
	}
	tests := []struct {
		name                  string
		args                  args
		expectValidAPI        bool
		expectValidCancelFunc bool
		expectError           bool
	}{
		{
			name: "Valid context and wait group",
			args: args{
				ctx: context.Background(),
				wg:  &sync.WaitGroup{},
			},
			expectValidAPI:        true,
			expectValidCancelFunc: true,
		},
		{
			name: "Nil context, one created for us",
			args: args{
				ctx: nil,
				wg:  &sync.WaitGroup{},
			},
			expectValidAPI:        true,
			expectValidCancelFunc: true,
		},
		{
			name: "Nil wait group",
			args: args{
				ctx: context.Background(),
				wg:  nil,
			},
			expectError: true,
		},
		{
			name: "Nil context and wait group",
			args: args{
				ctx: nil,
				wg:  nil,
			},
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAPI, gotCancelFunc, gotErr := New(tt.args.ctx, tt.args.wg)

			if tt.expectError {
				if gotErr == nil {
					t.Errorf("New() gotErr = nil, expected an error")
				}
			} else {
				if gotErr != nil {
					t.Errorf("New() gotErr = %v, expected no error", gotErr)
				}
			}

			if tt.expectValidAPI {
				if gotAPI == nil {
					t.Errorf("New() gotAPI = nil, expected a valid API")
				}
			} else {
				if gotAPI != nil {
					t.Errorf("New() gotAPI = %v, expected an Invalid API", gotAPI)
				}
			}

			if tt.expectValidCancelFunc {
				if gotCancelFunc == nil {
					t.Errorf("New() gotCancelFunc = nil, expected a valid cancel function")
				}
			} else {
				if gotCancelFunc != nil {
					t.Errorf("New() gotCancelFunc = %v, expected an Invalid cancel function", gotCancelFunc)
				}
			}

			if gotCancelFunc != nil {
				// Call the cancel function to clean up
				gotCancelFunc()
			}
		})
	}
}

func TestAPI_FizzAndBuzz(t *testing.T) {

	// functionToTest is a type alias for the Fizz and Buzz functions, so we can iterate through them in this test.
	type functionToTest func(int) bool

	// We use the same API instance for both Fizz and Buzz tests, the server is assigned in each test case.
	api := API{
		ctx: context.Background(),
		wg:  &sync.WaitGroup{},
	}
	functionsToTest := []functionToTest{
		api.Fizz,
		api.Buzz,
	}

	// Iterate over the functions to test
	// This allows us to test both Fizz and Buzz in the same test function.
	for _, funcToTest := range functionsToTest {

		type fields struct {
			server     *httptest.Server
			funcToTest functionToTest
		}
		tests := []struct {
			name         string
			fields       fields
			expectedBool bool
		}{
			{
				name: "Divisible",
				fields: fields{
					server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(server.DivisionResult{Remainder: 0})
					})),
					funcToTest: funcToTest,
				},
				expectedBool: true,
			},
			{
				name: "Not divisible",
				fields: fields{
					server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(server.DivisionResult{Remainder: 1})
					})),
					funcToTest: funcToTest,
				},
				expectedBool: false,
			},
			{
				name: "Error from server",
				fields: fields{
					server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(w).Encode(server.ErrorResult{Message: "Error from server test case"})
					})),
					funcToTest: funcToTest,
				},
				expectedBool: false,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {

				api.server = tt.fields.server

				if got := tt.fields.funcToTest(ignoredValue); got != tt.expectedBool {
					t.Errorf("Function returned %v, want %v", got, tt.expectedBool)
				}
			})
		}
	}
}
