package math

type Math struct{}

// Fizz implements the repository.FizzBuzzer interface.
func (Math) Fizz(in int) bool {
	return in%3 == 0
}

// Buzz implements the repository.FizzBuzzer interface.
func (Math) Buzz(in int) bool {
	return in%5 == 0
}
