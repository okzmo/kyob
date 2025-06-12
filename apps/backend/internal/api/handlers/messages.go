package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/okzmo/kyob/db"
	"github.com/okzmo/kyob/internal/api/actors"
	services "github.com/okzmo/kyob/internal/service"
	"github.com/okzmo/kyob/internal/utils"
	proto "github.com/okzmo/kyob/types"
)

func CreateOrEditMessage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(db.User)
	channelId := chi.URLParam(r, "channel_id")
	serverId := chi.URLParam(r, "server_id")

	var body services.MessageBody

	body.Type = r.FormValue("type")
	body.MentionsUsers = r.Form["mentions_users[]"]
	contentJSON := r.FormValue("content")
	if err := json.Unmarshal([]byte(contentJSON), &body.Content); err != nil {
		slog.Error(err.Error())
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid message content.")
		return
	}

	files := r.MultipartForm.File["attachments[]"]
	if len(files) > 0 {
		attachmentService := services.NewAttachmentService()
		attachments, err := attachmentService.ProcessAttachments(files, 15<<20)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		body.Attachments = attachments

		err = validate.Struct(body)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	channelPID := actors.ServersEngine.Registry.GetPID(fmt.Sprintf("server/%s/channel", serverId), channelId)

	switch body.Type {
	case "SEND":
		mess := &proto.IncomingChatMessage{
			AuthorId:      user.ID,
			Content:       body.Content,
			ServerId:      serverId,
			ChannelId:     channelId,
			MentionsUsers: body.MentionsUsers,
			Attachments:   body.Attachments,
		}

		actors.ServersEngine.Send(channelPID, mess)
	case "EDIT":
		messageId := chi.URLParam(r, "message_id")

		mess := &proto.EditChatMessage{
			UserId:        user.ID,
			ServerId:      serverId,
			ChannelId:     channelId,
			MessageId:     messageId,
			Content:       body.Content,
			MentionsUsers: body.MentionsUsers,
		}
		actors.ServersEngine.Send(channelPID, mess)
	}

	utils.RespondWithJSON(w, http.StatusCreated, &DefaultResponse{Message: "success"})
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	serverId := chi.URLParam(r, "server_id")
	channelId := chi.URLParam(r, "channel_id")
	messageId := chi.URLParam(r, "message_id")
	user := r.Context().Value("user").(db.User)

	channelPID := actors.ServersEngine.Registry.GetPID(fmt.Sprintf("server/%s/channel", serverId), channelId)
	mess := &proto.DeleteChatMessage{
		UserId:    user.ID,
		ServerId:  serverId,
		ChannelId: channelId,
		MessageId: messageId,
	}

	actors.ServersEngine.Send(channelPID, mess)

	utils.RespondWithJSON(w, http.StatusCreated, &DefaultResponse{Message: "success"})
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	channelId := chi.URLParam(r, "channel_id")

	messages, err := services.GetMessages(r.Context(), channelId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, messages)
}
