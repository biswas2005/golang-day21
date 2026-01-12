// Table-Driven test cases for CRUD operations in User Management
package testCases

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test_createUser() verifies the conditions for createUser Handler
func Test_createUser(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		status int
	}{
		//Todo's for the test case
		{
			//Happy Path: Valid user should be Accepted
			name:   "Valid User",
			body:   `{"name":"abc","email":"abc@gmail.com"}`,
			status: http.StatusAccepted,
		},
		{
			name:   "Empty Name",
			body:   `{"name":"","email":"abc@gmail.com"}`,
			status: http.StatusBadRequest,
		},
		{
			name:   "Empty Email",
			body:   `{"name":"abc","email":""}`,
			status: http.StatusBadRequest,
		},
		{
			name:   "Invalid Name",
			body:   `{"name":123,"email":"abc@gmail.com"}`,
			status: http.StatusBadRequest,
		},
		{
			name:   "Invalid Email",
			body:   `{"name":"abc","email":"@gmail.com"}`,
			status: http.StatusBadRequest,
		},
		{
			name:   "Mail format invalid",
			body:   `{"name":"abc","email":"abc"}`,
			status: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Dummy request
			req := httptest.NewRequest(http.MethodPost, "/users/", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()
			createUser(w, req)

			res := w.Result()
			if res.StatusCode != tt.status {
				t.Errorf("got %d,want %d", res.StatusCode, tt.status)
			}
		})
	}
}

// Test_getUser() verifies the condition for getUser Handler
func Test_getUser(t *testing.T) {
	users = []User{
		{ID: 1, Name: "abc", Email: "abc@gmail.com"},
	}

	tests := []struct {
		name       string
		query      string
		wantStatus int
		wantBody   string
	}{
		//Todo's for the test
		{
			//Happy Path: Validates the condition
			name:       "Get all users",
			query:      "",
			wantStatus: http.StatusOK,
			wantBody:   `[{"id":1,"name":"abc","email":"abc@gmail.com"}]`,
		},
		{
			// Happy Path: Validates the condition for specific ID
			name:       "Get user by id",
			query:      "?id=1",
			wantStatus: http.StatusOK,
			wantBody:   `{"id":1,"name":"abc","email":"abc@gmail.com"}`,
		},
		{
			name:       "Invalid ID format",
			query:      "?id=abc",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"Error":"Invalid ID"}`,
		},
		{
			name:       "Non existing ID",
			query:      "?id=999",
			wantStatus: http.StatusNotFound,
			wantBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/users"+tt.query, nil)
			w := httptest.NewRecorder()

			getUser(w, r)

			res := w.Result()
			if res.StatusCode != tt.wantStatus {
				t.Errorf("Got %d ,want %d", res.StatusCode, tt.wantStatus)
			}

			body, _ := io.ReadAll(res.Body)

			if strings.TrimSpace(string(body)) != strings.TrimSpace(tt.wantBody) {
				t.Errorf("Got body %s, want body %s", body, tt.wantBody)
			}
		})
	}
}

// Test_updateUser() verifies the condition for updateUser Handler
func Test_updateUser(t *testing.T) {

	users = []User{
		{ID: 1, Name: "abc", Email: "abc@gmail.com"},
	}
	tests := []struct {
		name       string
		url        string
		body       string
		wantStatus int
		wantBody   string
	}{
		//Todo's for the test cases
		{
			name:       "Invalid Path",
			url:        "/users/",
			body:       `{"id":1,"name":"new","email":"new@gmail.com"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Invalid ID Format",
			url:        "/users/abc",
			body:       `{"id":1,"name":"new","email":"new@gmail.com"}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"Error":"Invalid ID"}`,
		},
		{
			name:       "Invalid JSON body",
			url:        "/users/1",
			body:       `Invalid JSON`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"Error":"Invalid Body"}`,
		},
		{
			name:       "validation Fails-Invalid Name",
			url:        "/users/1",
			body:       `{"id":1,"name":"","email":"new@gmail.com"}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"Error":"Name cannot be empty."}`,
		},
		{
			name:       "Invalid ID-range missing",
			url:        "/users/99",
			body:       `{"id":99,"name":"fgd","email":"fgd@gmail.com"}`,
			wantStatus: http.StatusNotFound,
		},
		{
			//Happy Path:Validates the condition
			name:       "Successful Update",
			url:        "/users/1",
			body:       `{"id":1,"name":"updated abc","email":"abcnew@gmail.com"}`,
			wantStatus: http.StatusOK,
			wantBody:   `{"id":1,"name":"updated abc","email":"abcnew@gmail.com"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, tt.url, bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			updateUser(w, req)
			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatus {
				t.Errorf("Expected status %d, got status %d", tt.wantStatus, res.StatusCode)
			}
			if tt.wantBody != "" {
				gotBody, _ := io.ReadAll(res.Body)
				got := strings.TrimSpace(string(gotBody))
				if got != tt.wantBody {
					t.Errorf("Wanted body %s, got body %s", tt.wantBody, got)
				}
			}

		})
	}
}

// Test_deleteUser() verifies the condition for deleteUser Handler
func Test_deleteUser(t *testing.T) {
	users = []User{
		{ID: 1, Name: "abc", Email: "abc@gmail.com"},
	}
	tests := []struct {
		name       string
		url        string
		body       string
		wantStatus int
		wantBody   string
	}{
		//Todo's for the test case
		{
			name:       "Invalid Path",
			url:        "/users",
			body:       `{"id":1,"name":"abc","email":"abc@gmail.com"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Invalid ID Format",
			url:        "/users/abc",
			body:       `{"id":1,"name":"abc","email":"abc@gmail.com"}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"Error":"Invalid ID"}`,
		},
		{
			name:       "User not Found",
			url:        "/users/99",
			body:       `{"id":1,"name":"abc","email":"abc@gmail.com"}`,
			wantStatus: http.StatusNotFound,
		},
		{
			//Happy path:Validates the condition
			name:       "Successful Delete",
			url:        "/users/1",
			body:       `{"id":1,"name":"abc","email":"abc@gmail.com"}`,
			wantStatus: http.StatusOK,
			wantBody:   `{"message":"user 1 deleted"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, tt.url, nil)
			w := httptest.NewRecorder()
			deleteUser(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatus {
				t.Errorf("Wanted status %d , Got status %d", tt.wantStatus, res.StatusCode)
			}
			if tt.wantBody != "" {
				gotBody, _ := io.ReadAll(res.Body)
				got := strings.TrimSpace(string(gotBody))
				if got != tt.wantBody {
					t.Errorf("Wanted body %s, Got body %s", tt.wantBody, got)
				}
			}
		})
	}
}
