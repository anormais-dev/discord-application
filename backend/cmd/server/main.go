package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/anormais-dev/discord-application/backend/internal/config"
	"github.com/anormais-dev/discord-application/backend/internal/devauth"
	"github.com/anormais-dev/discord-application/backend/internal/hub"
	"github.com/anormais-dev/discord-application/backend/pkg/discord"
)

func main() {
	cfg := config.Load()

	dc := discord.NewClient(cfg.DiscordClientID, cfg.DiscordClientSecret)

	var users hub.UserLookup = dc
	if cfg.DevAuth {
		log.Print("ATENÇÃO: DEV_AUTH ativo, qualquer um entra com qualquer nome. Não use em produção.")
		dev := devauth.Users{}
		if cfg.HasDiscordCredentials() {
			dev.Next = dc
		}
		users = dev
	} else if !cfg.HasDiscordCredentials() {
		log.Fatal("DISCORD_CLIENT_ID e DISCORD_CLIENT_SECRET são obrigatórios (ou use DEV_AUTH=true para testar fora do Discord)")
	}
	h := hub.New(users)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/token", tokenHandler(dc))
	mux.HandleFunc("GET /api/ws", h.ServeWS)
	mux.Handle("/", http.FileServer(http.Dir(cfg.StaticDir)))

	log.Printf("ouvindo em :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}

// tokenHandler troca o code do SDK pelo access token, que precisa do client secret.
func tokenHandler(dc *discord.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Code string `json:"code"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Code == "" {
			http.Error(w, "code obrigatório", http.StatusBadRequest)
			return
		}
		token, err := dc.ExchangeCode(r.Context(), body.Code)
		if err != nil {
			log.Printf("troca de token: %v", err)
			http.Error(w, "falha ao trocar o code", http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"access_token": token})
	}
}
