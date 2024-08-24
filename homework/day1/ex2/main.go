package main

import (
	"fmt"
	"os"
)

func ex2(s string) bool {
	return len(s) % 2 == 0
}

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Must input only 1 string")
		return
	}
	fmt.Println("Length of input string is divisible by 2:", ex2(os.Args[1]))
}