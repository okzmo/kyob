package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"
	"os"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/okzmo/kyob/db"
	"github.com/okzmo/kyob/internal/utils"
)

var (
	ErrUnauthorizedServerEdition  = errors.New("cannot edit server")
	ErrUnauthorizedServerDeletion = errors.New("cannot delete this server")
	ErrServerNotFound             = errors.New("server not found")
	ErrNoIdInInvite               = errors.New("failed to find id in invite url")
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
	Crop        Crop   `validate:"required" json:"crop"`
	X           int    `validate:"required" json:"x"`
	Y           int    `validate:"required" json:"y"`
}

type EditServerBody struct {
	Name        string `validate:"max=50" json:"name"`
	Avatar      string `json:"avatar"`
	Banner      string `json:"banner"`
	Description string `validate:"max=280" json:"description"`
}

type JoinServerBody struct {
	InviteUrl string `validate:"required" json:"invite_url"`
	X         int    `validate:"required" json:"x"`
	Y         int    `validate:"required" json:"y"`
}

type ServerResponse struct {
	ID          int64       `json:"id"`
	OwnerID     int64       `json:"owner_id"`
	Name        string      `json:"name"`
	Avatar      pgtype.Text `json:"avatar"`
	Banner      pgtype.Text `json:"banner"`
	Description pgtype.Text `json:"description"`
	Private     bool        `json:"private"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	X           int         `json:"x"`
	Y           int         `json:"y"`
}

type JoinServerResponse struct {
	Server ServerWithChannels `json:"server"`
}

type ServerInviteResponse struct {
	InviteLink string `json:"invite_link"`
}

func CreateServer(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, server *CreateServerBody) (*ServerResponse, error) {
	image, err := utils.CropImage(file, server.Crop.X, server.Crop.Y, server.Crop.Width, server.Crop.Height)
	if err != nil {
		slog.Error("image cropping error", "err", err)
		return nil, err
	}

	id := utils.GenerateRandomId(8)
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

	user := ctx.Value("user").(db.User)
	newServer, err := db.Query.CreateServer(ctx, db.CreateServerParams{
		OwnerID:     user.ID,
		Name:        server.Name,
		Avatar:      pgtype.Text{String: fmt.Sprintf("%s/%s", os.Getenv("CDN_URL"), imgFileName), Valid: true},
		Description: pgtype.Text{String: server.Description, Valid: true},
		Private:     server.Private,
	})
	if err != nil {
		return nil, err
	}

	err = db.Query.JoinServer(ctx, db.JoinServerParams{
		ServerID: newServer.ID,
		UserID:   user.ID,
		X:        int32(server.X),
		Y:        int32(server.Y),
	})
	if err != nil {
		return nil, err
	}

	return &ServerResponse{
		ID:          newServer.ID,
		OwnerID:     newServer.OwnerID,
		Name:        newServer.Name,
		Avatar:      newServer.Avatar,
		Banner:      newServer.Banner,
		Description: newServer.Description,
		Private:     newServer.Private,
		CreatedAt:   newServer.CreatedAt,
		UpdatedAt:   newServer.UpdatedAt,
		X:           server.X,
		Y:           server.Y,
	}, nil
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

func DeleteServer(ctx context.Context, id int, userId int64) error {
	res, err := db.Query.DeleteServer(ctx, db.DeleteServerParams{
		ID:      int64(id),
		OwnerID: userId,
	})
	if err != nil || res.RowsAffected() == 0 {
		return ErrUnauthorizedServerDeletion
	}

	return nil
}

func CreateServerInvite(ctx context.Context, serverId int) (*string, error) {
	inviteId := utils.GenerateRandomId(10)

	res, err := db.Query.CreateInvite(ctx, db.CreateInviteParams{
		ServerID: int64(serverId),
		InviteID: inviteId,
		ExpireAt: time.Now().Add(time.Minute * 15),
	})
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func JoinServer(ctx context.Context, body JoinServerBody) (*ServerWithChannels, error) {
	user := ctx.Value("user").(db.User)

	pattern := regexp.MustCompile(`^(?:https:\/\/kyob\.app\/invite\/|)([a-zA-Z0-9]{10})$`)
	matches := pattern.FindStringSubmatch(body.InviteUrl)
	if matches == nil {
		return nil, ErrNoIdInInvite
	}

	serverId, err := db.Query.CheckInvite(ctx, matches[1])
	if err != nil {
		return nil, ErrServerNotFound
	}

	err = db.Query.JoinServer(ctx, db.JoinServerParams{
		UserID:   user.ID,
		ServerID: serverId,
		X:        int32(body.X),
		Y:        int32(body.Y),
	})
	if err != nil {
		return nil, err
	}

	channelMap := make(map[int64]db.Channel)
	channels, err := db.Query.GetChannelsFromServer(ctx, serverId)
	if err != nil {
		return nil, err
	}

	for _, channel := range channels {
		channelMap[channel.ID] = channel
	}

	server, err := db.Query.GetServerWithChannels(ctx, db.GetServerWithChannelsParams{
		ServerID: serverId,
		UserID:   user.ID,
	})
	if err != nil {
		return nil, err
	}

	users, err := db.Query.GetServerMembers(ctx, server.ID)
	if err != nil {
		return nil, err
	}

	s := ServerWithChannels{
		db.Server{
			ID:          server.ID,
			OwnerID:     server.OwnerID,
			Name:        server.Name,
			Avatar:      server.Avatar,
			Banner:      server.Banner,
			Description: server.Description,
			Private:     server.Private,
			CreatedAt:   server.CreatedAt,
			UpdatedAt:   server.UpdatedAt,
		},
		int(server.X),
		int(server.Y),
		channelMap,
		int(server.MemberCount),
		users,
	}

	return &s, nil
}

func LeaveServer(ctx context.Context, serverId int, userId int64) error {
	err := db.Query.LeaveServer(ctx, db.LeaveServerParams{
		UserID:   userId,
		ServerID: int64(serverId),
	})
	if err != nil {
		return err
	}

	return nil
}
