package model

import "time"

type Participant struct {
	ID       string
	Name     string
	Avatar   string
	Online   bool
	JoinedAt time.Time
}
