package math

type FizzBuzzer struct{}

func (FizzBuzzer) Fizz(in int) bool {
	return in%3 == 0
}

func (FizzBuzzer) Buzz(in int) bool {
	return in%5 == 0
}
