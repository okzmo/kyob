package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/okzmo/kyob/db"
	"github.com/okzmo/kyob/internal/api/actors"
	services "github.com/okzmo/kyob/internal/service"
	"github.com/okzmo/kyob/internal/utils"
	proto "github.com/okzmo/kyob/types"
)

func CreateChannel(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "server_id")
	serverId, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var body services.CreateChannelBody

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = validate.Struct(body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user := r.Context().Value("user").(db.User)
	channelMessage := &proto.BodyChannelCreation{
		CreatorId:   user.ID,
		ServerId:    int32(serverId),
		Name:        body.Name,
		Type:        string(body.Type),
		Description: body.Description,
		X:           body.X,
		Y:           body.Y,
	}

	channelPID := actors.ServersEngine.Registry.GetPID("server", strconv.Itoa(serverId))
	actors.ServersEngine.Send(channelPID, channelMessage)

	utils.RespondWithJSON(w, http.StatusCreated, DefaultResponse{Message: "success"})
}

func EditChannel(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var body services.EditChannelBody

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = validate.Struct(body)
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

	utils.RespondWithJSON(w, http.StatusContinue, &DefaultResponse{Message: "success"})
}

func DeleteChannel(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(db.User)
	channelIdParam := chi.URLParam(r, "channel_id")
	serverIdParam := chi.URLParam(r, "server_id")
	channelId, err1 := strconv.Atoi(channelIdParam)
	serverId, err2 := strconv.Atoi(serverIdParam)
	if err1 != nil || err2 != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	protoMessage := &proto.BodyChannelRemoved{
		ServerId:  int32(serverId),
		ChannelId: int32(channelId),
		UserId:    user.ID,
	}
	serverPID := actors.ServersEngine.Registry.GetPID("server", serverIdParam)
	actors.ServersEngine.Send(serverPID, protoMessage)

	utils.RespondWithJSON(w, http.StatusContinue, &DefaultResponse{Message: "success"})
}
