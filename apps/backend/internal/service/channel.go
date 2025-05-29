package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/okzmo/kyob/db"
	"github.com/okzmo/kyob/internal/utils"
	proto "github.com/okzmo/kyob/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrUnauthorizedChannelCreation = errors.New("cannot create a channel in this server")
	ErrUnauthorizedChannelEdition  = errors.New("cannot edit this channel")
	ErrUnauthorizedChannelDeletion = errors.New("cannot delete this channel")
)

type CreateChannelBody struct {
	Name        string         `validate:"required,max=50" json:"name"`
	Type        db.ChannelType `validate:"required,oneof=textual voice" json:"type"`
	Description string         `validate:"max=280" json:"description"`
	Users       []string       `json:"users"`
	Roles       []string       `json:"roles"`
	X           int32          `json:"x"`
	Y           int32          `json:"y"`
	Id          *string        `json:"id"`
}

type EditChannelBody struct {
	ServerID    string `validate:"required" json:"server_id"`
	Name        string `validate:"max=50" json:"name"`
	Description string `validate:"max=280" json:"description"`
}

type DeleteChannelBody struct {
	ServerID int `validate:"required" json:"server_id"`
}

func CreateChannel(ctx context.Context, creatorId string, serverId string, channel *CreateChannelBody) (*proto.BroadcastChannelCreation, error) {
	if creatorId != "global" {
		res, err := db.Query.OwnServer(ctx, db.OwnServerParams{
			ID:      serverId,
			OwnerID: creatorId,
		})
		if err != nil || res.RowsAffected() == 0 {
			return nil, ErrUnauthorizedChannelCreation
		}
	}

	channelParams := db.CreateChannelParams{
		ServerID:    serverId,
		Name:        channel.Name,
		Type:        channel.Type,
		Description: pgtype.Text{String: channel.Description, Valid: true},
		Users:       channel.Users,
		Roles:       channel.Roles,
		X:           channel.X,
		Y:           channel.Y,
	}

	if channel.Id != nil {
		channelParams.ID = *channel.Id
	} else {
		channelParams.ID = utils.Node.Generate().String()
	}

	c, err := db.Query.CreateChannel(ctx, channelParams)
	if err != nil {
		return nil, err
	}

	newChannel := &proto.BroadcastChannelCreation{
		Id:          c.ID,
		ServerId:    serverId,
		Name:        c.Name,
		Description: &c.Description.String,
		Type:        string(c.Type),
		X:           c.X,
		Y:           c.Y,
		CreatedAt:   timestamppb.New(c.CreatedAt),
		UpdatedAt:   timestamppb.New(c.UpdatedAt),
	}

	return newChannel, nil
}

func EditChannel(ctx context.Context, id string, body *EditChannelBody) error {
	user := ctx.Value("user").(db.User)
	res, err := db.Query.OwnServer(ctx, db.OwnServerParams{
		ID:      body.ServerID,
		OwnerID: user.ID,
	})
	if err != nil || res.RowsAffected() == 0 {
		return ErrUnauthorizedChannelEdition
	}

	if body.Name != "" {
		err := db.Query.UpdateChannelName(ctx, db.UpdateChannelNameParams{
			ID:   id,
			Name: body.Name,
		})
		if err != nil {
			return err
		}
	}

	if body.Description != "" {
		err := db.Query.UpdateChannelDescription(ctx, db.UpdateChannelDescriptionParams{
			ID:          id,
			Description: pgtype.Text{String: body.Description, Valid: true},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteChannel(ctx context.Context, serverId string, channelId string, userId string) error {
	res, err := db.Query.OwnServer(ctx, db.OwnServerParams{
		ID:      serverId,
		OwnerID: userId,
	})
	if err != nil || res.RowsAffected() == 0 {
		return ErrUnauthorizedChannelDeletion
	}

	err = db.Query.DeleteChannel(ctx, channelId)
	if err != nil {
		return err
	}

	return nil
}
