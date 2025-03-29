package main

import (
	"fmt"
)

func main() {
	fmt.Println(ValidParentheses("()"))     // true
	fmt.Println(ValidParentheses("()[]{}")) // true
	fmt.Println(ValidParentheses("(]"))     // false
	fmt.Println(ValidParentheses("([)]"))   // false
	fmt.Println(ValidParentheses("{[]}"))   // true
	fmt.Println(ValidParentheses("[{}]"))   // true
	fmt.Println(ValidParentheses("[{]}"))   // false
	fmt.Println(ValidParentheses("((((((((((((((((((((((((((((((((((((((((((((((((("))
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
