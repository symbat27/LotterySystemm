package utils

import (
	"math/rand"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateWinningNumbers() []int {
	numbers := make(map[int]bool)
	result := make([]int, 0, 6)

	for len(result) < 6 {
		num := rng.Intn(49) + 1 // 1-49
		if !numbers[num] {
			numbers[num] = true
			result = append(result, num)
		}
	}

	return result
}

func ValidateNumbers(numbers []int) bool {
	if len(numbers) != 6 {
		return false
	}

	seen := make(map[int]bool)
	for _, num := range numbers {
		if num < 1 || num > 49 {
			return false
		}
		if seen[num] {
			return false // duplicate
		}
		seen[num] = true
	}

	return true
}

func CountMatches(ticketNumbers, winningNumbers []int) int {
	winMap := make(map[int]bool)
	for _, num := range winningNumbers {
		winMap[num] = true
	}

	matches := 0
	for _, num := range ticketNumbers {
		if winMap[num] {
			matches++
		}
	}

	return matches
}
