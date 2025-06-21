package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/okzmo/kyob/internal/api/actors"
	services "github.com/okzmo/kyob/internal/service"
	"github.com/okzmo/kyob/internal/utils"
	proto "github.com/okzmo/kyob/types"
)

func CreateRole(w http.ResponseWriter, r *http.Request) {
	var body services.BodyRoleCreation
	serverID := chi.URLParam(r, "id")

	err := utils.ParseAndValidate(r, validate, &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	message := &proto.CreateRole{
		Name:      body.Name,
		Idx:       int32(body.Index),
		ServerId:  serverID,
		Color:     body.Color,
		Abilities: body.Abilities,
	}

	serverPID := actors.ServersEngine.Registry.GetPID("server", serverID)
	actors.ServersEngine.Send(serverPID, message)

	utils.RespondWithJSON(w, http.StatusOK, DefaultResponse{Message: "ok"})
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

func AddRoleMember(w http.ResponseWriter, r *http.Request) {
	var body services.BodyAddOrRemoveRole
	serverID := chi.URLParam(r, "id")

	err := utils.ParseAndValidate(r, validate, &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	message := &proto.AddRoleMember{
		UserId:   body.UserID,
		Id:       body.RoleID,
		ServerId: serverID,
	}

	serverPID := actors.ServersEngine.Registry.GetPID("server", serverID)
	actors.ServersEngine.Send(serverPID, message)

	utils.RespondWithJSON(w, http.StatusOK, DefaultResponse{Message: "ok"})
}

func RemoveRoleMember(w http.ResponseWriter, r *http.Request) {
	var body services.BodyAddOrRemoveRole
	serverID := chi.URLParam(r, "id")

	err := utils.ParseAndValidate(r, validate, &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	message := &proto.RemoveRoleMember{
		UserId:   body.UserID,
		Id:       body.RoleID,
		ServerId: serverID,
	}

	serverPID := actors.ServersEngine.Registry.GetPID("server", serverID)
	actors.ServersEngine.Send(serverPID, message)

	utils.RespondWithJSON(w, http.StatusOK, DefaultResponse{Message: "ok"})
}

func MoveRole(w http.ResponseWriter, r *http.Request) {
	var body services.BodyMoveRole
	serverID := chi.URLParam(r, "id")

	err := utils.ParseAndValidate(r, validate, &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	message := &proto.ChangeRoleRanking{
		Id:       body.RoleID,
		ServerId: serverID,
		From:     int32(body.From),
		To:       int32(body.To),
	}

	serverPID := actors.ServersEngine.Registry.GetPID("server", serverID)
	actors.ServersEngine.Send(serverPID, message)

	utils.RespondWithJSON(w, http.StatusOK, DefaultResponse{Message: "ok"})
}

func DeleteRole(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "role_id")
	serverID := chi.URLParam(r, "id")

	err := services.DeleteRole(r.Context(), roleID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	message := &proto.RemoveRoleMember{
		Id:       roleID,
		ServerId: serverID,
	}

	serverPID := actors.ServersEngine.Registry.GetPID("server", serverID)
	actors.ServersEngine.Send(serverPID, message)

	utils.RespondWithJSON(w, http.StatusOK, DefaultResponse{Message: "ok"})
}
