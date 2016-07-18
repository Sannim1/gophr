package main

import "time"

// Image represents an image
type Image struct {
	ID          string
	UserID      string
	Name        string
	Location    string
	Size        string
	CreatedAt   time.Time
	Description string
}
