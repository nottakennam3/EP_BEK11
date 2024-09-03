package main

import (
	"fmt"
)

func twoSum(nums []int, target int) []int {
	valueMap := make(map[int]int, len(nums))
	for i, n := range nums {
		need := target - n
		if val, ok := valueMap[need]; ok {
			return []int{val, i}
		}
		valueMap[n] = i
	}
	return []int{}
}

func main() {
	nums := []int{3, 5, 3}
	target := 6
	res := twoSum(nums, target)
	fmt.Println("Result:", res)
}