package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/anormais-dev/discord-application/backend/internal/dto"
	"github.com/anormais-dev/discord-application/backend/internal/response"
	"github.com/anormais-dev/discord-application/backend/internal/service"
)

type TokenHandler struct {
	auth *service.AuthService
}

func NewTokenHandler(auth *service.AuthService) *TokenHandler {
	return &TokenHandler{auth: auth}
}

// ExchangeCode troca o code do SDK pelo access token, que precisa do client secret.
func (h *TokenHandler) ExchangeCode(w http.ResponseWriter, r *http.Request) {
	var req dto.TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Code == "" {
		response.Error(w, http.StatusBadRequest, "code obrigatório")
		return
	}

	token, err := h.auth.ExchangeCode(r.Context(), req.Code)
	if err != nil {
		log.Printf("troca de token: %v", err)
		response.Error(w, http.StatusBadGateway, "falha ao trocar o code")
		return
	}

	response.JSON(w, http.StatusOK, dto.TokenResponse{AccessToken: token})
}
