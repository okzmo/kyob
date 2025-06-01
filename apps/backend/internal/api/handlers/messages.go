package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	channelId := chi.URLParam(r, "channel_id")
	serverId := chi.URLParam(r, "server_id")

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

	channelPID := actors.ServersEngine.Registry.GetPID(fmt.Sprintf("server/%s/channel", serverId), channelId)

	links, err := db.Query.GetUserLinks(r.Context(), user.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	facts, err := db.Query.GetUserFacts(r.Context(), user.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var protoLinks []*proto.UserLinksRow
	for _, link := range links {
		protoLink := &proto.UserLinksRow{
			Id:    link.ID,
			Label: link.Label.String,
			Url:   link.Url.String,
		}
		protoLinks = append(protoLinks, protoLink)
	}

	var protoFacts []*proto.UserFactsRow
	for _, fact := range facts {
		protoFact := &proto.UserFactsRow{
			Id:    fact.ID,
			Label: fact.Label.String,
			Value: fact.Value.String,
		}
		protoFacts = append(protoFacts, protoFact)
	}

	switch body.Type {
	case "SEND":
		mess := &proto.IncomingChatMessage{
			Author: &proto.User{
				Id:          user.ID,
				Email:       user.Email,
				Username:    user.Username,
				DisplayName: user.DisplayName,
				Avatar:      &user.Avatar.String,
				Banner:      &user.Banner.String,
				MainColor:   &user.MainColor.String,
				About:       user.About,
				Links:       protoLinks,
				Facts:       protoFacts,
				CreatedAt:   timestamppb.New(user.CreatedAt),
			},
			Content:       body.Content,
			ServerId:      serverId,
			ChannelId:     channelId,
			MentionsUsers: body.MentionsUsers,
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

	utils.RespondWithJSON(w, http.StatusContinue, messages)
}
