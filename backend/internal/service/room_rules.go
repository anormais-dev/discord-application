package service

import (
	"math/rand/v2"
	"slices"
	"time"

	"github.com/anormais-dev/discord-application/backend/internal/model"
)

// SpinDuration é quanto tempo a animação da roda dura nos clientes.
const SpinDuration = 5 * time.Second

func join(r *model.Room, p model.Participant) error {
	if r.Phase != model.PhaseLobby {
		return ErrWrongPhase
	}
	if indexOf(r, p.ID) >= 0 {
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

func leave(r *model.Room, userID string) error {
	if indexOf(r, userID) < 0 {
		return ErrNotParticipant
	}
	if r.Phase == model.PhaseDrafting {
		return ErrWrongPhase
	}
	remove(r, userID)
	return nil
}

func setFormat(r *model.Room, userID string, f model.Format) error {
	if userID != r.AdminID {
		return ErrNotAdmin
	}
	if r.Phase != model.PhaseLobby {
		return ErrWrongPhase
	}
	if !slices.Contains(model.Formats, f) {
		return ErrInvalidFormat
	}
	r.Format = f
	return nil
}

func setSpinner(r *model.Room, userID, target string, on bool) error {
	if userID != r.AdminID {
		return ErrNotAdmin
	}
	if indexOf(r, target) < 0 {
		return ErrNotParticipant
	}
	if on {
		r.Spinners[target] = true
	} else {
		delete(r.Spinners, target)
	}
	return nil
}

func setReady(r *model.Room, userID string, on bool) error {
	if r.Phase != model.PhaseLobby {
		return ErrWrongPhase
	}
	if indexOf(r, userID) < 0 {
		return ErrNotParticipant
	}
	if on {
		r.Ready[userID] = true
	} else {
		delete(r.Ready, userID)
	}
	return nil
}

func setAutoSpin(r *model.Room, userID string, on bool) error {
	if !canSpin(r, userID) {
		return ErrCannotSpin
	}
	if r.Phase == model.PhaseFinished {
		return ErrWrongPhase
	}
	r.AutoSpin = on
	return nil
}

func canSpin(r *model.Room, userID string) bool {
	return userID != "" && (userID == r.AdminID || r.Spinners[userID])
}

// spin só sorteia; o resultado entra nos times em commitSpin, depois da animação.
func spin(r *model.Room, rng *rand.Rand, userID string) (*model.Spin, error) {
	if !canSpin(r, userID) {
		return nil, ErrCannotSpin
	}
	if r.CurrentSpin != nil {
		return nil, ErrSpinInProgress
	}
	switch r.Phase {
	case model.PhaseLobby:
		if len(r.Participants) < r.Format.Slots() {
			return nil, ErrNotEnoughPlayers
		}
		if !allReady(r) {
			return nil, ErrNotAllReady
		}
		r.Teams = make([][]model.Participant, r.Format.Teams)
		for i := range r.Teams {
			r.Teams[i] = []model.Participant{}
		}
		r.Picks = 0
		r.Phase = model.PhaseDrafting
	case model.PhaseDrafting:
	default:
		return nil, ErrWrongPhase
	}

	return draw(r, rng), nil
}

// autoSpin gira sozinho entre um sorteio e outro, quando o giro automático está ligado.
func autoSpin(r *model.Room, rng *rand.Rand) bool {
	if !r.AutoSpin || r.Phase != model.PhaseDrafting || r.CurrentSpin != nil {
		return false
	}
	draw(r, rng)
	return true
}

func draw(r *model.Room, rng *rand.Rand) *model.Spin {
	r.CurrentSpin = newSpin(rng, r.Pool)
	return r.CurrentSpin
}

func newSpin(rng *rand.Rand, pool []string) *model.Spin {
	return &model.Spin{
		WinnerID:     pool[rng.IntN(len(pool))],
		PoolSnapshot: slices.Clone(pool),
		TargetOffset: 0.15 + 0.7*rng.Float64(),
		Rotations:    5 + rng.IntN(3),
		StartedAt:    time.Now(),
		DurationMs:   int(SpinDuration.Milliseconds()),
	}
}

// spinMap sorteia o mapa numa roda à parte, em qualquer fase, sem mexer nos times.
func spinMap(r *model.Room, rng *rand.Rand, userID string, maps []string) (*model.Spin, error) {
	if userID != r.AdminID {
		return nil, ErrNotAdmin
	}
	if r.MapSpin != nil {
		return nil, ErrSpinInProgress
	}
	if len(maps) == 0 {
		return nil, ErrMapsUnavailable
	}
	r.MapSpin = newSpin(rng, maps)
	return r.MapSpin, nil
}

// setMapOpen abre a roda de mapas para todo mundo; não fecha no meio de um giro.
func setMapOpen(r *model.Room, userID string, on bool) error {
	if userID != r.AdminID {
		return ErrNotAdmin
	}
	if !on && r.MapSpin != nil {
		return ErrSpinInProgress
	}
	r.MapOpen = on
	return nil
}

func commitMapSpin(r *model.Room) error {
	if r.MapSpin == nil {
		return ErrNoSpinInProgress
	}
	r.Map = r.MapSpin.WinnerID
	r.MapSpin = nil
	return nil
}

// commitSpin coloca o sorteado no próximo time. Se sobrar só uma pessoa na
// roda, ela entra direto, sem precisar de outro giro.
func commitSpin(r *model.Room) error {
	if r.CurrentSpin == nil {
		return ErrNoSpinInProgress
	}
	winner := r.CurrentSpin.WinnerID
	r.CurrentSpin = nil
	assign(r, winner)
	if r.Phase == model.PhaseDrafting && len(r.Pool) == 1 {
		assign(r, r.Pool[0])
	}
	return nil
}

// assign tira da roda e coloca no próximo time, alternando.
func assign(r *model.Room, userID string) {
	r.Pool = slices.DeleteFunc(r.Pool, func(id string) bool { return id == userID })

	team := r.Picks % r.Format.Teams
	p := r.Participants[indexOf(r, userID)]
	r.Teams[team] = append(r.Teams[team], p)
	r.Picks++

	if r.Picks == r.Format.Slots() {
		r.Phase = model.PhaseFinished
	}
}

// reset limpa inclusive os participantes; o próximo a entrar vira admin.
func reset(r *model.Room, userID string) error {
	if userID != r.AdminID {
		return ErrNotAdmin
	}
	if r.Phase != model.PhaseFinished {
		return ErrWrongPhase
	}
	*r = *model.NewRoom()
	return nil
}

func connect(r *model.Room, userID string) {
	i := indexOf(r, userID)
	if i < 0 {
		return
	}
	r.Participants[i].Online = true
	if r.AdminID == "" {
		r.AdminID = userID
	}
}

// disconnect tira da roda quem sai antes do primeiro giro; depois disso a pessoa
// fica offline, mas continua no time.
func disconnect(r *model.Room, userID string) {
	i := indexOf(r, userID)
	if i < 0 {
		return
	}
	if r.Phase == model.PhaseLobby {
		remove(r, userID)
		return
	}
	r.Participants[i].Online = false
	if r.AdminID == userID {
		transferAdmin(r)
	}
}

func remove(r *model.Room, userID string) {
	r.Participants = slices.DeleteFunc(r.Participants, func(p model.Participant) bool { return p.ID == userID })
	r.Pool = slices.DeleteFunc(r.Pool, func(id string) bool { return id == userID })
	delete(r.Spinners, userID)
	delete(r.Ready, userID)
	if r.AdminID == userID {
		transferAdmin(r)
	}
}

func allReady(r *model.Room) bool {
	for _, p := range r.Participants {
		if !r.Ready[p.ID] {
			return false
		}
	}
	return true
}

// transferAdmin passa o admin para o participante online mais antigo.
func transferAdmin(r *model.Room) {
	r.AdminID = ""
	for _, p := range r.Participants {
		if p.Online {
			r.AdminID = p.ID
			delete(r.Spinners, p.ID)
			return
		}
	}
}

func indexOf(r *model.Room, userID string) int {
	return slices.IndexFunc(r.Participants, func(p model.Participant) bool { return p.ID == userID })
}
