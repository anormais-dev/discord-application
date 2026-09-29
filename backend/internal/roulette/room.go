// Package roulette contém a lógica da roleta de times, sem dependência de rede.
// É usado pela Activity e, no futuro, pelo bot.
package roulette

import (
	crand "crypto/rand"
	"errors"
	"math/rand/v2"
	"slices"
	"time"
)

type Phase string

const (
	PhaseLobby    Phase = "lobby"
	PhaseDrafting Phase = "drafting"
	PhaseFinished Phase = "finished"
)

// SpinDuration é quanto tempo a animação da roda dura nos clientes.
const SpinDuration = 5 * time.Second

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
	Teams        [][]Participant `json:"teams"`
	Phase        Phase           `json:"phase"`
	CurrentSpin  *SpinState      `json:"spin"`
	Picks        int             `json:"picks"`

	rng *rand.Rand
}

func NewRoom() *Room {
	var seed [32]byte
	crand.Read(seed[:])
	return NewRoomWithRand(rand.New(rand.NewChaCha8(seed)))
}

func NewRoomWithRand(rng *rand.Rand) *Room {
	return &Room{
		Participants: []Participant{},
		Spinners:     map[string]bool{},
		Format:       DefaultFormat,
		Pool:         []string{},
		Teams:        [][]Participant{},
		Phase:        PhaseLobby,
		rng:          rng,
	}
}

// Join adiciona um participante. Só é permitido antes do primeiro giro.
func (r *Room) Join(p Participant) error {
	if r.Phase != PhaseLobby {
		return ErrWrongPhase
	}
	if r.indexOf(p.ID) >= 0 {
		return ErrAlreadyJoined
	}
	p.Online = true
	if p.JoinedAt.IsZero() {
		p.JoinedAt = time.Now()
	}
	r.Participants = append(r.Participants, p)
	r.Pool = append(r.Pool, p.ID)
	if r.AdminID == "" {
		r.AdminID = p.ID
	}
	return nil
}

// Leave remove o participante. Não é permitido enquanto os times estão sendo montados.
func (r *Room) Leave(userID string) error {
	if r.indexOf(userID) < 0 {
		return ErrNotParticipant
	}
	if r.Phase == PhaseDrafting {
		return ErrWrongPhase
	}
	r.remove(userID)
	return nil
}

func (r *Room) SetFormat(userID string, f Format) error {
	if userID != r.AdminID {
		return ErrNotAdmin
	}
	if r.Phase != PhaseLobby {
		return ErrWrongPhase
	}
	if !validFormat(f) {
		return ErrInvalidFormat
	}
	r.Format = f
	return nil
}

// SetSpinner dá ou tira de um participante a permissão de girar.
func (r *Room) SetSpinner(userID, target string, on bool) error {
	if userID != r.AdminID {
		return ErrNotAdmin
	}
	if r.indexOf(target) < 0 {
		return ErrNotParticipant
	}
	if on {
		r.Spinners[target] = true
	} else {
		delete(r.Spinners, target)
	}
	return nil
}

func (r *Room) CanSpin(userID string) bool {
	return userID != "" && (userID == r.AdminID || r.Spinners[userID])
}

// Spin sorteia o próximo da roda. O resultado só é aplicado em CommitSpin,
// depois que a animação termina.
func (r *Room) Spin(userID string) (*SpinState, error) {
	if !r.CanSpin(userID) {
		return nil, ErrCannotSpin
	}
	if r.CurrentSpin != nil {
		return nil, ErrSpinInProgress
	}
	switch r.Phase {
	case PhaseLobby:
		if len(r.Participants) < r.Format.Slots() {
			return nil, ErrNotEnoughPlayers
		}
		r.Teams = make([][]Participant, r.Format.Teams)
		for i := range r.Teams {
			r.Teams[i] = []Participant{}
		}
		r.Picks = 0
		r.Phase = PhaseDrafting
	case PhaseDrafting:
	default:
		return nil, ErrWrongPhase
	}

	r.CurrentSpin = &SpinState{
		WinnerID:     r.Pool[r.rng.IntN(len(r.Pool))],
		PoolSnapshot: slices.Clone(r.Pool),
		TargetOffset: 0.15 + 0.7*r.rng.Float64(),
		Rotations:    5 + r.rng.IntN(3),
		StartedAt:    time.Now(),
		DurationMs:   int(SpinDuration.Milliseconds()),
	}
	return r.CurrentSpin, nil
}

// CommitSpin tira o sorteado da roda e coloca no próximo time, alternando.
func (r *Room) CommitSpin() error {
	if r.CurrentSpin == nil {
		return ErrNoSpinInProgress
	}
	winner := r.CurrentSpin.WinnerID
	r.CurrentSpin = nil
	r.Pool = slices.DeleteFunc(r.Pool, func(id string) bool { return id == winner })

	team := r.Picks % r.Format.Teams
	p := r.Participants[r.indexOf(winner)]
	r.Teams[team] = append(r.Teams[team], p)
	r.Picks++

	if r.Picks == r.Format.Slots() {
		r.Phase = PhaseFinished
	}
	return nil
}

// Reset limpa tudo, inclusive os participantes. O próximo a entrar vira admin.
func (r *Room) Reset(userID string) error {
	if userID != r.AdminID {
		return ErrNotAdmin
	}
	if r.Phase != PhaseFinished {
		return ErrWrongPhase
	}
	*r = *NewRoomWithRand(r.rng)
	return nil
}

// Connect marca como online um participante que reconectou.
func (r *Room) Connect(userID string) {
	i := r.indexOf(userID)
	if i < 0 {
		return
	}
	r.Participants[i].Online = true
	if r.AdminID == "" {
		r.AdminID = userID
	}
}

// Disconnect trata a saída da atividade. Antes do primeiro giro a pessoa sai da roda;
// depois disso fica offline, mas continua no time.
func (r *Room) Disconnect(userID string) {
	i := r.indexOf(userID)
	if i < 0 {
		return
	}
	if r.Phase == PhaseLobby {
		r.remove(userID)
		return
	}
	r.Participants[i].Online = false
	if r.AdminID == userID {
		r.transferAdmin()
	}
}

// HasOnline indica se ainda há algum participante conectado.
func (r *Room) HasOnline() bool {
	return slices.ContainsFunc(r.Participants, func(p Participant) bool { return p.Online })
}

func (r *Room) remove(userID string) {
	r.Participants = slices.DeleteFunc(r.Participants, func(p Participant) bool { return p.ID == userID })
	r.Pool = slices.DeleteFunc(r.Pool, func(id string) bool { return id == userID })
	delete(r.Spinners, userID)
	if r.AdminID == userID {
		r.transferAdmin()
	}
}

// transferAdmin passa o admin para o participante online mais antigo.
func (r *Room) transferAdmin() {
	r.AdminID = ""
	for _, p := range r.Participants {
		if p.Online {
			r.AdminID = p.ID
			delete(r.Spinners, p.ID)
			return
		}
	}
}

func (r *Room) indexOf(userID string) int {
	return slices.IndexFunc(r.Participants, func(p Participant) bool { return p.ID == userID })
}
