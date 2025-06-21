package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	services "github.com/okzmo/kyob/internal/service"
	"github.com/okzmo/kyob/internal/utils"
)

func CreateRole(w http.ResponseWriter, r *http.Request) {
	var body services.BodyRoleCreation
	serverID := chi.URLParam(r, "id")

	err := utils.ParseAndValidate(r, validate, &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	role, err := services.CreateRole(r.Context(), serverID, &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, role)
}

func GetRoles(w http.ResponseWriter, r *http.Request) {
	serverID := chi.URLParam(r, "id")

	roles, err := services.GetRoles(r.Context(), serverID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, roles)
}

func MoveRole(w http.ResponseWriter, r *http.Request) {
	var body services.BodyMoveRole

	err := utils.ParseAndValidate(r, validate, &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = services.MoveRole(r.Context(), &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, DefaultResponse{Message: "ok"})
}

func DeleteRole(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "role_id")

	err := services.DeleteRole(r.Context(), roleID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, DefaultResponse{Message: "ok"})
}
