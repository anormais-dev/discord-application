package handler

import (
	"net/http"

	"github.com/anormais-dev/discord-application/backend/internal/dto"
	"github.com/anormais-dev/discord-application/backend/internal/response"
)

type ConfigHandler struct {
	clientID string
}

func NewConfigHandler(clientID string) *ConfigHandler {
	return &ConfigHandler{clientID: clientID}
}

// Get entrega ao frontend o client ID, que é público, para ele não depender do build.
func (h *ConfigHandler) Get(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, dto.ConfigResponse{ClientID: h.clientID})
}
