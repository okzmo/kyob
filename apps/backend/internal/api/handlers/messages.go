package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/okzmo/kyob/db"
	"github.com/okzmo/kyob/internal/api/actors"
	services "github.com/okzmo/kyob/internal/service"
	"github.com/okzmo/kyob/internal/utils"
	proto "github.com/okzmo/kyob/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CreateOrEditMessage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(db.User)
	channelIdParam := chi.URLParam(r, "channel_id")
	serverIdParam := chi.URLParam(r, "server_id")
	channelId, err1 := strconv.Atoi(channelIdParam)
	serverId, err2 := strconv.Atoi(serverIdParam)
	if err1 != nil || err2 != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var body services.MessageBody

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

	channelPID := actors.ServersEngine.Registry.GetPID(fmt.Sprintf("server/%d/channel", serverId), strconv.Itoa(channelId))

	switch body.Type {
	case "SEND":
		mess := &proto.IncomingChatMessage{
			Author: &proto.User{
				Id:             int32(user.ID),
				Email:          user.Email,
				Username:       user.Username,
				DisplayName:    user.DisplayName,
				Avatar:         &user.Avatar.String,
				Banner:         &user.Banner.String,
				GradientTop:    &user.GradientTop.String,
				GradientBottom: &user.GradientBottom.String,
				About:          &user.About.String,
				CreatedAt:      timestamppb.New(user.CreatedAt),
			},
			Content:       body.Content,
			ServerId:      int32(serverId),
			ChannelId:     int32(channelId),
			MentionsUsers: body.MentionsUsers,
		}

		actors.ServersEngine.Send(channelPID, mess)
	case "EDIT":
		messageIdParam := chi.URLParam(r, "message_id")
		messageId, err := strconv.Atoi(messageIdParam)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid ID format")
			return
		}

		mess := &proto.EditChatMessage{
			UserId:        user.ID,
			ServerId:      int32(serverId),
			ChannelId:     int32(channelId),
			MessageId:     int32(messageId),
			Content:       body.Content,
			MentionsUsers: body.MentionsUsers,
		}
		actors.ServersEngine.Send(channelPID, mess)
	}

	utils.RespondWithJSON(w, http.StatusCreated, &DefaultResponse{Message: "success"})
}

// func EditMessage(w http.ResponseWriter, r *http.Request) {
// 	idParam := chi.URLParam(r, "id")
// 	id, err := strconv.Atoi(idParam)
// 	if err != nil {
// 		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
//
// 	var body services.EditMessageBody
//
// 	err = json.NewDecoder(r.Body).Decode(&body)
// 	if err != nil {
// 		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
//
// 	err = validate.Struct(body)
// 	if err != nil {
// 		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
//
// 	err = services.EditMessage(r.Context(), id, &body)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, services.ErrUnauthorizedServerEdition):
// 			utils.RespondWithError(w, http.StatusUnauthorized, "You cannot edit a message which is not yours.")
// 		default:
// 			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
// 		}
// 		return
// 	}
//
// 	utils.RespondWithJSON(w, http.StatusCreated, &DefaultResponse{Message: "success"})
// }

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	serverParam := chi.URLParam(r, "server_id")
	channelParam := chi.URLParam(r, "channel_id")
	messageParam := chi.URLParam(r, "message_id")
	serverId, err1 := strconv.Atoi(serverParam)
	channelId, err2 := strconv.Atoi(channelParam)
	messageId, err3 := strconv.Atoi(messageParam)
	user := r.Context().Value("user").(db.User)
	if err1 != nil || err2 != nil || err3 != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	channelPID := actors.ServersEngine.Registry.GetPID(fmt.Sprintf("server/%d/channel", serverId), strconv.Itoa(channelId))
	mess := &proto.DeleteChatMessage{
		UserId:    user.ID,
		ServerId:  int32(serverId),
		ChannelId: int32(channelId),
		MessageId: int32(messageId),
	}

	actors.ServersEngine.Send(channelPID, mess)

	// err := services.DeleteMessage(r.Context(), id)
	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, services.ErrUnauthorizedMessageDeletion):
	// 		utils.RespondWithError(w, http.StatusUnauthorized, "You cannot delete this message.")
	// 	default:
	// 		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	// 	}
	// 	return
	// }

	utils.RespondWithJSON(w, http.StatusCreated, &DefaultResponse{Message: "success"})
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "channel_id")
	channelId, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	messages, err := services.GetMessages(r.Context(), channelId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusContinue, messages)
}
