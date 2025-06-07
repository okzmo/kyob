package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/okzmo/kyob/db"
	"github.com/okzmo/kyob/internal/utils"
	proto "github.com/okzmo/kyob/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrUnauthorizedMessageCreation = errors.New("unauthorized message creation")
	ErrUnauthorizedMessageEdition  = errors.New("unauthorized message edition")
	ErrUnauthorizedMessageDeletion = errors.New("unauthorized message deletion")
)

type MessageBody struct {
	Content          json.RawMessage `validate:"required" json:"content"`
	MentionsUsers    []string        `json:"mentions_users"`
	MentionsChannels []string        `json:"mentions_channels"`
	Attachments      []string        `json:"attachments"`
	Type             string          `json:"type"`
}

type EditMessageBody struct {
	Content json.RawMessage `validate:"required" json:"content"`
}

type MessageResponse struct {
	ID               string          `json:"id"`
	Author           UserResponse    `json:"author"`
	ServerId         string          `json:"server_id"`
	ChannelId        string          `json:"channel_id"`
	Content          json.RawMessage `json:"content"`
	MentionsUsers    []string        `json:"mentions_users"`
	MentionsChannels []string        `json:"mentions_channels"`
	Attachments      []string        `json:"attachments"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

func CreateMessage(ctx context.Context, user *proto.User, serverId string, channelId string, body *MessageBody) (*proto.BroadcastChatMessage, error) {
	if serverId != "global" {
		res, err := db.Query.CheckChannelMembership(ctx, db.CheckChannelMembershipParams{
			ID:     channelId,
			UserID: user.Id,
		})
		if err != nil || res.RowsAffected() == 0 {
			return nil, ErrUnauthorizedMessageCreation
		}
	}

	m, err := db.Query.CreateMessage(ctx, db.CreateMessageParams{
		ID:               utils.Node.Generate().String(),
		AuthorID:         user.Id,
		ServerID:         serverId,
		ChannelID:        channelId,
		Content:          body.Content,
		MentionsUsers:    body.MentionsUsers,
		MentionsChannels: body.MentionsChannels,
		Attached:         body.Attachments,
	})
	if err != nil {
		return nil, err
	}

	message := &proto.BroadcastChatMessage{
		Id:               m.ID,
		Author:           user,
		ServerId:         m.ServerID,
		ChannelId:        m.ChannelID,
		Content:          m.Content,
		MentionsUsers:    body.MentionsUsers,
		MentionsChannels: body.MentionsChannels,
		Attachments:      body.Attachments,
		CreatedAt:        timestamppb.New(m.CreatedAt),
	}
	return message, nil
}

func EditMessage(ctx context.Context, userId string, serverId string, channelId string, messageId string, body *MessageBody) (*proto.BroadcastEditMessage, error) {
	res, err := db.Query.UpdateMessage(ctx, db.UpdateMessageParams{
		ID:               messageId,
		MentionsUsers:    body.MentionsUsers,
		MentionsChannels: body.MentionsChannels,
		Content:          body.Content,
		AuthorID:         userId,
	})
	if err != nil || res.RowsAffected() == 0 {
		return nil, ErrUnauthorizedMessageEdition
	}

	message := &proto.BroadcastEditMessage{
		MessageId:        messageId,
		ServerId:         serverId,
		ChannelId:        channelId,
		Content:          body.Content,
		MentionsUsers:    body.MentionsUsers,
		MentionsChannels: body.MentionsChannels,
		UpdatedAt:        timestamppb.New(time.Now()),
	}

	return message, nil
}

func DeleteMessage(ctx context.Context, messageId string, userId string) error {
	res, err := db.Query.DeleteMessage(ctx, db.DeleteMessageParams{
		ID:       messageId,
		AuthorID: userId,
	})
	if err != nil || res.RowsAffected() == 0 {
		return ErrUnauthorizedMessageDeletion
	}

	return nil
}

func GetMessages(ctx context.Context, channelId string) ([]MessageResponse, error) {
	var messages []MessageResponse

	m, err := db.Query.GetMessagesFromChannel(ctx, channelId)
	if err != nil {
		return nil, err
	}

	for _, message := range m {
		author, err := db.Query.GetUserById(ctx, message.AuthorID)
		if err != nil {
			return nil, err
		}

		messages = append(messages, MessageResponse{
			ID: message.ID,
			Author: UserResponse{
				ID:          author.ID,
				Username:    author.Username,
				DisplayName: author.DisplayName,
				Avatar:      author.Avatar,
				Banner:      author.Banner,
				About:       author.About,
				Links:       author.Links,
				Facts:       author.Facts,
				MainColor:   author.MainColor,
			},
			ServerId:         message.ServerID,
			ChannelId:        message.ChannelID,
			Content:          message.Content,
			MentionsUsers:    message.MentionsUsers,
			MentionsChannels: message.MentionsChannels,
			Attachments:      message.Attached,
			UpdatedAt:        message.UpdatedAt,
			CreatedAt:        message.CreatedAt,
		})
	}

	return messages, nil
}
