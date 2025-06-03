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
	Email    string `validate:"omitempty,email" json:"email"`
	Username string `validate:"omitempty,min=1,max=20" json:"username"`
}

type Link struct {
	Id    string `json:"id"`
	Label string `json:"label"`
	Url   string `json:"url"`
}

type Fact struct {
	Id    string `json:"id"`
	Label string `json:"label"`
	Value string `json:"value"`
}

type UpdateProfileBody struct {
	DisplayName string          `validate:"required,max=20" json:"display_name"`
	About       json.RawMessage `json:"about"`
	Links       []Link          `json:"links"`
	Facts       []Fact          `json:"facts"`
}

type UpdateAvatarBody struct {
	CropAvatar Crop   `validate:"required" json:"crop_avatar"`
	CropBanner Crop   `validate:"required" json:"crop_banner"`
	MainColor  string `json:"main_color"`
}

type UpdateAvatarResponse struct {
	Avatar    string `json:"avatar"`
	Banner    string `json:"banner"`
	MainColor string `json:"main_color"`
}

func GetUser(ctx context.Context, userId string) (*UserResponse, error) {
	user, err := db.Query.GetUserById(ctx, userId)
	if err != nil {
		return nil, ErrUserNotFound
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
		Links:       user.Links,
		Facts:       user.Facts,
	}

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

func UpdateProfile(ctx context.Context, body *UpdateProfileBody) error {
	user := ctx.Value("user").(db.User)

	if body.DisplayName != "" {
		err := db.Query.UpdateUserDisplayName(ctx, db.UpdateUserDisplayNameParams{
			ID:          user.ID,
			DisplayName: body.DisplayName,
		})
		if err != nil {
			return err
		}
	}

	if len(body.About) > 0 {
		err := db.Query.UpdateUserAbout(ctx, db.UpdateUserAboutParams{
			ID:    user.ID,
			About: body.About,
		})
		if err != nil {
			return err
		}
	}

	links := make([]Link, 0)
	for _, link := range body.Links {
		link.Id = utils.Node.Generate().String()
		links = append(links, link)
	}
	jsonLinks, err := json.Marshal(links)
	if err != nil {
		return err
	}

	err = db.Query.UpdateUserLinks(ctx, db.UpdateUserLinksParams{
		ID:    user.ID,
		Links: jsonLinks,
	})
	if err != nil {
		return err
	}

	facts := make([]Fact, 0)
	for _, fact := range body.Facts {
		fact.Id = utils.Node.Generate().String()
		facts = append(facts, fact)
	}
	jsonFacts, err := json.Marshal(facts)
	if err != nil {
		return err
	}

	err = db.Query.UpdateUserFacts(ctx, db.UpdateUserFactsParams{
		ID:    user.ID,
		Facts: jsonFacts,
	})
	if err != nil {
		return err
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

	randomId := utils.GenerateRandomId(8)
	avatarFileName := fmt.Sprintf("avatar-%s-%s.webp", user.ID, randomId)
	bannerFileName := fmt.Sprintf("banner-%s-%s.webp", user.ID, randomId)

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
	mainColor := pgtype.Text{String: body.MainColor, Valid: true}
	err = db.Query.UpdateUserAvatarNBanner(ctx, db.UpdateUserAvatarNBannerParams{
		ID:        user.ID,
		Avatar:    avatarUrl,
		Banner:    bannerUrl,
		MainColor: mainColor,
	})
	if err != nil {
		return nil, err
	}

	return &UpdateAvatarResponse{
		Banner:    bannerUrl.String,
		Avatar:    avatarUrl.String,
		MainColor: mainColor.String,
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
