package main

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID       string    `json:"id"`
	Status   string    `json:"status"` //waiting, queued, active, expired
	JoinedAt time.Time `json:"joined_at"`
	Position int       `json:"position"` //asigned after randomization
}

func NewUser() *User {
	return &User{
		ID:       uuid.New().String(),
		Status:   "waiting",
		JoinedAt: time.Now(),
	}
}
