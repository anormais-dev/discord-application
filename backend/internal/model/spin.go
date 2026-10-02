package model

import "time"

// Spin tem tudo que os clientes precisam para animar a roda e parar todos na mesma posição.
type Spin struct {
	WinnerID     string
	PoolSnapshot []string
	TargetOffset float64 // 0..1 dentro da fatia do vencedor
	Rotations    int
	StartedAt    time.Time
	DurationMs   int
}
