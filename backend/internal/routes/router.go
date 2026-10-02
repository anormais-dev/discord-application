package routes

import (
	"net/http"

	"github.com/anormais-dev/discord-application/backend/internal/handler"
)

type Handlers struct {
	Config *handler.ConfigHandler
	Token  *handler.TokenHandler
	WS     *handler.WSHandler
}

func NewRouter(h Handlers, staticDir string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/config", h.Config.Get)
	mux.HandleFunc("POST /api/token", h.Token.ExchangeCode)
	mux.HandleFunc("GET /api/ws", h.WS.ServeWS)
	mux.Handle("/", http.FileServer(http.Dir(staticDir)))

	return mux
}
