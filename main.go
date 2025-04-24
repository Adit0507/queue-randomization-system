package main

import (
	"encoding/json"
	"net/http"
)

var q = NewQueue()

func main() {
	http.HandleFunc("/join", joinHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/start", startHandler)
	http.HandleFunc("/", serveIndex	)

	http.HandleFunc("/users", usersPageHandler)
	http.HandleFunc("/users-data", usersDataHandler)

	http.ListenAndServe(":8080", nil)
}

// adds a user to a waiting room
func joinHandler(w http.ResponseWriter, r *http.Request) {
	userID := q.AddUser()
	json.NewEncoder(w).Encode(map[string]string{"user_id": userID})
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "./templates/index.html")
		return
	}
	if r.URL.Path == "/users" {
		http.ServeFile(w, r, "./templates/users.html")
		return
	}
	http.NotFound(w, r)
}

func usersPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./templates/users.html")
}

func usersDataHandler(w http.ResponseWriter, r *http.Request) {
	q.Mu.Lock()
	defer q.Mu.Unlock()

	var users []*User
	for _, u := range q.Users {
		users = append(users, u)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// returns user status
func statusHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	q.Mu.Lock()

	user, exists := q.Users[userID]
	q.Mu.Unlock()

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// triggers queue randomization and processing
func startHandler(w http.ResponseWriter, r *http.Request) {
	go q.ProcessQueues()
	w.Write([]byte("Queue Started!!"))
}