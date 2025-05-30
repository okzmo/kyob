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
	Members     []db.GetMembersFromServersRow  `json:"members"`
}

type ChannelsWithMembers struct {
	db.Channel
	Users []db.GetUsersByIdsRow `json:"users"`
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

	friends, err := db.Query.GetFriends(ctx, ctxUser.ID)
	if err != nil {
		return nil, err
	}

	servers, err := db.Query.GetServersFromUser(ctx, ctxUser.ID)
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

	for _, f := range friends {
		res.Friends = append(res.Friends, FriendResponse{
			ID:           f.ID,
			FriendshipID: f.FriendshipID,
			ChannelID:    f.ChannelID.String,
			DisplayName:  f.DisplayName,
			Avatar:       f.Avatar,
			About:        f.About,
			Accepted:     f.Accepted,
			Sender:       f.FriendshipSenderID == ctxUser.ID,
		})
	}

	res.Servers = make(map[string]ServerWithChannels)
	if len(servers) > 0 {
		serversMap, err := processServers(ctx, servers)
		if err != nil {
			return nil, err
		}

		res.Servers = serversMap
	}

	return &res, nil
}

func processServers(ctx context.Context, servers []db.GetServersFromUserRow) (map[string]ServerWithChannels, error) {
	serverIDs := make([]string, 0, len(servers))
	for _, server := range servers {
		serverIDs = append(serverIDs, server.ID)
	}

	allChannels, err := db.Query.GetChannelsFromServers(ctx, serverIDs)
	if err != nil {
		return nil, err
	}

	allMembers, err := db.Query.GetMembersFromServers(ctx, serverIDs)
	if err != nil {
		return nil, err
	}

	userIDSet := make(map[string]bool)
	for _, channel := range allChannels {
		for _, userID := range channel.Users {
			userIDSet[userID] = true
		}
	}

	userIDs := make([]string, 0, len(userIDSet))
	for userId := range userIDSet {
		userIDs = append(userIDs, userId)
	}

	var allUsers []db.GetUsersByIdsRow
	if len(userIDs) > 0 {
		allUsers, err = db.Query.GetUsersByIds(ctx, userIDs)
		if err != nil {
			return nil, err
		}
	}

	userMap := make(map[string]db.GetUsersByIdsRow)
	for _, user := range allUsers {
		userMap[user.ID] = user
	}

	channelsByServer := make(map[string][]db.Channel)
	for _, channel := range allChannels {
		channelsByServer[channel.ServerID] = append(channelsByServer[channel.ServerID], channel)
	}

	membersByServer := make(map[string][]db.GetMembersFromServersRow)
	for _, member := range allMembers {
		membersByServer[member.ServerID] = append(membersByServer[member.ServerID], member)
	}

	result := make(map[string]ServerWithChannels)
	for _, server := range servers {
		channelMap := make(map[string]ChannelsWithMembers)

		for _, channel := range channelsByServer[server.ID] {
			channelUsers := make([]db.GetUsersByIdsRow, 0, len(channel.Users))
			for _, userID := range channel.Users {
				if user, exists := userMap[userID]; exists {
					channelUsers = append(channelUsers, user)
				}
			}

			channelMap[channel.ID] = ChannelsWithMembers{
				channel,
				channelUsers,
			}
		}

		result[server.ID] = ServerWithChannels{
			Server: db.Server{
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
			X:           server.X.Int32,
			Y:           server.Y.Int32,
			Channels:    channelMap,
			MemberCount: int(server.MemberCount),
			Members:     membersByServer[server.ID],
		}
	}

	return result, nil
}
