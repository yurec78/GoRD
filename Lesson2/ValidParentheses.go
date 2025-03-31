package main

import (
	"fmt"
)

func main() {
	fmt.Println(ValidParentheses("()"))     // true
	fmt.Println(ValidParentheses("()[]{}")) // true
	fmt.Println(ValidParentheses2("(]"))    // false
	fmt.Println(ValidParentheses2("([)]"))  // false
	fmt.Println(ValidParentheses2("{[]}"))  // true
	fmt.Println(ValidParentheses2("[{}]"))  // true
	fmt.Println(ValidParentheses2("[{]}"))  // false
	fmt.Println(ValidParentheses2("((((((((((((((((((((((((((((((((((((((((((((((((("))
}

func ValidParentheses(s string) bool {
	stack := []rune{}
	brackets := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}

	for _, char := range s {
		if char == '(' || char == '[' || char == '{' {
			stack = append(stack, char)
		} else if open, ok := brackets[char]; ok {
			if len(stack) == 0 || stack[len(stack)-1] != open {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}

func ValidParentheses2(s string) bool {
	stack := []rune{}
	brackets := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}

	for _, char := range s {
		switch char {
		case '(', '[', '{':
			stack = append(stack, char)
		case ')', ']', '}':
			open, ok := brackets[char]
			if !ok {
				return false // Invalid closing bracket
			}
			if len(stack) == 0 || stack[len(stack)-1] != open {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0
}
