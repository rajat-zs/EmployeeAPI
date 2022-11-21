package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Department struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Departments []Department

func (ds Departments) FindByID(ID int) (Department, error) {
	for _, department := range ds {
		if department.ID == ID {
			return department, nil
		}
	}

	return Department{}, fmt.Errorf("couldn't find department with ID: %d", ID)
}

type DepartmentHandler struct {
	departments *Departments
}

func (ds DepartmentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		if len(*ds.departments) == 0 {
			http.Error(w, "Error: No department found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(ds.departments)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

type Employee struct {
	EmpDepartID int    `json:"employee_depart_id"`
	Name        string `json:"name"`
	Age         int    `json:"age"`
	Address     string `json:"address"`
}

type Employees []Employee

type employeeHandler struct {
	departments *Departments
	employees   *Employees
}

func (eh employeeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		var o Employee

		if len(*eh.departments) == 0 {
			http.Error(w, "Error: No department found", http.StatusNotFound)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&o)
		if err != nil {
			http.Error(w, "Can't decode body", http.StatusBadRequest)
			return
		}

		d, err := eh.departments.FindByID(o.EmpDepartID)
		fmt.Println("Your department is", d)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusBadRequest)
			return
		}
		*eh.employees = append(*eh.employees, o)
		json.NewEncoder(w).Encode(o)
	case http.MethodGet:
		json.NewEncoder(w).Encode(eh.employees)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	var employees Employees
	departments := Departments{
		Department{
			ID:   1,
			Name: "HR",
		},
		Department{
			ID:   2,
			Name: "Technology",
		},
		Department{
			ID:   3,
			Name: "Testing",
		},
	}

	mux := http.NewServeMux()
	mux.Handle("/departments", DepartmentHandler{&departments})
	mux.Handle("/employees", employeeHandler{&departments, &employees})

	log.Fatal(http.ListenAndServe(":8082", mux))
}
