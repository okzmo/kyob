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
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/okzmo/kyob/db"
	"github.com/okzmo/kyob/internal/utils"
)

var (
	ErrAddingItself              = errors.New("user adding itself")
	ErrUsernameInUse             = errors.New("username in use")
	ErrEmailInUse                = errors.New("email in use")
	ErrUnauthorizedEmojiDeletion = errors.New("user cannot delete this emoji")
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

type UploadEmojiBody struct {
	Shortcodes []string `validate:"required,max=20,dive,emoji_shortcode" json:"shortcode"`
}

type UploadEmojiResponse struct {
	Id        string `json:"id"`
	Url       string `json:"url"`
	Shortcode string `json:"shortcode"`
}

type UpdateEmojiBody struct {
	Shortcode string `validate:"emoji_shortcode" json:"shortcode"`
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

		_, err = db.Query.UpdateUserUsername(ctx, db.UpdateUserUsernameParams{
			ID:       user.ID,
			Username: body.Username,
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

func UpdateProfile(ctx context.Context, body *UpdateProfileBody) ([]byte, []byte, error) {
	user := ctx.Value("user").(db.User)

	if body.DisplayName != "" {
		err := db.Query.UpdateUserDisplayName(ctx, db.UpdateUserDisplayNameParams{
			ID:          user.ID,
			DisplayName: body.DisplayName,
		})
		if err != nil {
			return nil, nil, err
		}
	}

	if len(body.About) > 0 {
		err := db.Query.UpdateUserAbout(ctx, db.UpdateUserAboutParams{
			ID:    user.ID,
			About: body.About,
		})
		if err != nil {
			return nil, nil, err
		}
	}

	links := make([]Link, 0)
	for i, link := range body.Links {
		if i >= 2 {
			break
		}
		link.Id = utils.Node.Generate().String()
		links = append(links, link)
	}
	jsonLinks, err := json.Marshal(links)
	if err != nil {
		return nil, nil, err
	}

	err = db.Query.UpdateUserLinks(ctx, db.UpdateUserLinksParams{
		ID:    user.ID,
		Links: jsonLinks,
	})
	if err != nil {
		return nil, nil, err
	}

	facts := make([]Fact, 0)
	for i, fact := range body.Facts {
		if i >= 3 {
			break
		}
		fact.Id = utils.Node.Generate().String()
		facts = append(facts, fact)
	}
	jsonFacts, err := json.Marshal(facts)
	if err != nil {
		return nil, nil, err
	}

	err = db.Query.UpdateUserFacts(ctx, db.UpdateUserFactsParams{
		ID:    user.ID,
		Facts: jsonFacts,
	})
	if err != nil {
		return nil, nil, err
	}

	return jsonFacts, jsonLinks, nil
}

func UpdateAvatar(ctx context.Context, file []byte, fileHeader *multipart.FileHeader, body *UpdateAvatarBody) (*UpdateAvatarResponse, error) {
	user := ctx.Value("user").(db.User)
	s3Client := s3.NewFromConfig(GetAWSConfig())

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
	oldAvatarSplit := strings.Split(user.Avatar.String, "/")
	oldBannerSplit := strings.Split(user.Banner.String, "/")
	defaultAvatars := []string{"avatar_1.webp", "avatar_2.webp", "avatar_3.webp", "avatar_4.webp"}

	// upload new avatars
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Key:    &avatarFileName,
		Bucket: aws.String("nyo-files"),
		Body:   bytes.NewReader(avatar),
	})
	if err != nil {
		slog.Error("failed uploading user avatar", "err", err)
		return nil, err
	}

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Key:    &bannerFileName,
		Bucket: aws.String("nyo-files"),
		Body:   bytes.NewReader(banner),
	})
	if err != nil {
		slog.Error("failed uploading user banner", "err", err)
		return nil, err
	}

	if !slices.Contains(defaultAvatars, oldAvatarSplit[len(oldAvatarSplit)-1]) {
		// delete old avatar
		_, err = s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Key:    aws.String(oldAvatarSplit[len(oldAvatarSplit)-1]),
			Bucket: aws.String("nyo-files"),
		})
		if err != nil {
			slog.Error("failed deleting user avatar", "err", err)
			return nil, err
		}
	}

	if !slices.Contains(defaultAvatars, oldBannerSplit[len(oldBannerSplit)-1]) {
		// delete old banner
		_, err = s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Key:    aws.String(oldBannerSplit[len(oldBannerSplit)-1]),
			Bucket: aws.String("nyo-files"),
		})
		if err != nil {
			slog.Error("failed deleting user banner", "err", err)
			return nil, err
		}
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

func UploadEmojis(ctx context.Context, files []*multipart.FileHeader, body *UploadEmojiBody) ([]UploadEmojiResponse, error) {
	user := ctx.Value("user").(db.User)
	s3Client := s3.NewFromConfig(GetAWSConfig())

	var emojiData []db.CreateEmojiParams
	var responses []UploadEmojiResponse

	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		emojiImg, err := utils.ConvertToEmoji(file)
		if err != nil {
			slog.Error("failed converting emoji to webp", "err", err)
			return nil, err
		}

		randomId := utils.GenerateRandomId(8)
		emojiFileName := fmt.Sprintf("emoji-%s-%s.webp", user.ID, randomId)

		_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Key:    &emojiFileName,
			Bucket: aws.String("nyo-files"),
			Body:   bytes.NewReader(emojiImg),
		})
		if err != nil {
			slog.Error("failed uploading emoji", "err", err)
			return nil, err
		}

		emojiUrl := fmt.Sprintf("%s/%s", os.Getenv("CDN_URL"), emojiFileName)
		emojiID := utils.Node.Generate().String()

		emojiData = append(emojiData, db.CreateEmojiParams{
			ID:        emojiID,
			UserID:    user.ID,
			Url:       emojiUrl,
			Shortcode: body.Shortcodes[i],
		})

		responses = append(responses, UploadEmojiResponse{
			Id:        emojiID,
			Url:       emojiUrl,
			Shortcode: body.Shortcodes[i],
		})
	}

	if len(emojiData) > 0 {
		if _, err := db.Query.CreateEmoji(ctx, emojiData); err != nil {
			return nil, fmt.Errorf("failed to batch insert emojis: %w", err)
		}
	}

	return responses, nil
}

func UpdateEmoji(ctx context.Context, emojiId string, body *UpdateEmojiBody) error {
	user := ctx.Value("user").(db.User)

	err := db.Query.UpdateEmoji(ctx, db.UpdateEmojiParams{
		ID:        emojiId,
		UserID:    user.ID,
		Shortcode: body.Shortcode,
	})
	if err != nil {
		return ErrUnauthorizedEmojiDeletion
	}

	return nil
}

func DeleteEmoji(ctx context.Context, emojiId string) error {
	user := ctx.Value("user").(db.User)

	err := db.Query.DeleteEmoji(ctx, db.DeleteEmojiParams{
		ID:     emojiId,
		UserID: user.ID,
	})
	if err != nil {
		return ErrUnauthorizedEmojiDeletion
	}

	return nil
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
	user := ctx.Value("user").(db.User)
	err := db.Query.AcceptFriend(ctx, db.AcceptFriendParams{
		ID:     body.FriendshipID,
		UserID: user.ID,
	})
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
		Column2: body.UserID,
	})
	if err != nil {
		return "", err
	}

	return channel.ID, nil
}
