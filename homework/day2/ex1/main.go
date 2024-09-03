package main

import (
	"fmt"
	"time"
)

type Person struct{
	name	string
	job		string
	yob		int
}

func main() {
	p := createPerson("Nghia", "Software Developer", 2000)
	fmt.Printf("Name: %s\nAge:  %d\nJob:  %s\n", p.name, p.getAge(), p.job)
	fmt.Println("Is compatible:", p.isCompatible())
}

func createPerson(name, job string, yob int) Person{
	return Person{name: name, job: job, yob: yob}
}

func (p *Person) getAge() int {
	return time.Now().Year() - p.yob
}

func (p *Person) isCompatible() bool {
	return p.yob % len(p.name) == 0
}