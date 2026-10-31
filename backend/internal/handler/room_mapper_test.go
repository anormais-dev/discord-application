package handler

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/anormais-dev/discord-application/backend/internal/model"
)

func TestReconnectUntilOnlyForOfflineParticipants(t *testing.T) {
	until := time.Date(2026, 1, 1, 12, 0, 30, 0, time.UTC)
	out := toParticipants([]model.Participant{
		{ID: "online", Online: true},
		{ID: "caiu", Online: false, ReconnectUntil: until},
	})
	b, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	if out[0].ReconnectUntil != nil || strings.Count(string(b), "reconnectUntil") != 1 {
		t.Fatalf("só quem caiu deveria ter reconnectUntil, veio %s", b)
	}
	if !out[1].ReconnectUntil.Equal(until) {
		t.Fatalf("prazo = %v, esperado %v", out[1].ReconnectUntil, until)
	}
}
