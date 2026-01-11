package testCases

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_createUser(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		status int
	}{
		{
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
		{
			name:       "Get all users",
			query:      "",
			wantStatus: http.StatusOK,
			wantBody:   `[{"id":1,"name":"abc","email":"abc@gmail.com"}]`,
		},
		{
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
