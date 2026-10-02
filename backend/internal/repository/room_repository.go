package repository

import (
	"errors"
	"sync"

	"github.com/anormais-dev/discord-application/backend/internal/model"
)

var ErrRoomNotFound = errors.New("room não encontrada")

type RoomRepository struct {
	mu    sync.RWMutex
	rooms map[string]*model.Room
}

func NewRoomRepository() *RoomRepository {
	return &RoomRepository{
		rooms: make(map[string]*model.Room),
	}
}

func (r *RoomRepository) Get(instance string) (*model.Room, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	room, ok := r.rooms[instance]
	if !ok {
		return nil, ErrRoomNotFound
	}

	return room, nil
}

func (r *RoomRepository) GetOrCreate(instance string) *model.Room {
	r.mu.Lock()
	defer r.mu.Unlock()

	room, ok := r.rooms[instance]
	if !ok {
		room = model.NewRoom()
		r.rooms[instance] = room
	}

	return room
}

func (r *RoomRepository) Delete(instance string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.rooms, instance)
}
