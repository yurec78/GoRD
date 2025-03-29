package main

import (
	"fmt"
	"strconv"
)

func main() {

	//test for Fibonacci functions
	numbers := []int{1, 2, 3, 4, 5, 0, 9, 13, 34}
	for i := 0; i < len(numbers); i++ {
		fmt.Println(isBinaryPalindrome(numbers[i]))
	}
}

func isBinaryPalindrome(n int) (bool, string) {
	result := true
	binaryString := strconv.FormatInt(int64(n), 2)
	for i, j := 0, len(binaryString)-1; i < j; i, j = i+1, j-1 {
		if binaryString[i] != binaryString[j] {
			result = false
		}
	}
	return result, binaryString
}
