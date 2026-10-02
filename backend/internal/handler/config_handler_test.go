package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anormais-dev/discord-application/backend/internal/dto"
)

func TestConfigHandlerReturnsClientID(t *testing.T) {
	rec := httptest.NewRecorder()
	NewConfigHandler("123").Get(rec, httptest.NewRequest(http.MethodGet, "/api/config", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, esperado 200", rec.Code)
	}
	var res dto.ConfigResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil || res.ClientID != "123" {
		t.Fatalf("resposta inesperada: %+v %v", res, err)
	}
}
