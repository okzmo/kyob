package services

import (
	"context"

	"github.com/okzmo/kyob/db"
	"github.com/okzmo/kyob/internal/utils"
)

type AddFriendBody struct {
	Username string `validate:"required" json:"username"`
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
