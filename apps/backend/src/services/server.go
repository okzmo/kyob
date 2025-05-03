package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/okzmo/kyob/db"
)

var (
	ErrUnauthorizedServerEdition  = errors.New("cannot edit server")
	ErrUnauthorizedServerDeletion = errors.New("cannot delete this server")
)

type CreateServerBody struct {
	Name        string `validate:"required,max=50" json:"name"`
	Background  string `validate:"required" json:"background"`
	Description string `json:"description"`
	X           int32  `json:"x"`
	Y           int32  `json:"y"`
}

type EditServerBody struct {
	Name        string `validate:"max=50" json:"name"`
	Background  string `json:"background"`
	Description string `validate:"max=280" json:"description"`
}

func CreateServer(ctx context.Context, server *CreateServerBody) (*db.Server, error) {
	user := ctx.Value("user").(db.User)
	newServer, err := db.Query.CreateServer(ctx, db.CreateServerParams{
		OwnerID:     user.ID,
		Name:        server.Name,
		Background:  server.Background,
		Description: pgtype.Text{String: server.Description, Valid: true},
		X:           server.X,
		Y:           server.Y,
	})
	if err != nil {
		return nil, err
	}

	return &newServer, nil
}

func EditServer(ctx context.Context, id int, body *EditServerBody) error {
	user := ctx.Value("user").(db.User)
	res, err := db.Query.OwnServer(ctx, db.OwnServerParams{
		ID:      int64(id),
		OwnerID: user.ID,
	})
	if err != nil || res.RowsAffected() == 0 {
		return ErrUnauthorizedServerEdition
	}

	if body.Name != "" {
		err := db.Query.UpdateServerName(ctx, db.UpdateServerNameParams{
			ID:      int64(id),
			Name:    body.Name,
			OwnerID: user.ID,
		})
		if err != nil {
			return err
		}
	}

	if body.Description != "" {
		err := db.Query.UpdateServerDescription(ctx, db.UpdateServerDescriptionParams{
			ID:          int64(id),
			Description: pgtype.Text{String: body.Description, Valid: true},
			OwnerID:     user.ID,
		})
		if err != nil {
			return err
		}
	}

	if body.Background != "" {
		err := db.Query.UpdateServerBackground(ctx, db.UpdateServerBackgroundParams{
			ID:         int64(id),
			Background: body.Background,
			OwnerID:    user.ID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteServer(ctx context.Context, id int) error {
	user := ctx.Value("user").(db.User)

	res, err := db.Query.DeleteServer(ctx, db.DeleteServerParams{
		ID:      int64(id),
		OwnerID: user.ID,
	})
	if err != nil || res.RowsAffected() == 0 {
		return ErrUnauthorizedServerDeletion
	}

	return nil
}
