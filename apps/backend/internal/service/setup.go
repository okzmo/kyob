package services

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/okzmo/kyob/db"
)

type serverWithChannels struct {
	db.Server
	X int `json:"x"`
	Y int `json:"y"`
	// Roles    []db.Role         `json:"roles"`
	IsMember bool                 `json:"is_member"`
	Channels map[int64]db.Channel `json:"channels"`
}

type UserResponse struct {
	ID             int64                `json:"id"`
	Email          string               `json:"email"`
	Username       string               `json:"username"`
	DisplayName    string               `json:"display_name"`
	Avatar         pgtype.Text          `json:"avatar"`
	Banner         pgtype.Text          `json:"banner"`
	GradientTop    pgtype.Text          `json:"gradient_top"`
	GradientBottom pgtype.Text          `json:"gradient_bottom"`
	About          pgtype.Text          `json:"about"`
	CreatedAt      time.Time            `json:"created_at"`
	Links          []db.GetUserLinksRow `json:"links"`
	Facts          []db.GetUserFactsRow `json:"facts"`
}

type SetupResponse struct {
	User    UserResponse                 `json:"user"`
	Servers map[int64]serverWithChannels `json:"servers"`
}

func GetSetup(ctx context.Context) (*SetupResponse, error) {
	var res SetupResponse

	ctxUser := ctx.Value("user").(db.User)

	facts, err := db.Query.GetUserFacts(ctx, ctxUser.ID)
	if err != nil {
		return nil, err
	}

	links, err := db.Query.GetUserLinks(ctx, ctxUser.ID)
	if err != nil {
		return nil, err
	}

	res.User = UserResponse{
		ID:             ctxUser.ID,
		Email:          ctxUser.Email,
		Username:       ctxUser.Username,
		DisplayName:    ctxUser.DisplayName,
		Avatar:         ctxUser.Avatar,
		Banner:         ctxUser.Banner,
		GradientTop:    ctxUser.GradientTop,
		GradientBottom: ctxUser.GradientBottom,
		About:          ctxUser.About,
		CreatedAt:      ctxUser.CreatedAt,
		Facts:          facts,
		Links:          links,
	}
	res.Servers = make(map[int64]serverWithChannels)

	servers, err := db.Query.GetServersFromUser(ctx, ctxUser.ID)
	if err != nil {
		return nil, err
	}

	for _, server := range servers {
		memberRes, err := db.Query.IsMember(ctx, db.IsMemberParams{
			UserID:   ctxUser.ID,
			ServerID: server.ID,
		})
		if err != nil {
			return nil, err
		}

		channelMap := make(map[int64]db.Channel)
		channels, err := db.Query.GetChannelsFromServer(ctx, server.ID)
		if err != nil {
			return nil, err
		}

		for _, channel := range channels {
			channelMap[channel.ID] = channel
		}

		s := serverWithChannels{
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
			memberRes.RowsAffected() > 0,
			channelMap,
		}

		res.Servers[server.ID] = s
	}

	return &res, nil
}
