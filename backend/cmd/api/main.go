package main

import (
	"log"
	"net/http"

	"github.com/anormais-dev/discord-application/backend/internal/config"
	"github.com/anormais-dev/discord-application/backend/internal/server"
)

func main() {
	cfg := config.Load()

	app, err := server.Build(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("ouvindo em :%s", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, app); err != nil {
		log.Fatal(err)
	}
}
