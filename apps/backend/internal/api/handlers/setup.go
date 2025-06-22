package handlers

import (
	"net/http"

	queries "github.com/okzmo/kyob/db/gen_queries"
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

	utils.RespondWithJSON(w, http.StatusOK, *all)
}

func SaveLastState(w http.ResponseWriter, r *http.Request) {
	var body services.BodySaveState
	user := r.Context().Value("user").(queries.User)

	err := utils.ParseAndValidate(r, validate, &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	body.UserID = user.ID

	err = services.SaveLastState(r.Context(), body)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, DefaultResponse{Message: "ok"})
}
