package main

import "log"

func main() {

	for i := range 100 {

		switch {
		case i%3 == 0 && i%5 == 0:
			log.Println("FizzBuzz")

		case i%3 == 0:
			log.Println("Fizz")

		case i%5 == 0:
			log.Println("Buzz")

		default:
			log.Println(i)
		}
	}
}
