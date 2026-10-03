package service

import (
	crand "crypto/rand"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/anormais-dev/discord-application/backend/internal/model"
	"github.com/anormais-dev/discord-application/backend/internal/repository"
)

const (
	// Folga depois da animação antes de aplicar o resultado.
	defaultCommitDelay = SpinDuration + 500*time.Millisecond
	// Pausa entre o resultado de um giro e o próximo, com o giro automático ligado.
	defaultAutoSpinDelay = 2 * time.Second
	// Tempo que uma room sem ninguém conectado fica em memória.
	emptyRoomTTL = 5 * time.Minute
)

type Client interface {
	UserID() string
	// SendState é chamado com o lock do service travado: serialize na hora, sem guardar o ponteiro.
	SendState(room *model.Room)
}

type RoomService struct {
	repo          *repository.RoomRepository
	CommitDelay   time.Duration
	AutoSpinDelay time.Duration

	mu      sync.Mutex
	rng     *rand.Rand
	clients map[string]map[Client]struct{}
}

func NewRoomService(repo *repository.RoomRepository) *RoomService {
	var seed [32]byte
	crand.Read(seed[:])
	return NewRoomServiceWithRand(repo, rand.New(rand.NewChaCha8(seed)))
}

func NewRoomServiceWithRand(repo *repository.RoomRepository, rng *rand.Rand) *RoomService {
	return &RoomService{
		repo:          repo,
		CommitDelay:   defaultCommitDelay,
		AutoSpinDelay: defaultAutoSpinDelay,
		rng:           rng,
		clients:       map[string]map[Client]struct{}{},
	}
}

func (s *RoomService) Attach(instance string, c Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	room := s.repo.GetOrCreate(instance)
	if s.clients[instance] == nil {
		s.clients[instance] = map[Client]struct{}{}
	}
	s.clients[instance][c] = struct{}{}
	connect(room, c.UserID())
	s.broadcast(instance, room)
}

func (s *RoomService) Detach(instance string, c Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, err := s.repo.Get(instance)
	if err != nil {
		return
	}
	delete(s.clients[instance], c)
	if !s.hasUser(instance, c.UserID()) {
		disconnect(room, c.UserID())
	}
	s.broadcast(instance, room)

	if len(s.clients[instance]) == 0 {
		time.AfterFunc(emptyRoomTTL, func() { s.cleanup(instance, room) })
	}
}

func (s *RoomService) Join(instance string, p model.Participant) error {
	return s.update(instance, func(r *model.Room) error { return join(r, p) })
}

func (s *RoomService) Leave(instance, userID string) error {
	return s.update(instance, func(r *model.Room) error { return leave(r, userID) })
}

func (s *RoomService) SetFormat(instance, userID string, f model.Format) error {
	return s.update(instance, func(r *model.Room) error { return setFormat(r, userID, f) })
}

func (s *RoomService) SetSpinner(instance, userID, target string, on bool) error {
	return s.update(instance, func(r *model.Room) error { return setSpinner(r, userID, target, on) })
}

func (s *RoomService) SetReady(instance, userID string, on bool) error {
	return s.update(instance, func(r *model.Room) error { return setReady(r, userID, on) })
}

func (s *RoomService) Spin(instance, userID string) error {
	return s.update(instance, func(r *model.Room) error {
		if _, err := spin(r, s.rng, userID); err != nil {
			return err
		}
		s.scheduleCommit(instance, r)
		return nil
	})
}

// SetAutoSpin ligado no meio do sorteio já agenda o próximo giro.
func (s *RoomService) SetAutoSpin(instance, userID string, on bool) error {
	return s.update(instance, func(r *model.Room) error {
		if err := setAutoSpin(r, userID, on); err != nil {
			return err
		}
		s.scheduleAutoSpin(instance, r)
		return nil
	})
}

func (s *RoomService) Reset(instance, userID string) error {
	return s.update(instance, func(r *model.Room) error { return reset(r, userID) })
}

func (s *RoomService) update(instance string, fn func(*model.Room) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, err := s.repo.Get(instance)
	if err != nil {
		return err
	}
	if err := fn(room); err != nil {
		return err
	}
	s.broadcast(instance, room)
	return nil
}

func (s *RoomService) commit(instance string, room *model.Room) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if current, err := s.repo.Get(instance); err != nil || current != room {
		return
	}
	if err := commitSpin(room); err != nil {
		return
	}
	s.scheduleAutoSpin(instance, room)
	s.broadcast(instance, room)
}

func (s *RoomService) scheduleCommit(instance string, room *model.Room) {
	time.AfterFunc(s.CommitDelay, func() { s.commit(instance, room) })
}

func (s *RoomService) scheduleAutoSpin(instance string, room *model.Room) {
	if !room.AutoSpin || room.Phase != model.PhaseDrafting || room.CurrentSpin != nil {
		return
	}
	time.AfterFunc(s.AutoSpinDelay, func() { s.autoSpin(instance, room) })
}

func (s *RoomService) autoSpin(instance string, room *model.Room) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if current, err := s.repo.Get(instance); err != nil || current != room {
		return
	}
	if !autoSpin(room, s.rng) {
		return
	}
	s.scheduleCommit(instance, room)
	s.broadcast(instance, room)
}

func (s *RoomService) cleanup(instance string, room *model.Room) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, err := s.repo.Get(instance)
	if err != nil || current != room || len(s.clients[instance]) > 0 {
		return
	}
	s.repo.Delete(instance)
	delete(s.clients, instance)
}

func (s *RoomService) broadcast(instance string, room *model.Room) {
	for c := range s.clients[instance] {
		c.SendState(room)
	}
}

func (s *RoomService) hasUser(instance, userID string) bool {
	for c := range s.clients[instance] {
		if c.UserID() == userID {
			return true
		}
	}
	return false
}
