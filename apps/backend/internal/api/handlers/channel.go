package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	services "github.com/okzmo/kyob/internal/service"
	"github.com/okzmo/kyob/internal/utils"
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

	channel, err := services.CreateChannel(r.Context(), serverId, &body)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorizedChannelCreation):
			utils.RespondWithError(w, http.StatusUnauthorized, "You cannot create a channel in this server.")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, channel)
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
	channelIdParam := chi.URLParam(r, "channel_id")
	serverIdParam := chi.URLParam(r, "server_id")
	channelId, _ := strconv.Atoi(channelIdParam)
	serverId, _ := strconv.Atoi(serverIdParam)

	err := services.DeleteChannel(r.Context(), serverId, channelId)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorizedChannelDeletion):
			utils.RespondWithError(w, http.StatusUnauthorized, "You cannot delete this channel.")
		default:
			utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusContinue, &DefaultResponse{Message: "success"})
}

func GetChannels(w http.ResponseWriter, r *http.Request) {
}
