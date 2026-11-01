package model

import "time"

type Participant struct {
	ID       string
	Name     string
	Avatar   string
	Online   bool
	JoinedAt time.Time
	// ReconnectUntil é o prazo para voltar sem perder o lugar; zero enquanto online.
	ReconnectUntil time.Time
}
