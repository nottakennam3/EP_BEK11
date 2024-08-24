package main

import (
	"fmt"
	"os"
	"strconv"
)

func ex1(x, y int) (int, int) {
	return 2 * (x + y), x * y
}

func main(){
	args := os.Args
	if len(args) > 3 {
		fmt.Println("Must input only 2 numbers")
		return
	}
	length, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	width, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println(err)
		return
	}
	p, s := ex1(length, width)
	fmt.Printf("Perimeter: %v\nArea: %v\n", p, s)
}