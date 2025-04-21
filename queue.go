package main

import (
	"math/rand"
	"sync"
	"time"
)

type Queue struct {
	Users map[string]*User
	Mu    sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		Users: make(map[string]*User),
	}
}

// adds user to waiting room
func (q *Queue) AddUser() string {
	q.Mu.Lock()
	defer q.Mu.Unlock()

	user := NewUser()
	q.Users[user.ID] = user

	return user.ID
}

// shufffling users and assigns position
func (q *Queue) RandomizeQueue() []string {
	q.Mu.Lock()
	defer q.Mu.Unlock()

	var userIDs []string
	for id := range q.Users {
		userIDs = append(userIDs, id)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(userIDs), func(i, j int) {
		userIDs[i], userIDs[j] = userIDs[j], userIDs[i]
	})

	for pos, id := range userIDs {
		q.Users[id].Status = "queued"
		q.Users[id].Position = pos + 1
	}

	return userIDs
}

// users are handled one by one
func (q *Queue) ProcessQueues() {
	orderedUsers := q.RandomizeQueue()
	
	for _, id := range orderedUsers {
		q.Mu.Lock()

		user := q.Users[id]
		if user.Status == "queued" {
			user.Status = "active"
			q.Mu.Unlock()

			time.Sleep(1*time.Minute)
			q.Mu.Lock()

			user.Status = "expired"
		}

		q.Mu.Unlock()
	}
}