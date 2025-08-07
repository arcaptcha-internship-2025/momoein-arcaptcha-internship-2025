package main

import (
	"fmt"
	"math"
)

func main() {
	var input int
	fmt.Scan(&input)

	results := DiehardPasswords(input)
	for _, num := range results {
		fmt.Println(num)
	}
}

func DiehardPasswords(n int) []int {
	var results []int

	var dfs func(current int, length int)
	dfs = func(current int, length int) {
		if length == n {
			results = append(results, current)
			return
		}

		for digit := 0; digit <= 9; digit++ {
			next := current*10 + digit
			if IsPrime(next) {
				dfs(next, length+1)
			}
		}
	}

	for i := 2; i <= 9; i++ {
		if IsPrime(i) {
			dfs(i, 1)
		}
	}

	return results
}

func IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	sqrt := int(math.Sqrt(float64(n)))
	for i := 3; i <= sqrt; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}
