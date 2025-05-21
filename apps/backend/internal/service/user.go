package services

import (
	"context"

	"github.com/okzmo/kyob/db"
)

func GetUser(ctx context.Context, userId string) (*UserResponse, error) {
	user, err := db.Query.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
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
