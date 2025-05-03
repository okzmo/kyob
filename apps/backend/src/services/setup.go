package services

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/okzmo/kyob/db"
)

type serverWithChannels struct {
	db.Server
	Channels []db.Channel `json:"channels"`
}

type UserResponse struct {
	ID          int64       `json:"id"`
	Email       string      `json:"email"`
	Username    string      `json:"username"`
	DisplayName string      `json:"display_name"`
	Avatar      pgtype.Text `json:"avatar"`
	About       pgtype.Text `json:"about"`
	CreatedAt   time.Time   `json:"created_at"`
}

type SetupResponse struct {
	User    UserResponse         `json:"user"`
	Servers []serverWithChannels `json:"servers"`
}

func GetSetup(ctx context.Context) (*SetupResponse, error) {
	var res SetupResponse

	ctxUser := ctx.Value("user").(db.User)
	res.User = UserResponse{
		ID:          ctxUser.ID,
		Email:       ctxUser.Email,
		Username:    ctxUser.Username,
		DisplayName: ctxUser.DisplayName,
		Avatar:      ctxUser.Avatar,
		About:       ctxUser.About,
		CreatedAt:   ctxUser.CreatedAt,
	}

	servers, err := db.Query.GetServersFromUser(ctx, ctxUser.ID)
	if err != nil {
		return nil, err
	}

	for _, server := range servers {
		channels, err := db.Query.GetChannelsFromServer(ctx, db.GetChannelsFromServerParams{
			UserID:   ctxUser.ID,
			ServerID: server.ID,
		})
		if err != nil {
			return nil, err
		}

		s := serverWithChannels{
			server,
			[]db.Channel(channels),
		}

		res.Servers = append(res.Servers, s)
	}

	return &res, nil
}
