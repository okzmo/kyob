package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-chi/chi/v5"
	"github.com/okzmo/kyob/db"
	"github.com/okzmo/kyob/internal/api/actors"
	services "github.com/okzmo/kyob/internal/service"
	"github.com/okzmo/kyob/internal/utils"
	proto "github.com/okzmo/kyob/types"
)

func CreateServer(w http.ResponseWriter, r *http.Request) {
	var body services.CreateServerBody

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		slog.Error(err.Error())
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to parse given image.")
		return
	}

	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		slog.Error(err.Error())
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to get image.")
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		slog.Error("Failed to read file", "err", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to process file.")
		return
	}

	body.Name = r.FormValue("name")
	descriptionJSON := r.FormValue("description")
	if descriptionJSON != "" {
		if err := json.Unmarshal([]byte(descriptionJSON), &body.Description); err != nil {
			slog.Error(err.Error())
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid description.")
			return
		}
	}

	body.Private = r.FormValue("private") == "true"
	x, errX := strconv.Atoi(r.FormValue("x"))
	y, errY := strconv.Atoi(r.FormValue("y"))
	if errX != nil || errY != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid coordinates.")
	}

	body.X = x
	body.Y = y

	cropJSON := r.FormValue("crop")
	var crop services.Crop
	if err := json.Unmarshal([]byte(cropJSON), &crop); err != nil {
		slog.Error(err.Error())
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid crop data.")
		return
	}
	body.Crop = crop

	err = validate.Struct(body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	server, err := services.CreateServer(r.Context(), fileData, fileHeader, &body)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrTooManyServers):
			utils.RespondWithError(w, http.StatusForbidden, "You cannot create more than 200 servers.", "ERR_TOO_MANY_SERVERS")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	user := r.Context().Value("user").(db.User)
	serverPID := actors.ServersEngine.Spawn(actors.NewServer, "server", actor.WithID(server.ID))
	userPID := actors.UsersEngine.Registry.GetPID("user", user.ID)
	actors.UsersEngine.Send(userPID, &proto.NewServerCreated{
		ActorId:      serverPID.ID,
		ActorAddress: serverPID.Address,
	})

	utils.RespondWithJSON(w, http.StatusCreated, server)
}

func EditServer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body services.EditServerBody

	err := utils.ParseAndValidate(r, validate, &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = services.EditServer(r.Context(), id, &body)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorizedServerEdition):
			utils.RespondWithError(w, http.StatusUnauthorized, "You cannot edit this server.")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, &DefaultResponse{Message: "success"})
}

func DeleteServer(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(db.User)
	id := chi.URLParam(r, "id")

	protoMessage := &proto.BodyServerRemoved{
		ServerId: id,
		UserId:   user.ID,
	}
	serverPID := actors.ServersEngine.Registry.GetPID("server", id)
	actors.ServersEngine.Send(serverPID, protoMessage)

	utils.RespondWithJSON(w, http.StatusOK, &DefaultResponse{Message: "success"})
}

func CreateServerInvite(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	inviteId, err := services.CreateServerInvite(r.Context(), id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, &services.ServerInviteResponse{InviteLink: fmt.Sprintf("http://localhost:5173/invite/%s", *inviteId)})
}

func JoinServer(w http.ResponseWriter, r *http.Request) {
	var body services.JoinServerBody

	err := utils.ParseAndValidate(r, validate, &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	server, err := services.JoinServer(r.Context(), body)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrNoIdInInvite):
			utils.RespondWithError(w, http.StatusBadRequest, "The invite url is invalid.", "ERR_INVITE_MISSING_ID")
		case errors.Is(err, services.ErrServerNotFound):
			utils.RespondWithError(w, http.StatusNotFound, "The given url doesn't match any existing realm.")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	user := r.Context().Value("user").(db.User)

	userPID := actors.UsersEngine.Registry.GetPID("user", user.ID)
	serverPID := actors.ServersEngine.Registry.GetPID("server", server.ID)
	actors.ServersEngine.SendWithSender(serverPID, &proto.Connect{Type: "JOIN_SERVER"}, userPID)
	actors.ServersEngine.Send(serverPID, &proto.BodyNewUserInServer{
		ServerId: server.ID,
		User: &proto.User{
			Id:          user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Avatar:      &user.Avatar.String,
		},
	})

	utils.RespondWithJSON(w, http.StatusOK, services.JoinServerResponse{Server: *server})
}

func LeaveServer(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(db.User)
	serverId := chi.URLParam(r, "id")

	err := services.LeaveServer(r.Context(), serverId, user.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userPID := actors.UsersEngine.Registry.GetPID("user", user.ID)
	serverPID := actors.ServersEngine.Registry.GetPID("server", serverId)
	actors.ServersEngine.SendWithSender(serverPID, &proto.Disconnect{
		Type: "LEAVE_SERVER",
	}, userPID)

	utils.RespondWithJSON(w, http.StatusOK, DefaultResponse{Message: "success"})
}
