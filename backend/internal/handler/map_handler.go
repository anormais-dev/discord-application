package handler

import (
	"errors"
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

// List trata GET /api/maps. O carrossel de mapas usa essa lista e as imagens enquanto ninguém sorteia.
func (h *MapHandler) List(w http.ResponseWriter, r *http.Request) {
	maps, err := h.maps.CompetitiveMaps(r.Context())
	if err != nil {
		response.Error(w, http.StatusBadGateway, err.Error())
		return
	}
	// As imagens passam pelo backend: o Discord bloqueia media.valorant-api.com na atividade.
	images := map[string]dto.MapImagesResponse{}
	for name, img := range h.maps.Images() {
		images[name] = dto.MapImagesResponse{
			Card:   "/maps/" + img.UUID + "/card",
			Splash: "/maps/" + img.UUID + "/splash",
		}
	}
	response.JSON(w, http.StatusOK, dto.MapsResponse{Maps: maps, Images: images})
}

// Image trata GET /api/maps/{id}/{kind}: repassa a arte card ou splash do mapa.
func (h *MapHandler) Image(w http.ResponseWriter, r *http.Request) {
	img, err := h.maps.Image(r.Context(), r.PathValue("id"), r.PathValue("kind"))
	if errors.Is(err, service.ErrMapImageNotFound) {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		response.Error(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Content-Type", img.ContentType)
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(img.Data)
}
