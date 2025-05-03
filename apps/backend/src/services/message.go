package services

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/okzmo/kyob/db"
)

var (
	ErrUnauthorizedMessageEdition  = errors.New("unauthorized message edition")
	ErrUnauthorizedMessageDeletion = errors.New("unauthorized message deletion")
)

type CreateMessageBody struct {
	Content          json.RawMessage `validate:"required" json:"content"`
	MentionsUsers    []int64         `json:"mentions_users"`
	MentionsChannels []int64         `json:"mentions_channels"`
}

type EditMessageBody struct {
	Content json.RawMessage `validate:"required" json:"content"`
}

func CreateMessage(ctx context.Context, channelId int, body *CreateMessageBody) error {
	user := ctx.Value("user").(db.User)

	_, err := db.Query.CreateMessage(ctx, db.CreateMessageParams{
		AuthorID:         user.ID,
		ChannelID:        int64(channelId),
		Content:          body.Content,
		MentionsUsers:    body.MentionsUsers,
		MentionsChannels: body.MentionsChannels,
	})
	if err != nil {
		return err
	}

	return nil
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
