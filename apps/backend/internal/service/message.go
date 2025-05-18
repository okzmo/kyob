package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/okzmo/kyob/db"
	proto "github.com/okzmo/kyob/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrUnauthorizedMessageCreation = errors.New("unauthorized message creation")
	ErrUnauthorizedMessageEdition  = errors.New("unauthorized message edition")
	ErrUnauthorizedMessageDeletion = errors.New("unauthorized message deletion")
)

type CreateMessageBody struct {
	Content          json.RawMessage `validate:"required" json:"content"`
	MentionsUsers    []int32         `json:"mentions_users"`
	MentionsChannels []int32         `json:"mentions_channels"`
}

type EditMessageBody struct {
	Content json.RawMessage `validate:"required" json:"content"`
}

type MessageResponse struct {
	ID               int64           `json:"id"`
	Author           db.User         `json:"author"`
	ServerId         int64           `json:"server_id"`
	ChannelId        int64           `json:"channel_id"`
	Content          json.RawMessage `json:"content"`
	MentionsUsers    []int64         `json:"mentions_users"`
	MentionsChannels []int64         `json:"mentions_channels"`
	CreatedAt        time.Time       `json:"created_at"`
}

func CreateMessage(ctx context.Context, user *proto.User, serverId int32, channelId int32, body *CreateMessageBody) (*proto.BroadcastChatMessage, error) {
	res, err := db.Query.CheckChannelMembership(ctx, db.CheckChannelMembershipParams{
		ID:     int64(channelId),
		UserID: int64(user.Id),
	})
	if err != nil || res.RowsAffected() == 0 {
		return nil, ErrUnauthorizedMessageCreation
	}

	convertedMentionsUsers := make([]int64, len(body.MentionsUsers))
	fmt.Println(convertedMentionsUsers)
	for i, v := range body.MentionsUsers {
		convertedMentionsUsers[i] = int64(v)
	}

	convertedMentionsChannels := make([]int64, len(body.MentionsChannels))
	for i, v := range body.MentionsChannels {
		convertedMentionsChannels[i] = int64(v)
	}

	m, err := db.Query.CreateMessage(ctx, db.CreateMessageParams{
		AuthorID:         int64(user.Id),
		ServerID:         int64(serverId),
		ChannelID:        int64(channelId),
		Content:          body.Content,
		MentionsUsers:    convertedMentionsUsers,
		MentionsChannels: convertedMentionsChannels,
	})
	if err != nil {
		return nil, err
	}

	message := &proto.BroadcastChatMessage{
		Id:               int32(m.ID),
		Author:           user,
		ServerId:         int32(m.ServerID),
		ChannelId:        int32(m.ChannelID),
		Content:          m.Content,
		MentionsUsers:    body.MentionsUsers,
		MentionsChannels: body.MentionsChannels,
		CreatedAt:        timestamppb.New(m.CreatedAt),
	}

	return message, nil
}

func EditMessage(ctx context.Context, messageId int, body *EditMessageBody) error {
	user := ctx.Value("user").(db.User)

	if len(body.Content) > 0 {
		res, err := db.Query.UpdateMessageContent(ctx, db.UpdateMessageContentParams{
			ID:       int64(messageId),
			Content:  body.Content,
			AuthorID: user.ID,
		})
		if err != nil || res.RowsAffected() == 0 {
			return ErrUnauthorizedMessageEdition
		}
	}

	return nil
}

func DeleteMessage(ctx context.Context, messageId int) error {
	user := ctx.Value("user").(db.User)

	res, err := db.Query.DeleteMessage(ctx, db.DeleteMessageParams{
		ID:       int64(messageId),
		AuthorID: user.ID,
	})
	if err != nil || res.RowsAffected() == 0 {
		return ErrUnauthorizedMessageDeletion
	}

	return nil
}

func GetMessages(ctx context.Context, channelId int) ([]MessageResponse, error) {
	var messages []MessageResponse

	m, err := db.Query.GetMessagesFromChannel(ctx, int64(channelId))
	if err != nil {
		return nil, err
	}

	for _, message := range m {
		author, err := db.Query.GetUserById(ctx, message.AuthorID)
		if err != nil {
			return nil, err
		}

		messages = append(messages, MessageResponse{
			ID:               message.ID,
			Author:           author,
			ServerId:         message.ServerID,
			ChannelId:        message.ChannelID,
			Content:          message.Content,
			MentionsUsers:    message.MentionsUsers,
			MentionsChannels: message.MentionsChannels,
			CreatedAt:        message.CreatedAt,
		})
	}

	return messages, nil
}
