package main

import (
	"fmt"
	"os"
)

func getStrMap(s string) map[string]int {
	strMap := make(map[string]int, len(s))
	for _, c := range s {
		strMap[string(c)] += 1
	}
	return strMap
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Must input 1 string")
		return
	}
	input := os.Args[1]
	strMap := getStrMap(input)
	fmt.Printf("Chars count: %v\n", strMap)
}