package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
)

type DivisionResult struct {
	Remainder int `json:"remainder"`
}

type ErrorResult struct {
	Message string `json:"message"`
}

// newErrorJSON creates a JSON error response with the given message.
func newErrorJSON(message string) []byte {
	result := ErrorResult{
		Message: message,
	}
	errorJSON, err := json.Marshal(result)
	if err != nil {
		return []byte(`{"message":"Internal Server Error"}`)
	}
	return errorJSON
}

// getIntFromQuery retrieves an integer parameter from the query string.
func getIntFromQuery(q url.Values, param string) (int, error) {
	value, err := strconv.Atoi(q.Get(param))
	if err != nil {
		return 0, err
	}
	return value, nil
}

// New creates a new HTTP test server for handling requests.
func New() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			errorJSON  []byte
			resultJSON []byte
			statusCode = http.StatusBadRequest
		)

		defer func() {
			w.Header().Set("Content-Type", "application/json")

			switch {
			case len(errorJSON) > 0:
				w.WriteHeader(statusCode)
				w.Write(errorJSON) // if errorJSON is nil, it will write "{}"

			case statusCode == http.StatusOK:
				w.Write(resultJSON) // if resultJSON is nil, it will write "{}"

			default:
				// Safety net, it shouldn't be possible to reach here.
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(newErrorJSON("Internal Server Error"))
			}
		}()

		switch r.URL.Path {
		case "/divide":
			a, err := getIntFromQuery(r.URL.Query(), "a")
			if err != nil {
				errorJSON = newErrorJSON("Invalid query parameter: 'a'")
				statusCode = http.StatusBadRequest
				return
			}

			b, err := getIntFromQuery(r.URL.Query(), "b")
			if err != nil {
				errorJSON = newErrorJSON("Invalid query parameter: 'b'")
				statusCode = http.StatusBadRequest
				return
			}

			if b == 0 {
				errorJSON = newErrorJSON("Division by zero is not allowed")
				statusCode = http.StatusBadRequest
				return
			}

			result := DivisionResult{
				Remainder: a % b,
			}
			statusCode = http.StatusOK

			resultJSON, err = json.Marshal(result)
			if err != nil {
				errorJSON = newErrorJSON("Failed to marshal JSON")
				statusCode = http.StatusInternalServerError
			}

		default:
			errorJSON = newErrorJSON("Unsupported path")
			statusCode = http.StatusNotFound
		}
	}))
}
