// Package hub mantém uma room por instância da Activity e sincroniza os clientes via WebSocket.
package hub

import (
	"net/http"
	"sync"

	"github.com/anormais-dev/discord-application/backend/internal/discord"
	"github.com/anormais-dev/discord-application/backend/internal/roulette"
)

type Hub struct {
	discord *discord.Client

	mu    sync.Mutex
	rooms map[string]*roomHandle
}

type roomHandle struct {
	room *roulette.Room
}

func New(dc *discord.Client) *Hub {
	return &Hub{discord: dc, rooms: map[string]*roomHandle{}}
}

// ServeWS trata GET /api/ws?instance=<instanceId>.
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "não implementado", http.StatusNotImplemented)
}
