package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
)

type Person struct {
	name	string
	job		string
	yob		int
}

func main() {
	people := readFile("a.txt")
	fmt.Printf("Result: %v\n", people)
}

func readFile(fname string) []Person {
	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	res := []Person{}
	for scanner.Scan() {
		els := strings.Split(scanner.Text(), "|")
		res = append(res, createPerson(els))
	}
	return res
}

func createPerson(els []string) Person {
	name, job, yob := els[0], els[1], els[2]
	birthYear, err := strconv.Atoi(yob)
	if err != nil {
		log.Fatal(err)
	}
	return Person{
		name: strings.ToUpper(name), 
		job: strings.ToLower(job), 
		yob: birthYear,
	}
}