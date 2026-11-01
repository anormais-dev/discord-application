package model

import "time"

// Cancellation diz quem cancelou a roleta por último, para a sala mostrar o aviso.
type Cancellation struct {
	ByID   string
	ByName string
	At     time.Time
}
