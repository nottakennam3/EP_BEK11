package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Must input at least 2 numbers")
		return
	}
	var sum int
	nums := make([]int, 0, 2)
	for _, arg := range os.Args[1:] {
		n, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Println("Cannot convert string to integer:", arg)
			return
		}
		nums = append(nums, n)
		sum += n
	}
	fmt.Printf("Sum: %v\nMax: %v\nMin: %v\nAvg: %.2f\n", sum, nums[len(nums)-1], nums[0], float64(sum) / float64(len(nums)))
}