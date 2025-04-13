# fizzbuzz-golang

fizzbuzz-golang provides a Fizz Buzz command line (CLI) program.

The program only uses Golang standard packages (the go.mod is empty).

It's an opportunity to exercise my Golang knowledge and not necessarily the most efficient Fizz Buzz implementation possible!

Compared to my [Z80 Fizz Buzz](https://github.com/MarkSonghurstPersonal/fizzbuzz-z80) program, this took a fraction of the time to write! Ah, but it lacks the Spectrum's beeps!

*Mark Songhurst, April 2025*


## Repository Structure
This git repo follows the [Golang Standard for Project Layout](https://github.com/golang-standards/project-layout)

In short:

* `cmd/fizzbuzz` contains the main function for the CLI program, where it performs argument parsing.
* `pkg/repository` Contains the FizzBuzz interface (repository) which adapters must implement.
* `internal/app` Contains the business logic and mechanics of the program. You could use this even if the program was not CLI based.
* `internal/adapters/secondary` Contains implementations of the FizzBuzzer interface
    * `internal/adapters/secondary/math` Is a simple math based implementor. Arguably this is not a secondary adapter as it doesn't call out to anything external.
    * `internal/adapters/secondary/httpapi` Simulates an HTTP REST API which provides a divide endpoint. I've use httptest.Server to provide a local HTTP service.


## Testing
Almost all packages have unit tests, run this from the top level directory to exercise the code:
```bash
$ go test ./... --cover --race
```

## Building
I'll add a makefile at some point, but for now:
```bash
$ cd cmd/fizzbuzz
$ go build
```

It should build on any Golang support system.

## Golang Features Covered
* Channels and go routines.
* Contexts, including cancelation.
* Error wrapping.
* Interfaces.
* Logging using log/slog
* Unit testing, including parameterized tests.


## TODO
* Squeeze in Generics somewhere. Perhaps by adding support for float64 division?