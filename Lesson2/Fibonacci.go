package main

import (
	"fmt"
)

func main() {

	//test for Fibonacci functions
	numbers := []int{1, 2, 3, 4, 5, 0, 9, 13, 34}
	for i := 0; i < len(numbers); i++ {
		fmt.Println(fibonacciRecursive(numbers[i]))
		fmt.Println(fibonacciIterative(numbers[i]))
	}
}

func fibonacciRecursive(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacciRecursive(n-1) + fibonacciRecursive(n-2)
}

func fibonacciIterative(n int) int {
	if n <= 1 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}
