package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/anormais-dev/discord-application/backend/internal/config"
	"github.com/anormais-dev/discord-application/backend/internal/handler"
	"github.com/anormais-dev/discord-application/backend/internal/repository"
	"github.com/anormais-dev/discord-application/backend/internal/routes"
	"github.com/anormais-dev/discord-application/backend/internal/service"
	"github.com/anormais-dev/discord-application/backend/pkg/discord"
)

var ErrMissingCredentials = errors.New("DISCORD_CLIENT_ID e DISCORD_CLIENT_SECRET são obrigatórios (ou use DEV_AUTH=true para testar fora do Discord)")

func Build(cfg *config.Config) (http.Handler, error) {
	var dc service.DiscordClient
	if cfg.HasDiscordCredentials() {
		dc = discord.NewClient(cfg.DiscordClientID, cfg.DiscordClientSecret)
	} else if !cfg.DevAuth {
		return nil, ErrMissingCredentials
	}
	if cfg.DevAuth {
		log.Print("ATENÇÃO: DEV_AUTH ativo, qualquer um entra com qualquer nome. Não use em produção.")
	}

	authService := service.NewAuthService(dc, cfg.DevAuth)
	roomService := service.NewRoomService(repository.NewRoomRepository())

	handlers := routes.Handlers{
		Config: handler.NewConfigHandler(cfg.DiscordClientID),
		Token:  handler.NewTokenHandler(authService),
		WS:     handler.NewWSHandler(authService, roomService),
	}

	return routes.NewRouter(handlers, cfg.StaticDir), nil
}
