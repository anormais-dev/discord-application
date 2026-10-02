package main

import (
	"log"
	"net/http"

	"github.com/anormais-dev/discord-application/backend/internal/config"
	"github.com/anormais-dev/discord-application/backend/internal/handler"
	"github.com/anormais-dev/discord-application/backend/internal/repository"
	"github.com/anormais-dev/discord-application/backend/internal/service"
	"github.com/anormais-dev/discord-application/backend/pkg/discord"
)

func main() {
	cfg := config.Load()

	var dc service.DiscordClient
	if cfg.HasDiscordCredentials() {
		dc = discord.NewClient(cfg.DiscordClientID, cfg.DiscordClientSecret)
	} else if !cfg.DevAuth {
		log.Fatal("DISCORD_CLIENT_ID e DISCORD_CLIENT_SECRET são obrigatórios (ou use DEV_AUTH=true para testar fora do Discord)")
	}
	if cfg.DevAuth {
		log.Print("ATENÇÃO: DEV_AUTH ativo, qualquer um entra com qualquer nome. Não use em produção.")
	}
	auth := service.NewAuthService(dc, cfg.DevAuth)
	rooms := service.NewRoomService(repository.NewRoomRepository())

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/token", handler.NewTokenHandler(auth).ExchangeCode)
	mux.HandleFunc("GET /api/ws", handler.NewWSHandler(auth, rooms).ServeWS)
	mux.Handle("/", http.FileServer(http.Dir(cfg.StaticDir)))

	log.Printf("ouvindo em :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
