package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anormais-dev/discord-application/backend/internal/dto"
	"github.com/anormais-dev/discord-application/backend/internal/service"
)

func exchange(body string) *httptest.ResponseRecorder {
	h := NewTokenHandler(service.NewAuthService(fakeDiscord{}, false))
	rec := httptest.NewRecorder()
	h.ExchangeCode(rec, httptest.NewRequest(http.MethodPost, "/api/token", strings.NewReader(body)))
	return rec
}

func TestTokenHandlerExchangesCode(t *testing.T) {
	rec := exchange(`{"code":"abc"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, esperado 200", rec.Code)
	}
	var res dto.TokenResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil || res.AccessToken != "token-abc" {
		t.Fatalf("resposta inesperada: %+v %v", res, err)
	}
}

func TestTokenHandlerErrors(t *testing.T) {
	cases := map[string]int{
		`{}`:                  http.StatusBadRequest,
		`não é json`:          http.StatusBadRequest,
		`{"code":"invalido"}`: http.StatusBadGateway,
	}
	for body, want := range cases {
		if rec := exchange(body); rec.Code != want {
			t.Errorf("body %s: status = %d, esperado %d", body, rec.Code, want)
		}
	}
}
