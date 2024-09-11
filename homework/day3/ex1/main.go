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

func main() {
	url := "https://dummy.restapiexample.com/api/v1/employees"
	employeeData, err := fetchEmployeeData(url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", employeeData)
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