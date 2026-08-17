// Package model represents the data types used in this application.
package model

import "time"

// Project represents a singular project that contains problems.
type Project struct {
	ID        int
	Title     string
	Problems  []Problem
	CreatedAt time.Time
	UpdatedAt time.Time
}
