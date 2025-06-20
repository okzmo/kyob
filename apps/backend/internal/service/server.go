package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"
	"os"
	"regexp"
	"strings"
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
	ErrNoIDInInvite               = errors.New("failed to find id in invite url")
	ErrTooManyServers             = errors.New("servers limit")
)

type Crop struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Height int `json:"height"`
	Width  int `json:"width"`
}

type CreateServerBody struct {
	Name        string          `validate:"required,max=50" json:"name"`
	Description json.RawMessage `json:"description"`
	Private     bool            `json:"private"`
	Crop        Crop            `validate:"required" json:"crop"`
	X           int             `validate:"required" json:"x"`
	Y           int             `validate:"required" json:"y"`
}

type UpdateServerProfileBody struct {
	Name        string          `validate:"max=50" json:"name"`
	Description json.RawMessage `json:"description"`
}

type JoinServerBody struct {
	InviteURL string `validate:"required" json:"invite_url"`
	X         int    `validate:"required" json:"x"`
	Y         int    `validate:"required" json:"y"`
}

type ServerResponse struct {
	ID          string          `json:"id"`
	OwnerID     string          `json:"owner_id"`
	Name        string          `json:"name"`
	Avatar      pgtype.Text     `json:"avatar"`
	Banner      pgtype.Text     `json:"banner"`
	Description json.RawMessage `json:"description"`
	MainColor   string          `json:"main_color"`
	Private     bool            `json:"private"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	X           int             `json:"x"`
	Y           int             `json:"y"`
}

type JoinServerResponse struct {
	Server ServerWithChannels `json:"server"`
}

type ServerInviteResponse struct {
	InviteLink string `json:"invite_link"`
}

func CreateServer(ctx context.Context, file []byte, fileHeader *multipart.FileHeader, server *CreateServerBody) (*ServerResponse, error) {
	image, err := utils.CropImage(file, server.Crop.X, server.Crop.Y, server.Crop.Width, server.Crop.Height)
	if err != nil {
		slog.Error("image cropping error", "err", err)
		return nil, err
	}

	id := utils.GenerateRandomId(8)
	imgFileName := fmt.Sprintf("avatar-server-%s.webp", id)

	client := s3.NewFromConfig(GetAWSConfig())
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
	nbServers, _ := db.Query.GetServersCountFromUser(ctx, user.ID)
	if nbServers >= 200 {
		return nil, ErrTooManyServers
	}

	mainColor := pgtype.Text{String: "12,12,16", Valid: true}
	newServer, err := db.Query.CreateServer(ctx, db.CreateServerParams{
		ID:          utils.Node.Generate().String(),
		OwnerID:     user.ID,
		Name:        server.Name,
		Avatar:      pgtype.Text{String: fmt.Sprintf("%s/%s", os.Getenv("CDN_URL"), imgFileName), Valid: true},
		Description: server.Description,
		MainColor:   mainColor,
		Private:     server.Private,
	})
	if err != nil {
		return nil, err
	}

	err = db.Query.JoinServer(ctx, db.JoinServerParams{
		ID:       utils.Node.Generate().String(),
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

func UpdateServerProfile(ctx context.Context, id string, body *UpdateServerProfileBody) error {
	user := ctx.Value("user").(db.User)
	res, err := db.Query.OwnServer(ctx, db.OwnServerParams{
		ID:      id,
		OwnerID: user.ID,
	})
	if err != nil || res.RowsAffected() == 0 {
		return ErrUnauthorizedServerEdition
	}

	if body.Name != "" {
		err := db.Query.UpdateServerName(ctx, db.UpdateServerNameParams{
			ID:      id,
			Name:    body.Name,
			OwnerID: user.ID,
		})
		if err != nil {
			return err
		}
	}

	if len(body.Description) > 0 {
		err := db.Query.UpdateServerDescription(ctx, db.UpdateServerDescriptionParams{
			ID:          id,
			Description: body.Description,
			OwnerID:     user.ID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateServerAvatar(ctx context.Context, serverID string, file []byte, fileHeader *multipart.FileHeader, body *UpdateAvatarBody) (*UpdateAvatarResponse, error) {
	user := ctx.Value("user").(db.User)
	s3Client := s3.NewFromConfig(GetAWSConfig())

	res, err := db.Query.OwnServer(ctx, db.OwnServerParams{
		ID:      serverID,
		OwnerID: user.ID,
	})
	if err != nil || res.RowsAffected() == 0 {
		return nil, ErrUnauthorizedServerEdition
	}

	server, err := db.Query.GetServer(ctx, serverID)
	if err != nil || res.RowsAffected() == 0 {
		return nil, err
	}

	avatar, err := utils.CropImage(file, body.CropAvatar.X, body.CropAvatar.Y, body.CropAvatar.Width, body.CropAvatar.Height)
	if err != nil {
		slog.Error("avatar cropping error", "err", err)
		return nil, err
	}

	banner, err := utils.CropImage(file, body.CropBanner.X, body.CropBanner.Y, body.CropBanner.Width, body.CropBanner.Height)
	if err != nil {
		slog.Error("banner cropping error", "err", err)
		return nil, err
	}

	randomID := utils.GenerateRandomId(8)
	avatarFileName := fmt.Sprintf("avatar-%s-%s.webp", serverID, randomID)
	bannerFileName := fmt.Sprintf("banner-%s-%s.webp", serverID, randomID)
	oldAvatarSplit := strings.Split(server.Avatar.String, "/")
	oldBannerSplit := strings.Split(server.Banner.String, "/")

	// upload new avatar
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Key:    &avatarFileName,
		Bucket: aws.String("nyo-files"),
		Body:   bytes.NewReader(avatar),
	})
	if err != nil {
		slog.Error("failed uploading server avatar", "err", err)
		return nil, err
	}

	// upload new banner
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Key:    &bannerFileName,
		Bucket: aws.String("nyo-files"),
		Body:   bytes.NewReader(banner),
	})
	if err != nil {
		slog.Error("failed uploading server banner", "err", err)
		return nil, err
	}

	// delete old avatar
	if server.Avatar.String != "" {
		_, err = s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Key:    aws.String(oldAvatarSplit[len(oldAvatarSplit)-1]),
			Bucket: aws.String("nyo-files"),
		})
		if err != nil {
			slog.Error("failed deleting user avatar", "err", err)
			return nil, err
		}
	}

	// delete old banner
	if server.Banner.String != "" {
		_, err = s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Key:    aws.String(oldBannerSplit[len(oldBannerSplit)-1]),
			Bucket: aws.String("nyo-files"),
		})
		if err != nil {
			slog.Error("failed deleting user banner", "err", err)
			return nil, err
		}
	}

	avatarURL := pgtype.Text{String: fmt.Sprintf("%s/%s", os.Getenv("CDN_URL"), avatarFileName), Valid: true}
	bannerURL := pgtype.Text{String: fmt.Sprintf("%s/%s", os.Getenv("CDN_URL"), bannerFileName), Valid: true}
	mainColor := pgtype.Text{String: body.MainColor, Valid: true}
	err = db.Query.UpdateServerAvatarNBanner(ctx, db.UpdateServerAvatarNBannerParams{
		ID:        serverID,
		OwnerID:   user.ID,
		Avatar:    avatarURL,
		Banner:    bannerURL,
		MainColor: mainColor,
	})
	if err != nil {
		return nil, err
	}

	return &UpdateAvatarResponse{
		Banner:    bannerURL.String,
		Avatar:    avatarURL.String,
		MainColor: mainColor.String,
	}, nil
}

func DeleteServer(ctx context.Context, id string, userID string) error {
	res, err := db.Query.DeleteServer(ctx, db.DeleteServerParams{
		ID:      id,
		OwnerID: userID,
	})
	if err != nil || res.RowsAffected() == 0 {
		return ErrUnauthorizedServerDeletion
	}

	return nil
}

func CreateServerInvite(ctx context.Context, serverID string) (*string, error) {
	inviteID := utils.GenerateRandomId(10)

	res, err := db.Query.CreateInvite(ctx, db.CreateInviteParams{
		ID:       utils.Node.Generate().String(),
		ServerID: serverID,
		InviteID: inviteID,
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
	matches := pattern.FindStringSubmatch(body.InviteURL)
	if matches == nil {
		return nil, ErrNoIDInInvite
	}

	serverID, err := db.Query.CheckInvite(ctx, matches[1])
	if err != nil {
		return nil, ErrServerNotFound
	}

	err = db.Query.JoinServer(ctx, db.JoinServerParams{
		ID:       utils.Node.Generate().String(),
		UserID:   user.ID,
		ServerID: serverID,
		X:        int32(body.X),
		Y:        int32(body.Y),
	})
	if err != nil {
		return nil, err
	}

	channelMap := make(map[string]ChannelsWithMembers)
	channels, err := db.Query.GetChannelsFromServer(ctx, serverID)
	if err != nil {
		return nil, err
	}

	membersMap := make(map[string]db.GetMembersFromServersRow)
	allMembers, err := db.Query.GetMembersFromServers(ctx, []string{serverID})
	if err != nil {
		return nil, err
	}
	for _, member := range allMembers {
		membersMap[member.ID] = member
	}

	for _, channelRaw := range channels {
		channel := ChannelsWithMembers{
			channelRaw,
			"",
			"",
			json.RawMessage(`[]`),
			[]VoiceUser{},
		}

		channelMap[channel.ID] = channel
	}

	server, err := db.Query.GetServerWithChannels(ctx, db.GetServerWithChannelsParams{
		ServerID: serverID,
		UserID:   user.ID,
	})
	if err != nil {
		return nil, err
	}

	s := ServerWithChannels{
		ServerResponse{
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
		server.X,
		server.Y,
		channelMap,
		int(server.MemberCount),
		allMembers,
	}

	return &s, nil
}

func LeaveServer(ctx context.Context, serverID, userID string) error {
	err := db.Query.LeaveServer(ctx, db.LeaveServerParams{
		UserID:   userID,
		ServerID: serverID,
	})
	if err != nil {
		return err
	}

	return nil
}
