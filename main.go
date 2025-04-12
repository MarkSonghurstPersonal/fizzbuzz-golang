package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/MarkSonghurstPersonal/fizzbuzz-golang/pkg/internal/adapters/secondary/math"
	"github.com/MarkSonghurstPersonal/fizzbuzz-golang/pkg/internal/repository"
)

type fizzBuzz struct {
	Number int
	Fizz   bool
	Buzz   bool
}

func main() {

	// Use slog for structured logging in JSON format.
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		// Filter out the automatically added time field.
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})
	slog.SetDefault(slog.New(handler))

	var upperLimit int
	flag.IntVar(&upperLimit, "limit", 64, "The upper limit for the FizzBuzz sequence, defaults to 64, must be higher than 0")
	var adapterToUse string
	flag.StringVar(&adapterToUse, "adapter", "math", "The adapter to use for FizzBuzz logic, defaults to math")
	flag.Parse()

	if upperLimit < 1 {
		slog.Error("Invalid upper limit", slog.Int("upperLimit", upperLimit))
		os.Exit(1)
	}

	var repo repository.FizzBuzzer
	switch adapterToUse {
	case "math":
		repo = math.FizzBuzzer{}

	default:
		slog.Error("Invalid adapter", slog.String("adapter", adapterToUse))
		os.Exit(1)
	}

	slog.Info("Starting FizzBuzz", slog.Int("upperLimit", upperLimit), slog.String("adapter", adapterToUse))

	wg := sync.WaitGroup{}

	// The Channel buffer is limited to half the upper limit, to exercise the channel's blocking behavior.
	// This is to ensure that the channel does not grow indefinitely.
	chProcessor := make(chan fizzBuzz, upperLimit/2)

	// Goroutine to generate FizzBuzz values and send them to the channel.
	go func() {
		defer close(chProcessor)
		for i := 1; i <= upperLimit; i++ { // Fixed range loop to iterate correctly from 1 to upperLimit
			chProcessor <- fizzBuzz{
				Number: i,
				Fizz:   repo.Fizz(i),
				Buzz:   repo.Buzz(i),
			}
		}
	}()

	// Goroutine to process FizzBuzz values from the channel.
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			fb, ok := <-chProcessor
			if !ok {
				slog.Debug("Channel closed")
				return
			}
			switch {
			case fb.Fizz && fb.Buzz:
				slog.Info("FizzBuzz", slog.Int("number", fb.Number))

			case fb.Fizz:
				slog.Info("Fizz", slog.Int("number", fb.Number))

			case fb.Buzz:
				slog.Info("Buzz", slog.Int("number", fb.Number))

			default:
				slog.Info(fmt.Sprintf("%d", fb.Number))
			}
		}
	}()

	wg.Wait()
}
