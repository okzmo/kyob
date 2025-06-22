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
	queries "github.com/okzmo/kyob/db/gen_queries"
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

	user := r.Context().Value("user").(queries.User)
	serverPID := actors.ServersEngine.Spawn(actors.NewServer, "server", actor.WithID(server.ID))
	userPID := actors.UsersEngine.Registry.GetPID("user", user.ID)
	actors.UsersEngine.Send(userPID, &proto.NewServerCreated{
		ActorId:      serverPID.ID,
		ActorAddress: serverPID.Address,
	})

	utils.RespondWithJSON(w, http.StatusCreated, server)
}

func UpdateServerProfile(w http.ResponseWriter, r *http.Request) {
	var body services.UpdateServerProfileBody
	serverID := chi.URLParam(r, "id")

	err := utils.ParseAndValidate(r, validate, &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = services.UpdateServerProfile(r.Context(), serverID, &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	serverPID := actors.ServersEngine.Registry.GetPID("server", serverID)
	messageToSend := &proto.ServerChangedInformations{
		ServerId: serverID,
		ServerInformations: &proto.ServerInformations{
			Name:        &body.Name,
			Description: body.Description,
		},
	}

	actors.ServersEngine.Send(serverPID, messageToSend)

	utils.RespondWithJSON(w, http.StatusOK, DefaultResponse{Message: "success"})
}

func UpdateServerAvatar(w http.ResponseWriter, r *http.Request) {
	var body services.UpdateAvatarBody
	var cropAvatar, cropBanner services.Crop
	serverID := chi.URLParam(r, "id")

	config := utils.ImageValidationConfig{
		MaxSize: 10 << 20, // 10 MB
		AllowedMimeTypes: []string{
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/webp",
		},
		RequireValidHeader: true,
	}

	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		slog.Error(err.Error())
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to get image.")
		return
	}
	defer file.Close()

	if err := utils.ParseAndValidateImage(fileHeader, config); err != nil {
		slog.Error("Avatar image validation failed", "error", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Avatar's invalid.")
		return
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		slog.Error("Failed to read file", "err", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to process file.")
		return
	}

	cropAvatarJSON := r.FormValue("crop_avatar")
	cropBannerJSON := r.FormValue("crop_banner")
	mainColor := r.FormValue("main_color")

	if err := json.Unmarshal([]byte(cropAvatarJSON), &cropAvatar); err != nil {
		slog.Error(err.Error())
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid crop data.")
		return
	}

	if err := json.Unmarshal([]byte(cropBannerJSON), &cropBanner); err != nil {
		slog.Error(err.Error())
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid crop data.")
		return
	}

	body.CropAvatar = cropAvatar
	body.CropBanner = cropBanner
	body.MainColor = mainColor

	err = validate.Struct(body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := services.UpdateServerAvatar(r.Context(), serverID, fileData, fileHeader, &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	serverPID := actors.ServersEngine.Registry.GetPID("server", serverID)
	actors.ServersEngine.Send(serverPID, &proto.ServerChangedInformations{
		ServerId: serverID,
		ServerInformations: &proto.ServerInformations{
			Avatar:    &res.Avatar,
			Banner:    &res.Banner,
			MainColor: &body.MainColor,
		},
	})

	utils.RespondWithJSON(w, http.StatusOK, res)
}

func DeleteServer(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(queries.User)
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

	inviteID, err := services.CreateServerInvite(r.Context(), id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, &services.ServerInviteResponse{InviteLink: fmt.Sprintf("http://localhost:5173/invite/%s", *inviteID)})
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
		case errors.Is(err, services.ErrNoIDInInvite):
			utils.RespondWithError(w, http.StatusBadRequest, "The invite url is invalid.", "ERR_INVITE_MISSING_ID")
		case errors.Is(err, services.ErrServerNotFound):
			utils.RespondWithError(w, http.StatusNotFound, "The given url doesn't match any existing realm.")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	user := r.Context().Value("user").(queries.User)

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
	user := r.Context().Value("user").(queries.User)
	serverID := chi.URLParam(r, "id")

	err := services.LeaveServer(r.Context(), serverID, user.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userPID := actors.UsersEngine.Registry.GetPID("user", user.ID)
	serverPID := actors.ServersEngine.Registry.GetPID("server", serverID)
	actors.ServersEngine.SendWithSender(serverPID, &proto.Disconnect{
		Type: "LEAVE_SERVER",
	}, userPID)

	utils.RespondWithJSON(w, http.StatusOK, DefaultResponse{Message: "success"})
}
