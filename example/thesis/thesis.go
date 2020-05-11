package thesis

import "time"

// Thesis model
type Thesis struct {
	Title     string
	Body      string
	Author    string
	CreatedAt time.Time
}
