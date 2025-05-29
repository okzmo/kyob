package services

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/okzmo/kyob/db"
)

type ServerWithChannels struct {
	db.Server
	X int32 `json:"x"`
	Y int32 `json:"y"`
	// Roles    []db.Role         `json:"roles"`
	Channels    map[string]ChannelsWithMembers `json:"channels"`
	MemberCount int                            `json:"member_count"`
	Members     []db.GetServerMembersRow       `json:"members"`
}

type ChannelsWithMembers struct {
	db.Channel
	Users []db.GetServerMembersRow `json:"users"`
}

type UserResponse struct {
	ID             string               `json:"id"`
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

type FriendResponse struct {
	ID           string      `json:"id"`
	FriendshipID string      `json:"friendship_id"`
	ChannelID    string      `json:"channel_id"`
	DisplayName  string      `json:"display_name"`
	Avatar       pgtype.Text `json:"avatar"`
	About        pgtype.Text `json:"about"`
	Accepted     bool        `json:"accepted"`
	Sender       bool        `json:"sender"`
}

type SetupResponse struct {
	User    UserResponse                  `json:"user"`
	Friends []FriendResponse              `json:"friends"`
	Servers map[string]ServerWithChannels `json:"servers"`
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

	friends, err := db.Query.GetFriends(ctx, ctxUser.ID)
	if err != nil {
		return nil, err
	}

	for _, f := range friends {
		res.Friends = append(res.Friends, FriendResponse{
			ID:           f.ID,
			FriendshipID: f.FriendshipID,
			ChannelID:    f.ChannelID,
			DisplayName:  f.DisplayName,
			Avatar:       f.Avatar,
			About:        f.About,
			Accepted:     f.Accepted,
			Sender:       f.FriendshipSenderID == ctxUser.ID,
		})
	}

	res.Servers = make(map[string]ServerWithChannels)
	servers, err := db.Query.GetServersFromUser(ctx, ctxUser.ID)
	if err != nil {
		return nil, err
	}

	for _, server := range servers {
		channelMap := make(map[string]ChannelsWithMembers)
		channels, err := db.Query.GetChannelsFromServer(ctx, server.ID)
		if err != nil {
			return nil, err
		}

		for _, channelRaw := range channels {
			channel := ChannelsWithMembers{
				channelRaw,
				[]db.GetServerMembersRow{},
			}

			if len(channelRaw.Users) > 0 {
				for _, userId := range channelRaw.Users {
					user, err := db.Query.GetUserById(ctx, userId)
					channel.Users = append(channel.Users, db.GetServerMembersRow{
						ID:          user.ID,
						Username:    user.Username,
						DisplayName: user.DisplayName,
						Avatar:      user.Avatar,
					})
					if err != nil {
						return nil, err
					}
				}
			}

			channelMap[channel.ID] = channel
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
			server.X.Int32,
			server.Y.Int32,
			channelMap,
			int(server.MemberCount),
			users,
		}

		res.Servers[server.ID] = s
	}

	return &res, nil
}
