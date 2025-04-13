package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/MarkSonghurstPersonal/fizzbuzz-golang/internal/adapters/secondary/httpapi/server"
)

const (
	maxResponseSize = 1024 // 1 KB is the maximum supported response size.
	requestPath     = "%s/divide?a=%d&b=%d"
)

type API struct {
	server *httptest.Server
	ctx    context.Context
	wg     *sync.WaitGroup
}

// New creates a new API instance with an embedded httptest Server.
// The caller should supply a context to control when the server should be closed, or the function will create one for you.
// The caller is responsible for calling the returned cancel function to cleanly stop the server,
// and should wait on the supplied WaitGroup to ensure all the server has stopped.
func New(ctx context.Context, wg *sync.WaitGroup) (*API, context.CancelFunc, error) {
	if wg == nil {
		return nil, nil, fmt.Errorf("wait group cannot be nil")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	// Create a new context from the supplied one, with a cancel function that we'll return.
	ctx, cancel := context.WithCancel(ctx)

	api := API{
		server: server.New(),
		ctx:    ctx,
		wg:     wg,
	}

	// Ensure the cancel function is called when the context is done
	api.wg.Add(1)
	go func() {
		defer api.wg.Done()

		<-api.ctx.Done()
		api.server.Close()
	}()

	return &api, cancel, nil
}

// Fizz implements the repository.FizzBuzzer interface.
// Errors are logged and result in a false return value.
func (api *API) Fizz(in int) bool {
	return api.commonDivide(in, 3)
}

// Buzz implements the repository.FizzBuzzer interface.
// Errors are logged and result in a false return value.
func (api *API) Buzz(in int) bool {
	return api.commonDivide(in, 5)
}

func (api *API) commonDivide(in int, divisor int) bool {
	result, err := api.divide(in, divisor)
	if err != nil {
		slog.Error("Error calling divide API", slog.String("error", err.Error()))
		return false
	}
	return result == 0
}

// divide calls the internal httptest server to perform a division operation, simulating
// an external HTTP API call. Any errors returned from the server are logged and returned to the caller.
func (api *API) divide(a, b int) (int, error) {

	// Construct the API URL
	url := fmt.Sprintf(requestPath, api.server.URL, a, b)

	// Submit the HTTP GET request to the server
	resp, err := api.server.Client().Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to call API: %w", err)
	}
	defer resp.Body.Close()

	// Check the size of the response before we slurp it all into memory.
	if resp.ContentLength > maxResponseSize {
		return 0, fmt.Errorf("response too large: %d bytes", resp.ContentLength)
	}
	if resp.ContentLength < 0 {
		// Note: unable to test this scenario without intercepting the response within this function.
		return 0, fmt.Errorf("response size unknown, will not process")
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body on response status %s: %w", resp.Status, err)
	}

	if resp.StatusCode == http.StatusOK {
		// Decode the JSON response
		var result server.DivisionResult
		if err := json.Unmarshal(bodyBytes, &result); err != nil {
			return 0, fmt.Errorf("failed to decode result response: %w", err)
		}
		return result.Remainder, nil
	}

	// Handle error responses
	var er server.ErrorResult
	if err := json.Unmarshal(bodyBytes, &er); err != nil {
		return 0, fmt.Errorf("failed to decode error response for status code %d: %w", resp.StatusCode, err)
	}

	switch {
	case resp.StatusCode == http.StatusBadRequest:
		return 0, fmt.Errorf("%s: %s", resp.Status, er.Message)

	case resp.StatusCode >= 400 && resp.StatusCode <= 499:
		return 0, fmt.Errorf("%s: %s", resp.Status, er.Message)

	case resp.StatusCode == http.StatusInternalServerError:
		return 0, fmt.Errorf("%s: %s", resp.Status, er.Message)

	default:
		return 0, fmt.Errorf("unexpected status code: %s: %s", resp.Status, er.Message)
	}
}
