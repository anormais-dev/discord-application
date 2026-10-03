package handler

import (
	"net/http"

	"github.com/anormais-dev/discord-application/backend/internal/dto"
	"github.com/anormais-dev/discord-application/backend/internal/response"
	"github.com/anormais-dev/discord-application/backend/internal/service"
)

type MapHandler struct {
	maps *service.MapService
}

func NewMapHandler(maps *service.MapService) *MapHandler {
	return &MapHandler{maps: maps}
}

// List trata GET /api/maps. A roda de mapas usa essa lista enquanto ninguém gira.
func (h *MapHandler) List(w http.ResponseWriter, r *http.Request) {
	maps, err := h.maps.CompetitiveMaps(r.Context())
	if err != nil {
		response.Error(w, http.StatusBadGateway, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, dto.MapsResponse{Maps: maps})
}
