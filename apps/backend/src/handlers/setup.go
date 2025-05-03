package handlers

import (
	"net/http"

	"github.com/okzmo/kyob/src/services"
	"github.com/okzmo/kyob/src/utils"
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
