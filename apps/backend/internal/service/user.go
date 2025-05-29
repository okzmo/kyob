package services

import (
	"context"
	"errors"

	"github.com/okzmo/kyob/db"
	"github.com/okzmo/kyob/internal/utils"
)

var ErrAddingItself = errors.New("user adding itself")

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
		ID:             user.ID,
		Email:          user.Email,
		Username:       user.Username,
		DisplayName:    user.DisplayName,
		Avatar:         user.Avatar,
		Banner:         user.Banner,
		GradientTop:    user.GradientTop,
		GradientBottom: user.GradientBottom,
		About:          user.About,
		CreatedAt:      user.CreatedAt,
		Links:          links,
		Facts:          facts,
	}

	return res, nil
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
