package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	services "github.com/okzmo/kyob/internal/service"
	"github.com/okzmo/kyob/internal/utils"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "user_id")

	user, err := services.GetUser(r.Context(), userId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	utils.RespondWithJSON(w, http.StatusContinue, user)
}
