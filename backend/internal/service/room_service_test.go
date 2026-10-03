package service

import (
	"errors"
	"testing"
	"time"

	"github.com/anormais-dev/discord-application/backend/internal/model"
	"github.com/anormais-dev/discord-application/backend/internal/repository"
)

type fakeClient struct {
	id     string
	states []model.Room
}

func (c *fakeClient) UserID() string { return c.id }

func (c *fakeClient) SendState(room *model.Room) { c.states = append(c.states, *room) }

func (c *fakeClient) last() model.Room { return c.states[len(c.states)-1] }

func TestRoomServiceBroadcastsToInstance(t *testing.T) {
	s := NewRoomService(repository.NewRoomRepository())
	a := &fakeClient{id: "a"}
	b := &fakeClient{id: "b"}
	other := &fakeClient{id: "c"}
	s.Attach("abc", a)
	s.Attach("abc", b)
	s.Attach("xyz", other)

	if err := s.Join("abc", model.Participant{ID: "a"}); err != nil {
		t.Fatal(err)
	}
	if b.last().AdminID != "a" {
		t.Fatalf("b deveria ver a como admin, veio %+v", b.last())
	}
	if len(other.states) != 1 {
		t.Fatalf("outra instância não deveria receber o estado, recebeu %d", len(other.states))
	}
	if err := s.Join("abc", model.Participant{ID: "a"}); !errors.Is(err, ErrAlreadyJoined) {
		t.Fatalf("erro = %v, esperado ErrAlreadyJoined", err)
	}
}

func TestRoomServiceCommitsSpinAfterDelay(t *testing.T) {
	s := NewRoomService(repository.NewRoomRepository())
	s.CommitDelay = 10 * time.Millisecond
	a := &fakeClient{id: "a"}
	s.Attach("abc", a)
	s.Join("abc", model.Participant{ID: "a"})
	s.Join("abc", model.Participant{ID: "b"})
	s.Join("abc", model.Participant{ID: "c"})
	s.SetReady("abc", "a", true)
	s.SetReady("abc", "b", true)
	s.SetReady("abc", "c", true)
	s.SetFormat("abc", "a", model.Format{Teams: 2, Size: 1})

	if err := s.Spin("abc", "a"); err != nil {
		t.Fatal(err)
	}
	if a.last().CurrentSpin == nil {
		t.Fatalf("estado deveria ter o giro em andamento")
	}

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		s.mu.Lock()
		picks := a.last().Picks
		s.mu.Unlock()
		if picks == 1 {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("o giro deveria ter sido aplicado depois do CommitDelay")
}

func TestRoomServiceDetachInLobbyRemovesParticipant(t *testing.T) {
	s := NewRoomService(repository.NewRoomRepository())
	a := &fakeClient{id: "a"}
	b := &fakeClient{id: "b"}
	s.Attach("abc", a)
	s.Attach("abc", b)
	s.Join("abc", model.Participant{ID: "a"})
	s.Join("abc", model.Participant{ID: "b"})

	s.Detach("abc", a)

	if got := b.last(); len(got.Participants) != 1 || got.AdminID != "b" {
		t.Fatalf("a deveria ter saído e b virado admin, veio %+v", got)
	}
}

func TestRoomServiceAutoSpinRunsUntilFinished(t *testing.T) {
	s := NewRoomService(repository.NewRoomRepository())
	s.CommitDelay = 5 * time.Millisecond
	s.AutoSpinDelay = 5 * time.Millisecond
	a := &fakeClient{id: "a"}
	s.Attach("abc", a)
	for _, id := range []string{"a", "b", "c", "d", "e"} {
		s.Join("abc", model.Participant{ID: id})
		s.SetReady("abc", id, true)
	}
	s.SetFormat("abc", "a", model.Format{Teams: 2, Size: 2})
	if err := s.SetAutoSpin("abc", "a", true); err != nil {
		t.Fatal(err)
	}
	if err := s.Spin("abc", "a"); err != nil {
		t.Fatal(err)
	}

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		s.mu.Lock()
		phase := a.last().Phase
		s.mu.Unlock()
		if phase == model.PhaseFinished {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("o giro automático deveria ter fechado os times")
}
