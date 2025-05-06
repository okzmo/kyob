package handlers

import (
	"net/http"

	services "github.com/okzmo/kyob/internal/service"
	"github.com/okzmo/kyob/internal/utils"
)

type DefaultResponse struct {
	Message string `json:"message"`
}

func Setup(w http.ResponseWriter, r *http.Request) {
	all, err := services.GetSetup(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusContinue, *all)
}
