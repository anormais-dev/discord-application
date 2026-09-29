// Package roulette contém a lógica da roleta de times, sem dependência de rede.
// É usado pela Activity e, no futuro, pelo bot.
package roulette

import (
	"errors"
	"time"
)

type Phase string

const (
	PhaseLobby    Phase = "lobby"
	PhaseDrafting Phase = "drafting"
	PhaseFinished Phase = "finished"
)

var (
	ErrNotAdmin         = errors.New("apenas o admin pode fazer isso")
	ErrCannotSpin       = errors.New("você não tem permissão para girar")
	ErrWrongPhase       = errors.New("ação não permitida neste momento")
	ErrAlreadyJoined    = errors.New("você já está participando")
	ErrNotParticipant   = errors.New("você não está participando")
	ErrSpinInProgress   = errors.New("a roleta já está girando")
	ErrNotEnoughPlayers = errors.New("participantes insuficientes para o formato")
	ErrInvalidFormat    = errors.New("formato inválido")
	ErrNoSpinInProgress = errors.New("nenhum giro em andamento")
	ErrNotImplemented   = errors.New("não implementado")
)

type Participant struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Avatar   string    `json:"avatar"`
	Online   bool      `json:"online"`
	JoinedAt time.Time `json:"joinedAt"`
}

// SpinState descreve um giro em andamento. Os clientes animam a roda
// a partir destes dados, então todos param na mesma posição.
type SpinState struct {
	WinnerID     string    `json:"winnerId"`
	PoolSnapshot []string  `json:"poolSnapshot"`
	TargetOffset float64   `json:"targetOffset"` // 0..1 dentro da fatia do vencedor
	Rotations    int       `json:"rotations"`
	StartedAt    time.Time `json:"startedAt"`
	DurationMs   int       `json:"durationMs"`
}

type Room struct {
	Participants []Participant   `json:"participants"`
	AdminID      string          `json:"adminId"`
	Spinners     map[string]bool `json:"spinners"`
	Format       Format          `json:"format"`
	Pool         []string        `json:"pool"`
	Teams        [][]string      `json:"teams"`
	Phase        Phase           `json:"phase"`
	CurrentSpin  *SpinState      `json:"spin"`
	Picks        int             `json:"picks"`
}

func NewRoom() *Room {
	return &Room{
		Spinners: map[string]bool{},
		Format:   DefaultFormat,
		Phase:    PhaseLobby,
	}
}

func (r *Room) Join(p Participant) error                        { return ErrNotImplemented }
func (r *Room) Leave(userID string) error                       { return ErrNotImplemented }
func (r *Room) SetFormat(userID string, f Format) error         { return ErrNotImplemented }
func (r *Room) SetSpinner(userID, target string, on bool) error { return ErrNotImplemented }
func (r *Room) Spin(userID string) (*SpinState, error)          { return nil, ErrNotImplemented }
func (r *Room) CommitSpin() error                               { return ErrNotImplemented }
func (r *Room) Reset(userID string) error                       { return ErrNotImplemented }
func (r *Room) Disconnect(userID string)                        {}
