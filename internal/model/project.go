package model

import "time"

// Project represents a singular project that contains problems.
type Project struct {
	title     string
	problems  []Problem
	createdAt time.Time
	updatedAt time.Time
}
