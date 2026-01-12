package testCases

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = []User{
	{ID: 1, Name: "abc", Email: "abc@gmail.com"},
}
var idCounter = 2

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, `{"Error":"Invalid Body"}`, http.StatusBadRequest)
		return
	}

	if !validation(newUser, w) {
		return
	}
	newUser.ID = idCounter
	idCounter++

	users = append(users, newUser)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(newUser)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idstr := r.URL.Query().Get("id")
	if idstr != "" {

		id, err := strconv.Atoi(idstr)
		if err != nil {
			http.Error(w, `{"Error":"Invalid ID"}`, http.StatusBadRequest)
		}

		for _, user := range users {
			if user.ID == id {
				json.NewEncoder(w).Encode(user)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, `{"Error":"Invalid ID"}`, http.StatusBadRequest)
		return
	}

	var updatedUser User
	err1 := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err1 != nil {
		http.Error(w, `{"Error":"Invalid Body"}`, http.StatusBadRequest)
		return
	}

	if !validation(updatedUser, w) {
		return
	}

	for i, user := range users {
		if user.ID == id {
			users[i].Name = updatedUser.Name
			users[i].Email = updatedUser.Email
			json.NewEncoder(w).Encode(users[i])
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, `{"Error":"Invalid ID"}`, http.StatusBadRequest)
	}

	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message":"user %d deleted"}`, id)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		createUser(w, r)

	case http.MethodGet:
		getUser(w, r)

	case http.MethodPut:
		updateUser(w, r)

	case http.MethodDelete:
		deleteUser(w, r)

	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func validation(user User, w http.ResponseWriter) bool {
	w.Header().Set("Content-Type", "application/json")

	if strings.TrimSpace(user.Name) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"Error": "Name cannot be empty."})
		return false
	}

	if strings.TrimSpace(user.Email) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"Error": "Email cannot be empty."})
		return false
	}

	if !strings.HasSuffix(user.Email, "@gmail.com") {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"Error": "Wrong Email format."})
		return false
	}

	prefix := strings.TrimSuffix(user.Email, "@gmail.com")
	if prefix == "" {
		http.Error(w, `{"Error:mail cannot be empty"}`, http.StatusBadRequest)
		return false
	}
	return true
}

func UserManagement() {

	http.HandleFunc("/users/", userHandler)
	fmt.Println("Server running on Path 8080.")
	http.ListenAndServe(":8080", nil)
}
