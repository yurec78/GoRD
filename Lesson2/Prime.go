package main

import (
	"fmt"
	"math"
)

func main() {

	numbers := []int{1, 2, 3, 4, 5, 0, 9, 13, 34, 7}

	//test for simple number status
	for i := 0; i < len(numbers); i++ {

		fmt.Println(numbers[i], isPrime(numbers[i]))
	}
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
