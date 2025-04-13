package app

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/MarkSonghurstPersonal/fizzbuzz-golang/pkg/repository"
)

type FizzBuzz struct {
	upperLimit int
	repo       repository.FizzBuzzer
}

type result struct {
	Number int
	Fizz   bool
	Buzz   bool
}

func New(upperLimit int, repo repository.FizzBuzzer) *FizzBuzz {
	return &FizzBuzz{
		upperLimit: upperLimit,
		repo:       repo,
	}
}

func (fb FizzBuzz) Run(ctx context.Context) {

	wg := sync.WaitGroup{}

	// The Channel buffer is limited to half the upper limit, to exercise the channel's blocking behavior.
	// This also ensures that the channel does not grow indefinitely if upperLimit is set to a very large number.
	chProcessor := make(chan result, fb.upperLimit/2)

	// Goroutine to generate FizzBuzz values and send them to the channel.
	go func() {
		defer close(chProcessor)
		for i := 1; i <= fb.upperLimit; i++ { // Fixed range loop to iterate correctly from 1 to upperLimit
			chProcessor <- result{
				Number: i,
				Fizz:   fb.repo.Fizz(i),
				Buzz:   fb.repo.Buzz(i),
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
