package main

import (
	"fmt"
	"strconv"
)

func main() {

	testCases := []struct {
		input    string
		expected int
	}{
		{"0", 1},
		{"1", 2},
		{"10", 3},
		{"11", 4},
		{"100", 5},
		{"101", 6},
		{"110", 7},
		{"111", 8},
		{"1000", 9},
		{"1111", 16},
		{"101010", 43},
		{"111111", 64},
	}

	for _, tc := range testCases {
		result := Increment(tc.input)
		if result != tc.expected {
			fmt.Printf("Помилка: Increment(%q) = %d, очікувалося %d\n", tc.input, result, tc.expected)
		} else {
			fmt.Printf("Успішно: Increment(%q) = %d\n", tc.input, result)
		}
	}
}

func Increment(num string) int {
	// 1. Перетворюємо двійковий рядок в ціле число
	intValue, _ := strconv.ParseInt(num, 2, 64)

	// 2. Додаємо одиницю
	intValue++

	// 3. Перетворюємо результат назад в двійковий рядок
	binaryString := strconv.FormatInt(intValue, 2)

	// 4. Перетворюємо двійковий рядок в int та повертаємо результат.
	result, _ := strconv.Atoi(binaryString)

	return result
}
