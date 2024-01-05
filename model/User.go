package model

import "time"

type User struct {
	UserId              int
	Username            string
	Email               string
	FirstName           string
	LastName            string
	Created             time.Time
	LastLogin           time.Time
	HashedPassword      string
	FailedLoginAttempts int
	LoginBlockedUntil   time.Time
}
