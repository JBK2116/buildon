package model

import "time"

// Problem represents a singular problem within a project.
type Problem struct {
	title     string
	content   string
	createdAt time.Time
	updatedAt time.Time
}
