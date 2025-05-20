package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	services "github.com/okzmo/kyob/internal/service"
	"github.com/okzmo/kyob/internal/utils"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "user_id")
	userId, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	user, err := services.GetUser(r.Context(), int64(userId))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	utils.RespondWithJSON(w, http.StatusContinue, user)
}
