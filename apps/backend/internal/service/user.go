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
	ErrAddingItself  = errors.New("user adding itself")
	ErrUsernameInUse = errors.New("username in use")
	ErrEmailInUse    = errors.New("email in use")
)

type AddFriendBody struct {
	Username string `validate:"required" json:"username"`
}

type AcceptFriendBody struct {
	FriendshipID string `validate:"required" json:"friendship_id"`
	UserID       string `validate:"required" json:"user_id"`
	FriendID     string `validate:"required" json:"friend_id"`
}

type RemoveFriendBody struct {
	FriendshipID string `validate:"required" json:"friendship_id"`
	FriendID     string `validate:"required" json:"friend_id"`
	UserID       string `validate:"required" json:"user_id"`
}

type UpdateAccountBody struct {
	Email       string `validate:"omitempty,email" json:"email"`
	Username    string `validate:"omitempty,min=1,max=20" json:"username"`
	DisplayName string `validate:"omitempty,min=1,max=20" json:"display_name"`
}

type UpdateAvatarBody struct {
	CropAvatar Crop `validate:"required" json:"crop_avatar"`
	CropBanner Crop `validate:"required" json:"crop_banner"`
}

type UpdateAvatarResponse struct {
	Banner string `json:"avatar"`
	Avatar string `json:"banner"`
}

func GetUser(ctx context.Context, userId string) (*UserResponse, error) {
	user, err := db.Query.GetUserById(ctx, userId)
	if err != nil {
		return nil, ErrUserNotFound
	}

	links, err := db.Query.GetUserLinks(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	facts, err := db.Query.GetUserFacts(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	res := &UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Avatar:      user.Avatar,
		Banner:      user.Banner,
		MainColor:   user.MainColor,
		About:       user.About,
		CreatedAt:   user.CreatedAt,
		Links:       links,
		Facts:       facts,
	}

	fmt.Println(res)

	return res, nil
}

func UpdateAccount(ctx context.Context, body *UpdateAccountBody) error {
	user := ctx.Value("user").(db.User)

	if body.Username != "" {
		_, err := db.Query.GetUser(ctx, db.GetUserParams{
			Username: body.Username,
		})
		if err == nil {
			return ErrUsernameInUse
		}

		res, err := db.Query.UpdateUserUsername(ctx, db.UpdateUserUsernameParams{
			ID:       user.ID,
			Username: body.Username,
		})
		fmt.Println(res)
		if err != nil {
			return err
		}
	}

	if body.DisplayName != "" {
		err := db.Query.UpdateUserDisplayName(ctx, db.UpdateUserDisplayNameParams{
			ID:          user.ID,
			DisplayName: body.DisplayName,
		})
		if err != nil {
			return err
		}
	}

	if body.Email != "" {
		_, err := db.Query.GetUser(ctx, db.GetUserParams{
			Email: body.Email,
		})
		if err == nil {
			return ErrEmailInUse
		}

		err = db.Query.UpdateUserEmail(ctx, db.UpdateUserEmailParams{
			ID:    user.ID,
			Email: body.Email,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateAvatar(ctx context.Context, file []byte, fileHeader *multipart.FileHeader, body *UpdateAvatarBody) (*UpdateAvatarResponse, error) {
	user := ctx.Value("user").(db.User)
	client := s3.NewFromConfig(GetS3Config())

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

	avatarFileName := fmt.Sprintf("avatar-%s.webp", user.ID)
	bannerFileName := fmt.Sprintf("banner-%s.webp", user.ID)

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Key:    &avatarFileName,
		Bucket: aws.String("nyo-files"),
		Body:   bytes.NewReader(avatar),
	})
	if err != nil {
		slog.Error("failed uploading user avatar", "err", err)
		return nil, err
	}

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Key:    &bannerFileName,
		Bucket: aws.String("nyo-files"),
		Body:   bytes.NewReader(banner),
	})
	if err != nil {
		slog.Error("failed uploading user banner", "err", err)
		return nil, err
	}

	avatarUrl := pgtype.Text{String: fmt.Sprintf("%s/%s", os.Getenv("CDN_URL"), avatarFileName), Valid: true}
	bannerUrl := pgtype.Text{String: fmt.Sprintf("%s/%s", os.Getenv("CDN_URL"), bannerFileName), Valid: true}
	err = db.Query.UpdateUserAvatarNBanner(ctx, db.UpdateUserAvatarNBannerParams{
		ID:     user.ID,
		Avatar: avatarUrl,
		Banner: bannerUrl,
	})
	if err != nil {
		return nil, err
	}

	return &UpdateAvatarResponse{
		Banner: bannerUrl.String,
		Avatar: avatarUrl.String,
	}, nil
}

func AddFriend(ctx context.Context, body *AddFriendBody) (string, string, error) {
	user := ctx.Value("user").(db.User)

	friend, err := db.Query.GetUser(ctx, db.GetUserParams{
		Username: body.Username,
	})
	if err != nil {
		return "", "", ErrUserNotFound
	}

	if user.ID == friend.ID {
		return "", "", ErrAddingItself
	}

	invite, err := db.Query.AddFriend(ctx, db.AddFriendParams{
		ID:       utils.Node.Generate().String(),
		UserID:   user.ID,
		FriendID: friend.ID,
	})
	if err != nil {
		return "", "", err
	}

	return invite.ID, friend.ID, nil
}

func AcceptFriend(ctx context.Context, body *AcceptFriendBody) (*db.User, *db.Channel, error) {
	err := db.Query.AcceptFriend(ctx, body.FriendshipID)
	if err != nil {
		return nil, nil, err
	}

	friend, err := db.Query.GetUserById(ctx, body.UserID)
	if err != nil {
		return nil, nil, err
	}

	existingChannel, err := db.Query.GetExistingChannel(ctx, db.GetExistingChannelParams{
		Column1: body.FriendID,
		Column2: body.UserID,
	})
	if err != nil {
		return &friend, nil, nil
	}

	return &friend, &existingChannel, nil
}

func DeleteFriend(ctx context.Context, body *RemoveFriendBody) (string, error) {
	err := db.Query.DeleteFriend(ctx, body.FriendshipID)
	if err != nil {
		return "", err
	}

	channel, err := db.Query.DeactivateChannel(ctx, db.DeactivateChannelParams{
		Column1: body.FriendID,
		Column2: body.FriendID,
	})
	if err != nil {
		return "", err
	}

	return channel.ID, nil
}
