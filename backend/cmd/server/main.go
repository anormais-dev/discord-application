package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/anormais-dev/discord-application/backend/internal/devauth"
	"github.com/anormais-dev/discord-application/backend/internal/discord"
	"github.com/anormais-dev/discord-application/backend/internal/hub"
)

func main() {
	loadEnvFile("../.env")

	port := getenv("PORT", "3000")
	dc := discord.NewClient(os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_CLIENT_SECRET"))
	hasCredentials := dc.ClientID != "" && dc.ClientSecret != ""

	var users hub.UserLookup = dc
	if os.Getenv("DEV_AUTH") == "true" {
		log.Print("ATENÇÃO: DEV_AUTH ativo, qualquer um entra com qualquer nome. Não use em produção.")
		dev := devauth.Users{}
		if hasCredentials {
			dev.Next = dc
		}
		users = dev
	} else if !hasCredentials {
		log.Fatal("DISCORD_CLIENT_ID e DISCORD_CLIENT_SECRET são obrigatórios (ou use DEV_AUTH=true para testar fora do Discord)")
	}
	h := hub.New(users)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/token", tokenHandler(dc))
	mux.HandleFunc("GET /api/ws", h.ServeWS)
	mux.Handle("/", http.FileServer(http.Dir(getenv("STATIC_DIR", "../frontend/dist"))))

	log.Printf("ouvindo em :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
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

// loadEnvFile carrega KEY=VALUE de um arquivo, sem sobrescrever o que já está no ambiente.
func loadEnvFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, strings.Trim(strings.TrimSpace(value), `"'`))
		}
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
