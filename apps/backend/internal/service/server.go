package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/okzmo/kyob/db"
	"github.com/okzmo/kyob/internal/utils"
)

var (
	ErrUnauthorizedServerEdition  = errors.New("cannot edit server")
	ErrUnauthorizedServerDeletion = errors.New("cannot delete this server")
)

type Crop struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Height int `json:"height"`
	Width  int `json:"width"`
}

type CreateServerBody struct {
	Name        string `validate:"required,max=50" json:"name"`
	Description string `validate:"max=280" json:"description"`
	Private     bool   `json:"private"`
	Crop        Crop   `json:"crop"`
}

type EditServerBody struct {
	Name        string `validate:"max=50" json:"name"`
	Avatar      string `json:"avatar"`
	Banner      string `json:"banner"`
	Description string `validate:"max=280" json:"description"`
}

func CreateServer(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, server *CreateServerBody) (*db.Server, error) {
	image, err := utils.CropImage(file, server.Crop.X, server.Crop.Y, server.Crop.Width, server.Crop.Height)
	if err != nil {
		slog.Error("image cropping error", "err", err)
		return nil, err
	}

	id, err := utils.GenerateRandomId(8)
	if err != nil {
		slog.Error("failed generating random id", "err", err)
		return nil, err
	}
	imgFileName := fmt.Sprintf("avatar-server-%s.webp", id)

	client := s3.NewFromConfig(GetS3Config())
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Key:    &imgFileName,
		Bucket: aws.String("nyo-files"),
		Body:   bytes.NewReader(image),
	})
	if err != nil {
		slog.Error("failed uploading server avatar", "err", err)
		return nil, err
	}

	x, y, err := utils.GenerateRandomCoordinates()
	if err != nil {
		slog.Error("failed generating coords", "err", err)
		return nil, err
	}

	user := ctx.Value("user").(db.User)
	newServer, err := db.Query.CreateServer(ctx, db.CreateServerParams{
		OwnerID:     user.ID,
		Name:        server.Name,
		Avatar:      pgtype.Text{String: fmt.Sprintf("%s/%s", os.Getenv("CDN_URL"), imgFileName), Valid: true},
		Description: pgtype.Text{String: server.Description, Valid: true},
		Private:     server.Private,
		X:           x,
		Y:           y,
	})
	if err != nil {
		return nil, err
	}

	err = db.Query.JoinServer(ctx, db.JoinServerParams{
		ServerID: newServer.ID,
		UserID:   user.ID,
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

	if body.Avatar != "" {
		err := db.Query.UpdateServerAvatar(ctx, db.UpdateServerAvatarParams{
			ID:      int64(id),
			Avatar:  pgtype.Text{String: body.Avatar, Valid: true},
			OwnerID: user.ID,
		})
		if err != nil {
			return err
		}
	}

	if body.Banner != "" {
		err := db.Query.UpdateServerBanner(ctx, db.UpdateServerBannerParams{
			ID:      int64(id),
			Banner:  pgtype.Text{String: body.Banner, Valid: true},
			OwnerID: user.ID,
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
