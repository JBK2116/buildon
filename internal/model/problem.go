package model

import "time"

// Problem represents a singular problem within a project.
type Problem struct {
	ID        int
	ProjectID int
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
