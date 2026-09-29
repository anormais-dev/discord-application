// Package hub mantém uma room por instância da Activity e sincroniza os clientes via WebSocket.
package hub

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"

	"github.com/anormais-dev/discord-application/backend/internal/discord"
	"github.com/anormais-dev/discord-application/backend/internal/roulette"
)

const (
	authTimeout  = 10 * time.Second
	writeTimeout = 5 * time.Second
	// Folga depois da animação antes de aplicar o resultado.
	defaultCommitDelay = roulette.SpinDuration + 500*time.Millisecond
	// Tempo que uma room sem ninguém conectado fica em memória.
	emptyRoomTTL = 5 * time.Minute
)

// UserLookup resolve o usuário dono de um access token.
type UserLookup interface {
	CurrentUser(ctx context.Context, accessToken string) (*discord.User, error)
}

type Hub struct {
	users UserLookup
	// CommitDelay é quanto tempo depois do giro o resultado é aplicado.
	CommitDelay time.Duration

	mu    sync.Mutex
	rooms map[string]*roomHandle
}

type roomHandle struct {
	commitDelay time.Duration

	mu      sync.Mutex
	room    *roulette.Room
	clients map[*client]struct{}
	closed  bool
}

type client struct {
	user *discord.User
	send chan []byte
}

// inMessage cobre todas as mensagens que o cliente envia.
type inMessage struct {
	Type        string          `json:"type"`
	AccessToken string          `json:"accessToken"`
	Format      roulette.Format `json:"format"`
	UserID      string          `json:"userId"`
	Enabled     bool            `json:"enabled"`
}

type stateMessage struct {
	Type    string            `json:"type"`
	Room    json.RawMessage   `json:"room"`
	You     string            `json:"you"`
	Formats []roulette.Format `json:"formats"`
}

type errorMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func New(users UserLookup) *Hub {
	return &Hub{users: users, CommitDelay: defaultCommitDelay, rooms: map[string]*roomHandle{}}
}

// ServeWS trata GET /api/ws?instance=<instanceId>.
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	instance := r.URL.Query().Get("instance")
	if instance == "" {
		http.Error(w, "instance obrigatório", http.StatusBadRequest)
		return
	}
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*.discordsays.com"},
	})
	if err != nil {
		return
	}
	defer conn.CloseNow()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user, err := h.authenticate(ctx, conn)
	if err != nil {
		log.Printf("ws auth: %v", err)
		conn.Close(websocket.StatusPolicyViolation, "autenticação falhou")
		return
	}

	c := &client{user: user, send: make(chan []byte, 32)}
	rh := h.attach(instance, c)
	defer h.detach(instance, rh, c)

	go writeLoop(ctx, cancel, conn, c)

	for {
		var msg inMessage
		if err := wsjson.Read(ctx, conn, &msg); err != nil {
			return
		}
		rh.handle(c, msg)
	}
}

func (h *Hub) authenticate(ctx context.Context, conn *websocket.Conn) (*discord.User, error) {
	ctx, cancel := context.WithTimeout(ctx, authTimeout)
	defer cancel()
	var msg inMessage
	if err := wsjson.Read(ctx, conn, &msg); err != nil {
		return nil, err
	}
	if msg.Type != "auth" || msg.AccessToken == "" {
		return nil, errors.New("primeira mensagem precisa ser auth")
	}
	return h.users.CurrentUser(ctx, msg.AccessToken)
}

// attach registra o cliente na room da instância, criando a room se preciso.
func (h *Hub) attach(instance string, c *client) *roomHandle {
	for {
		h.mu.Lock()
		rh, ok := h.rooms[instance]
		if !ok {
			rh = &roomHandle{commitDelay: h.CommitDelay, room: roulette.NewRoom(), clients: map[*client]struct{}{}}
			h.rooms[instance] = rh
		}
		h.mu.Unlock()

		rh.mu.Lock()
		if rh.closed {
			rh.mu.Unlock()
			continue
		}
		rh.clients[c] = struct{}{}
		rh.room.Connect(c.user.ID)
		rh.broadcastState()
		rh.mu.Unlock()
		return rh
	}
}

func (h *Hub) detach(instance string, rh *roomHandle, c *client) {
	rh.mu.Lock()
	delete(rh.clients, c)
	if !rh.hasUser(c.user.ID) {
		rh.room.Disconnect(c.user.ID)
	}
	rh.broadcastState()
	empty := len(rh.clients) == 0
	rh.mu.Unlock()

	if empty {
		time.AfterFunc(emptyRoomTTL, func() { h.cleanup(instance, rh) })
	}
}

func (h *Hub) cleanup(instance string, rh *roomHandle) {
	h.mu.Lock()
	defer h.mu.Unlock()
	rh.mu.Lock()
	defer rh.mu.Unlock()
	if len(rh.clients) == 0 && h.rooms[instance] == rh {
		rh.closed = true
		delete(h.rooms, instance)
	}
}

func (rh *roomHandle) handle(c *client, msg inMessage) {
	rh.mu.Lock()
	defer rh.mu.Unlock()

	uid := c.user.ID
	var err error
	switch msg.Type {
	case "join":
		err = rh.room.Join(roulette.Participant{
			ID:     uid,
			Name:   c.user.DisplayName(),
			Avatar: c.user.AvatarURL(),
		})
	case "leave":
		err = rh.room.Leave(uid)
	case "set_format":
		err = rh.room.SetFormat(uid, msg.Format)
	case "set_spinner":
		err = rh.room.SetSpinner(uid, msg.UserID, msg.Enabled)
	case "spin":
		if _, err = rh.room.Spin(uid); err == nil {
			time.AfterFunc(rh.commitDelay, rh.commitSpin)
		}
	case "reset":
		err = rh.room.Reset(uid)
	default:
		err = errors.New("mensagem desconhecida")
	}

	if err != nil {
		c.enqueue(mustMarshal(errorMessage{Type: "error", Message: err.Error()}))
		return
	}
	rh.broadcastState()
}

func (rh *roomHandle) commitSpin() {
	rh.mu.Lock()
	defer rh.mu.Unlock()
	if err := rh.room.CommitSpin(); err != nil {
		return
	}
	rh.broadcastState()
}

// broadcastState envia o estado para todos. Precisa ser chamado com rh.mu travado.
func (rh *roomHandle) broadcastState() {
	room := mustMarshal(rh.room)
	for c := range rh.clients {
		c.enqueue(mustMarshal(stateMessage{
			Type:    "state",
			Room:    room,
			You:     c.user.ID,
			Formats: roulette.Formats,
		}))
	}
}

func (rh *roomHandle) hasUser(userID string) bool {
	for c := range rh.clients {
		if c.user.ID == userID {
			return true
		}
	}
	return false
}

// enqueue nunca bloqueia; um cliente lento demais perde mensagens e recebe o próximo estado completo.
func (c *client) enqueue(b []byte) {
	select {
	case c.send <- b:
	default:
	}
}

func writeLoop(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn, c *client) {
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		case b := <-c.send:
			wctx, wcancel := context.WithTimeout(ctx, writeTimeout)
			err := conn.Write(wctx, websocket.MessageText, b)
			wcancel()
			if err != nil {
				return
			}
		}
	}
}

func mustMarshal(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
