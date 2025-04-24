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
	for id, user := range q.Users {
		// Only randomize users who are still waiting
		if user.Status == "waiting" {
			userIDs = append(userIDs, id)
		}
	}

	// Shuffle logic
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(userIDs), func(i, j int) {
		userIDs[i], userIDs[j] = userIDs[j], userIDs[i]
	})

	// Assign positions and mark as queued
	for pos, id := range userIDs {
		q.Users[id].Status = "queued"
		q.Users[id].Position = pos + 1
	}

	return userIDs
}

// users are handled one by one
func (q *Queue) ProcessQueues() {
	orderedUsers := q.RandomizeQueue()

	go func() {
		for _, userID := range orderedUsers {
			q.Mu.Lock()
			user := q.Users[userID]

			if user.Status == "queued" {
				user.Status = "active"
				q.Mu.Unlock()

				time.Sleep(1 * time.Minute)

				q.Mu.Lock()
				user.Status = "expired"
				q.Mu.Unlock()
			} else {
				q.Mu.Unlock()
			}
		}
	}()
}
