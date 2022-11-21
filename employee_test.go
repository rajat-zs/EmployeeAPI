package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDepartmentHandler(t *testing.T) {
	tt := []struct {
		name       string
		method     string
		input      *Departments
		want       string
		statusCode int
	}{
		{
			name:       "without Departments",
			method:     http.MethodGet,
			input:      &Departments{},
			want:       "Error: No department found",
			statusCode: http.StatusNotFound,
		},
		{
			name:   "with Departments",
			method: http.MethodGet,
			input: &Departments{
				Department{
					ID:   2,
					Name: "Technology",
				},
			},
			want:       `[{"id":2,"name":"Technology"}]`,
			statusCode: http.StatusOK,
		},
		{
			name:       "with bad method",
			method:     http.MethodPost,
			input:      &Departments{},
			want:       "Method not allowed",
			statusCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(tc.method, "/employees", nil)
			responseRecorder := httptest.NewRecorder()

			DepartmentHandler{tc.input}.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			if strings.TrimSpace(responseRecorder.Body.String()) != tc.want {
				t.Errorf("Want '%s', got '%s'", tc.want, responseRecorder.Body)
			}
		})
	}
}

func TestEmployeeHandler(t *testing.T) {
	tt := []struct {
		name        string
		method      string
		departments *Departments
		employees   *Employees
		body        string
		want        string
		statusCode  int
	}{
		{
			name:   "with a department ID and name",
			method: http.MethodPost,
			departments: &Departments{
				Department{
					ID:   1,
					Name: "HR",
				},
			},
			employees:  &Employees{},
			body:       `{"employee_depart_id": 1,"name": "Raja"}`,
			want:       `{"employee_depart_id":1,"name":"Raja","age":0,"address":""}`,
			statusCode: http.StatusOK,
		},
		{
			name:   "with a department and no name",
			method: http.MethodPost,
			departments: &Departments{
				Department{
					ID:   1,
					Name: "HR",
				},
			},
			employees:  &Employees{},
			body:       `{"employee_depart_id": 1}`,
			want:       `{"employee_depart_id":1,"name":"","age":0,"address":""}`,
			statusCode: http.StatusOK,
		},
		{
			name:        "with no departments on list",
			method:      http.MethodPost,
			departments: &Departments{},
			employees:   &Employees{},
			body:        `{"employee_depart_id":1,"name":"Raja","age":23,"address":"ward 2"}`,
			want:        "Error: No department found",
			statusCode:  http.StatusNotFound,
		},
		{
			name:        "with GET method and no employee in memory",
			method:      http.MethodGet,
			departments: &Departments{},
			employees:   &Employees{},
			body:        "",
			want:        "[]",
			statusCode:  http.StatusOK,
		},
		{
			name:        "with GET method and with employees already in memory",
			method:      http.MethodGet,
			departments: &Departments{},
			employees: &Employees{
				Employee{
					EmpDepartID: 1,
					Name:        "Raja",
					Age:         23,
					Address:     "ward 2",
				},
			},
			body:       "",
			want:       `[{"employee_depart_id":1,"name":"Raja","age":23,"address":"ward 2"}]`,
			statusCode: http.StatusOK,
		},
		{
			name:        "with bad HTTP method",
			method:      http.MethodDelete,
			departments: &Departments{},
			employees:   &Employees{},
			body:        "",
			want:        "Method not allowed",
			statusCode:  http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(tc.method, "http://localhost:8082/employees", strings.NewReader(tc.body))
			responseRecorder := httptest.NewRecorder()

			handler := employeeHandler{tc.departments, tc.employees}
			handler.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}
			if strings.TrimSpace(responseRecorder.Body.String()) != tc.want {
				t.Errorf("Want '%s', got '%s'", tc.want, responseRecorder.Body)
			}
		})
	}
}
