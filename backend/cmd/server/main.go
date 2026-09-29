package main

import (
	"log"
	"net/http"
	"os"

	"github.com/anormais-dev/discord-application/backend/internal/discord"
	"github.com/anormais-dev/discord-application/backend/internal/hub"
)

func main() {
	port := getenv("PORT", "3000")
	dc := &discord.Client{
		ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
	}
	h := hub.New(dc)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/token", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "não implementado", http.StatusNotImplemented)
	})
	mux.HandleFunc("GET /api/ws", h.ServeWS)
	mux.Handle("/", http.FileServer(http.Dir(getenv("STATIC_DIR", "../frontend/dist"))))

	log.Printf("ouvindo em :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
