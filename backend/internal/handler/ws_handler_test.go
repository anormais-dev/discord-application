package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"

	"github.com/anormais-dev/discord-application/backend/internal/dto"
	"github.com/anormais-dev/discord-application/backend/internal/repository"
	"github.com/anormais-dev/discord-application/backend/internal/service"
	"github.com/anormais-dev/discord-application/backend/pkg/discord"
)

// fakeDiscord usa o próprio token como ID do usuário.
type fakeDiscord struct{}

func (fakeDiscord) ExchangeCode(_ context.Context, code string) (string, error) {
	if code == "invalido" {
		return "", errors.New("code inválido")
	}
	return "token-" + code, nil
}

func (fakeDiscord) CurrentUser(_ context.Context, token string) (*discord.User, error) {
	if token == "invalido" {
		return nil, errors.New("token inválido")
	}
	return &discord.User{ID: token, Username: token}, nil
}

type received struct {
	Type    string           `json:"type"`
	Room    dto.RoomResponse `json:"room"`
	You     string           `json:"you"`
	Message string           `json:"message"`
}

type testClient struct {
	t    *testing.T
	conn *websocket.Conn
}

func dial(t *testing.T, srv *httptest.Server, token string) *testClient {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	url := "ws" + strings.TrimPrefix(srv.URL, "http") + "/api/ws?instance=abc"
	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { conn.CloseNow() })
	c := &testClient{t: t, conn: conn}
	c.send(map[string]any{"type": "auth", "accessToken": token})
	return c
}

func (c *testClient) send(v any) {
	c.t.Helper()
	if err := wsjson.Write(context.Background(), c.conn, v); err != nil {
		c.t.Fatalf("write: %v", err)
	}
}

// until lê mensagens até uma satisfazer a condição.
func (c *testClient) until(cond func(received) bool) received {
	c.t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for {
		var raw json.RawMessage
		if err := wsjson.Read(ctx, c.conn, &raw); err != nil {
			c.t.Fatalf("read: %v", err)
		}
		var m received
		if err := json.Unmarshal(raw, &m); err != nil {
			c.t.Fatalf("unmarshal: %v", err)
		}
		if cond(m) {
			return m
		}
	}
}

func TestWSSyncsTwoClients(t *testing.T) {
	rooms := service.NewRoomService(repository.NewRoomRepository())
	rooms.CommitDelay = 50 * time.Millisecond
	h := NewWSHandler(service.NewAuthService(fakeDiscord{}, false), rooms)
	srv := httptest.NewServer(http.HandlerFunc(h.ServeWS))
	defer srv.Close()

	a := dial(t, srv, "alice")
	a.until(func(m received) bool { return m.Type == "state" })
	b := dial(t, srv, "bob")
	b.until(func(m received) bool { return m.Type == "state" })

	a.send(map[string]any{"type": "join"})
	b.until(func(m received) bool { return m.Type == "state" && m.Room.AdminID == "alice" })
	b.send(map[string]any{"type": "join"})
	a.until(func(m received) bool { return m.Type == "state" && len(m.Room.Participants) == 2 })

	b.send(map[string]any{"type": "spin"})
	b.until(func(m received) bool { return m.Type == "error" })

	a.send(map[string]any{"type": "set_format", "format": map[string]int{"teams": 2, "size": 1}})
	a.send(map[string]any{"type": "spin"})
	b.until(func(m received) bool { return m.Type == "state" && m.Room.Spin != nil })
	b.until(func(m received) bool { return m.Type == "state" && m.Room.Picks == 1 })

	b.send(map[string]any{"type": "leave"})
	b.until(func(m received) bool { return m.Type == "error" })

	a.send(map[string]any{"type": "spin"})
	final := b.until(func(m received) bool { return m.Type == "state" && m.Room.Phase == "finished" })
	if len(final.Room.Teams[0]) != 1 || len(final.Room.Teams[1]) != 1 {
		t.Fatalf("times inesperados: %+v", final.Room.Teams)
	}
}

func TestWSRejectsInvalidToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(NewWSHandler(service.NewAuthService(fakeDiscord{}, false), service.NewRoomService(repository.NewRoomRepository())).ServeWS))
	defer srv.Close()
	c := dial(t, srv, "invalido")
	var m json.RawMessage
	err := wsjson.Read(context.Background(), c.conn, &m)
	if websocket.CloseStatus(err) != websocket.StatusPolicyViolation {
		t.Fatalf("esperado fechamento por policy violation, veio %v", err)
	}
}
