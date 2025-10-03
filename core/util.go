package core

import (
	"math/rand"
)

func randomKFromN(k, n int) []int {
	if k > n {
		panic("k greater than n")
	}

	nums := make([]int, n)
	for i := range n {
		nums[i] = i + 1
	}

	rand.Shuffle(n, func(i, j int) {
		nums[i], nums[j] = nums[j], nums[i]
	})

	return nums[:k]
}
