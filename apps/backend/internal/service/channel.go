package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/okzmo/kyob/db"
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
	Users       []int64        `json:"users"`
	Roles       []int64        `json:"roles"`
	X           int32          `json:"x"`
	Y           int32          `json:"y"`
}

type EditChannelBody struct {
	ServerID    int    `validate:"required" json:"server_id"`
	Name        string `validate:"max=50" json:"name"`
	Description string `validate:"max=280" json:"description"`
}

type DeleteChannelBody struct {
	ServerID int `validate:"required" json:"server_id"`
}

func CreateChannel(ctx context.Context, creatorId int64, serverId int32, channel *CreateChannelBody) (*proto.BroadcastChannelCreation, error) {
	res, err := db.Query.OwnServer(ctx, db.OwnServerParams{
		ID:      int64(serverId),
		OwnerID: creatorId,
	})
	if err != nil || res.RowsAffected() == 0 {
		return nil, ErrUnauthorizedChannelCreation
	}

	c, err := db.Query.CreateChannel(ctx, db.CreateChannelParams{
		ServerID:    int64(serverId),
		Name:        channel.Name,
		Type:        channel.Type,
		Description: pgtype.Text{String: channel.Description, Valid: true},
		Users:       channel.Users,
		Roles:       channel.Roles,
		X:           channel.X,
		Y:           channel.Y,
	})
	if err != nil {
		return nil, err
	}

	newChannel := &proto.BroadcastChannelCreation{
		Id:          int32(c.ID),
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

func EditChannel(ctx context.Context, id int, body *EditChannelBody) error {
	user := ctx.Value("user").(db.User)
	res, err := db.Query.OwnServer(ctx, db.OwnServerParams{
		ID:      int64(body.ServerID),
		OwnerID: user.ID,
	})
	if err != nil || res.RowsAffected() == 0 {
		return ErrUnauthorizedChannelEdition
	}

	if body.Name != "" {
		err := db.Query.UpdateChannelName(ctx, db.UpdateChannelNameParams{
			ID:   int64(id),
			Name: body.Name,
		})
		if err != nil {
			return err
		}
	}

	if body.Description != "" {
		err := db.Query.UpdateChannelDescription(ctx, db.UpdateChannelDescriptionParams{
			ID:          int64(id),
			Description: pgtype.Text{String: body.Description, Valid: true},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteChannel(ctx context.Context, serverId int, channelId int, userId int64) error {
	res, err := db.Query.OwnServer(ctx, db.OwnServerParams{
		ID:      int64(serverId),
		OwnerID: userId,
	})
	if err != nil || res.RowsAffected() == 0 {
		return ErrUnauthorizedChannelDeletion
	}

	err = db.Query.DeleteChannel(ctx, int64(channelId))
	if err != nil {
		return err
	}

	return nil
}
