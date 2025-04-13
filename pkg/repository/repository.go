package repository

// FizzBuzzer is an interface for implementors of the super complicated FizzBuzz algorithm.
type FizzBuzzer interface {
	// Fizz checks if the given number is divisible by 3 and if so, return false.
	Fizz(int) bool
	// Buzz checks if the given number is divisible by 5, and if so return true.
	Buzz(int) bool
}
