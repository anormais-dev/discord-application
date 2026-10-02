package response

import (
	"encoding/json"
	"net/http"

	"github.com/anormais-dev/discord-application/backend/internal/dto"
)

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}

	_ = json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, dto.ErrorResponse{Error: message})
}
