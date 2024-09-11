package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Employee struct {
	ID				int		`json:"id"`
	EmployeeName	string	`json:"employee_name"`
	EmployeeSalary	int		`json:"employee_salary"`
	EmployeeAge		int		`json:"employee_age"`
	ProfileImage	string	`json:"profile_image"`
}

type APIResponse struct {
	Status	string		`json:"status"`
	Data	[]Employee	`json:"data"`
	Message	string		`json:"message"`
}

const (
	url = "https://dummy.restapiexample.com/api/v1/employees"
	MAX_WORKER = 3
)

func main() {
	employeeData, err := fetchEmployeeData(url)
	if err != nil {
		log.Fatal(err)
	}
	results := getResults(employeeData)
	fmt.Println(results)
}

func fetchEmployeeData(url string) ([]Employee, error){
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong response status %d", resp.StatusCode)
	}
	var res APIResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func getResults(employees []Employee) []float64 {
	jobChan := make(chan *Employee, len(employees))
	resChan := make(chan float64, len(employees))
	for w := 0; w < MAX_WORKER; w++ {
		go worker(jobChan, resChan)
	}
	for _, e := range employees {
		jobChan <- &e
	}
    close(jobChan)
	res := []float64{}
    for a := 0; a < len(employees); a++ {
		res = append(res, <-resChan)
	}

	return res
}

func worker(job <-chan *Employee, res chan<- float64) {
	for e := range job {
		a := float64(e.EmployeeSalary) / float64(e.EmployeeAge)
		res <- float64(int(a * 100)) / 100
	}
}