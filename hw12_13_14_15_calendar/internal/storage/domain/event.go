package domain

import (
	"time"
)

type Event struct {
	ID          string
	Title       string
	DateTime    time.Time
	Duration    time.Duration
	Description string
	UserID      string
	NotifyIn    time.Time
}

type Notification struct {
	ID     string
	Title  string
	Date   time.Time
	UserID string
}
