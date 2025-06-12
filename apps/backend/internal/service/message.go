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
	Attachments      json.RawMessage `json:"attachments"`
	Type             string          `json:"type"`
}

type EditMessageBody struct {
	Content json.RawMessage `validate:"required" json:"content"`
}

type MessageResponse struct {
	ID               string          `json:"id"`
	AuthorId         string          `json:"author_id"`
	ServerId         string          `json:"server_id"`
	ChannelId        string          `json:"channel_id"`
	Content          json.RawMessage `json:"content"`
	MentionsUsers    []string        `json:"mentions_users"`
	MentionsChannels []string        `json:"mentions_channels"`
	Attachments      json.RawMessage `json:"attachments"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

func CreateMessage(ctx context.Context, userId string, serverId string, channelId string, body *MessageBody) (*proto.BroadcastChatMessage, error) {
	if serverId != "global" {
		res, err := db.Query.CheckChannelMembership(ctx, db.CheckChannelMembershipParams{
			ID:     channelId,
			UserID: userId,
		})
		if err != nil || res.RowsAffected() == 0 {
			return nil, ErrUnauthorizedMessageCreation
		}
	}

	m, err := db.Query.CreateMessage(ctx, db.CreateMessageParams{
		ID:               utils.Node.Generate().String(),
		AuthorID:         userId,
		ServerID:         serverId,
		ChannelID:        channelId,
		Content:          body.Content,
		MentionsUsers:    body.MentionsUsers,
		MentionsChannels: body.MentionsChannels,
		Attachments:      body.Attachments,
	})
	if err != nil {
		return nil, err
	}

	message := &proto.BroadcastChatMessage{
		Id:               m.ID,
		AuthorId:         userId,
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
	s3Client := s3.NewFromConfig(GetAWSConfig())

	mess, err := db.Query.GetMessage(ctx, messageId)
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
			attachmentSplit := strings.Split(attachment.Url, "/")
			s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Key:    aws.String(attachmentSplit[len(attachmentSplit)-1]),
				Bucket: aws.String("nyo-files"),
			})
		}
	}

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
		messages = append(messages, MessageResponse{
			ID:               message.ID,
			AuthorId:         message.AuthorID,
			ServerId:         message.ServerID,
			ChannelId:        message.ChannelID,
			Content:          message.Content,
			MentionsUsers:    message.MentionsUsers,
			MentionsChannels: message.MentionsChannels,
			Attachments:      message.Attachments,
			UpdatedAt:        message.UpdatedAt,
			CreatedAt:        message.CreatedAt,
		})
	}

	return messages, nil
}
