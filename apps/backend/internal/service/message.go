package services

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/okzmo/kyob/db"
	queries "github.com/okzmo/kyob/db/gen_queries"
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
	Everyone         bool            `json:"everyone"`
	MentionsUsers    []string        `json:"mentions_users"`
	MentionsChannels []string        `json:"mentions_channels"`
	Attachments      json.RawMessage `json:"attachments"`
	Type             string          `json:"type"`
}

type EditMessageBody struct {
	Content json.RawMessage `validate:"required" json:"content"`
}

type MessageResponse struct {
	ID               string          `json:"id"`
	AuthorID         string          `json:"author_id"`
	ServerID         string          `json:"server_id"`
	ChannelID        string          `json:"channel_id"`
	Content          json.RawMessage `json:"content"`
	Everyone         bool            `json:"everyone"`
	MentionsUsers    []string        `json:"mentions_users"`
	MentionsChannels []string        `json:"mentions_channels"`
	Attachments      json.RawMessage `json:"attachments"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

func CreateMessage(ctx context.Context, userID, serverID, channelID string, body *MessageBody) (*proto.BroadcastChatMessage, error) {
	if serverID != "global" {
		res, err := db.Query.CheckChannelMembership(ctx, queries.CheckChannelMembershipParams{
			ID:     channelID,
			UserID: userID,
		})
		if err != nil || res.RowsAffected() == 0 {
			return nil, ErrUnauthorizedMessageCreation
		}
	}

	m, err := db.Query.CreateMessage(ctx, queries.CreateMessageParams{
		ID:               utils.Node.Generate().String(),
		AuthorID:         userID,
		ServerID:         serverID,
		ChannelID:        channelID,
		Content:          body.Content,
		Everyone:         body.Everyone,
		MentionsUsers:    body.MentionsUsers,
		MentionsChannels: body.MentionsChannels,
		Attachments:      body.Attachments,
	})
	if err != nil {
		return nil, err
	}

	message := &proto.BroadcastChatMessage{
		Id:               m.ID,
		AuthorId:         userID,
		ServerId:         m.ServerID,
		ChannelId:        m.ChannelID,
		Content:          m.Content,
		Everyone:         body.Everyone,
		MentionsUsers:    body.MentionsUsers,
		MentionsChannels: body.MentionsChannels,
		Attachments:      body.Attachments,
		CreatedAt:        timestamppb.New(m.CreatedAt),
	}
	return message, nil
}

func EditMessage(ctx context.Context, userID, serverID, channelID, messageID string, body *MessageBody) (*proto.BroadcastEditMessage, error) {
	res, err := db.Query.UpdateMessage(ctx, queries.UpdateMessageParams{
		ID:               messageID,
		Everyone:         body.Everyone,
		MentionsUsers:    body.MentionsUsers,
		MentionsChannels: body.MentionsChannels,
		Content:          body.Content,
		AuthorID:         userID,
	})
	if err != nil || res.RowsAffected() == 0 {
		return nil, ErrUnauthorizedMessageEdition
	}

	message := &proto.BroadcastEditMessage{
		MessageId:        messageID,
		ServerId:         serverID,
		ChannelId:        channelID,
		Content:          body.Content,
		Everyone:         body.Everyone,
		MentionsUsers:    body.MentionsUsers,
		MentionsChannels: body.MentionsChannels,
		UpdatedAt:        timestamppb.New(time.Now()),
	}

	return message, nil
}

func DeleteMessage(ctx context.Context, messageID, userID string) error {
	s3Client := s3.NewFromConfig(GetAWSConfig())

	mess, err := db.Query.GetMessage(ctx, messageID)
	if err != nil {
		return err
	}

	if len(mess.Attachments) > 0 {
		var attachments []Attachment
		err := json.Unmarshal(mess.Attachments, &attachments)
		if err != nil {
			return err
		}

		for _, attachment := range attachments {
			attachmentSplit := strings.Split(attachment.URL, "/")
			s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Key:    aws.String(attachmentSplit[len(attachmentSplit)-1]),
				Bucket: aws.String("nyo-files"),
			})
		}
	}

	res, err := db.Query.DeleteMessage(ctx, queries.DeleteMessageParams{
		ID:       messageID,
		AuthorID: userID,
	})
	if err != nil || res.RowsAffected() == 0 {
		return ErrUnauthorizedMessageDeletion
	}

	return nil
}

func GetMessages(ctx context.Context, channelID string) ([]MessageResponse, error) {
	var messages []MessageResponse

	m, err := db.Query.GetMessagesFromChannel(ctx, channelID)
	if err != nil {
		return nil, err
	}

	for _, message := range m {
		messages = append(messages, MessageResponse{
			ID:               message.ID,
			AuthorID:         message.AuthorID,
			ServerID:         message.ServerID,
			ChannelID:        message.ChannelID,
			Content:          message.Content,
			Everyone:         message.Everyone,
			MentionsUsers:    message.MentionsUsers,
			MentionsChannels: message.MentionsChannels,
			Attachments:      message.Attachments,
			UpdatedAt:        message.UpdatedAt,
			CreatedAt:        message.CreatedAt,
		})
	}

	return messages, nil
}
