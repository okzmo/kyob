package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	queries "github.com/okzmo/kyob/db/gen_queries"
	"github.com/okzmo/kyob/internal/api/actors"
	services "github.com/okzmo/kyob/internal/service"
	"github.com/okzmo/kyob/internal/utils"
	proto "github.com/okzmo/kyob/types"
)

func CreateChannel(w http.ResponseWriter, r *http.Request) {
	serverID := chi.URLParam(r, "server_id")
	var body services.CreateChannelBody

	err := utils.ParseAndValidate(r, validate, &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user := r.Context().Value("user").(queries.User)
	channelMessage := &proto.BodyChannelCreation{
		CreatorId:   user.ID,
		ServerId:    serverID,
		Name:        body.Name,
		Type:        string(body.Type),
		Description: body.Description,
		X:           body.X,
		Y:           body.Y,
	}

	serverPID := actors.ServersEngine.Registry.GetPID("server", serverID)
	actors.ServersEngine.Send(serverPID, channelMessage)

	utils.RespondWithJSON(w, http.StatusCreated, DefaultResponse{Message: "success"})
}

func EditChannel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body services.EditChannelBody

	err := utils.ParseAndValidate(r, validate, &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = services.EditChannel(r.Context(), id, &body)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorizedChannelEdition):
			utils.RespondWithError(w, http.StatusUnauthorized, "You cannot edit this channel.")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, &DefaultResponse{Message: "success"})
}

func DeleteChannel(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(queries.User)
	channelID := chi.URLParam(r, "channel_id")
	serverID := chi.URLParam(r, "server_id")

	protoMessage := &proto.BodyChannelRemoved{
		ServerId:  serverID,
		ChannelId: channelID,
		UserId:    user.ID,
	}
	serverPID := actors.ServersEngine.Registry.GetPID("server", serverID)
	actors.ServersEngine.Send(serverPID, protoMessage)

	utils.RespondWithJSON(w, http.StatusOK, &DefaultResponse{Message: "success"})
}

func ConnectToCall(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(queries.User)
	channelID := chi.URLParam(r, "channel_id")
	serverID := chi.URLParam(r, "server_id")

	token, err := services.GenerateCallToken(channelID, user.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	protoMessage := &proto.ConnectToCall{
		UserId:    user.ID,
		ServerId:  serverID,
		ChannelId: channelID,
	}

	channelPID := actors.ServersEngine.Registry.GetPID(fmt.Sprintf("server/%s/channel", serverID), channelID)
	actors.ServersEngine.Send(channelPID, protoMessage)

	utils.RespondWithJSON(w, http.StatusOK, &services.LivekitResponse{Token: token})
}

func DisconnectFromCall(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(queries.User)
	channelID := chi.URLParam(r, "channel_id")
	serverID := chi.URLParam(r, "server_id")

	protoMessage := &proto.DisconnectFromCall{
		UserId:    user.ID,
		ServerId:  serverID,
		ChannelId: channelID,
	}

	channelPID := actors.ServersEngine.Registry.GetPID(fmt.Sprintf("server/%s/channel", serverID), channelID)
	actors.ServersEngine.Send(channelPID, protoMessage)

	utils.RespondWithJSON(w, http.StatusOK, &DefaultResponse{Message: "success"})
}
