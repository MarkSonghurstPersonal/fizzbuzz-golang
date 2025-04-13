package app

import (
	"context"
	"sync"
	"testing"
)

// mockFizzBuzzer is a mock implementation of the FizzBuzzer interface.
type mockFizzBuzzer struct{}

func (m mockFizzBuzzer) Fizz(n int) bool {
	return n%3 == 0
}

func (m mockFizzBuzzer) Buzz(n int) bool {
	return n%5 == 0
}

func TestRun(t *testing.T) {
	tests := []struct {
		name       string
		upperLimit int
	}{
		{
			name:       "call Run with upperLimit 15",
			upperLimit: 15,
		},
		{
			name:       "Call Run with upperLimit 1",
			upperLimit: 1,
		},
		{
			name:       "Call Run with upperLimit 0 (no output expected)",
			upperLimit: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new FizzBuzz instance with the upper limit and a mock FizzBuzzer
			fb := New(tt.upperLimit, mockFizzBuzzer{})

			// We need to call the Run function to exercise the FizzBuzz logic in a separate goroutine,
			// this is required because the Run function starts it's own goroutine and then immediately returns.
			// Use a WaitGroup to wait until the Run function completes
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				fb.Run(context.Background())
			}()

			// Wait for the Run function to complete
			wg.Wait()

			// TODO
			// No test assertions are made here because the Run function just logs using log/slog we've got no
			// function output to validate.
			// However, we could redirect slog's output at the start of the test to an io.Writer such as a
			// bytes.Buffer and then validate the contents.
			/*
				var buf bytes.Buffer
				handler := slog.NewJSONHandler(&buf, nil)
				slog.SetDefault(slog.New(handler))
			*/
		})
	}
}
