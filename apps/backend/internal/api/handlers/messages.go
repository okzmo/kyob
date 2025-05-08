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

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	channelIdParam := chi.URLParam(r, "channel_id")
	serverIdParam := chi.URLParam(r, "server_id")
	channelId, _ := strconv.Atoi(channelIdParam)
	serverId, _ := strconv.Atoi(serverIdParam)

	var body services.CreateMessageBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = validate.Struct(body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	message, err := services.CreateMessage(r.Context(), serverId, channelId, &body)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorizedMessageCreation):
			utils.RespondWithError(w, http.StatusUnauthorized, "You cannot send a message in this channel.")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, message)
}

func EditMessage(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var body services.EditMessageBody

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

	err = services.EditMessage(r.Context(), id, &body)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorizedServerEdition):
			utils.RespondWithError(w, http.StatusUnauthorized, "You cannot edit a message which is not yours.")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, &DefaultResponse{Message: "success"})
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = services.DeleteMessage(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorizedMessageDeletion):
			utils.RespondWithError(w, http.StatusUnauthorized, "You cannot delete this message.")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, &DefaultResponse{Message: "success"})
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "channel_id")
	channelId, _ := strconv.Atoi(idParam)

	messages, err := services.GetMessages(r.Context(), channelId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusContinue, messages)
}
