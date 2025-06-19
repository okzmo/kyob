package services

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/okzmo/kyob/db"
)

type channelState struct {
	ChannelID     string `json:"channel_id"`
	LastMessageID string `json:"last_message_id"`
}

type BodySaveState struct {
	UserID         string            `json:"user_id"`
	ChannelIDs     []string          `json:"channel_ids"`
	LastMessageIDs []string          `json:"last_message_ids"`
	MentionsIds    []json.RawMessage `json:"mentions_ids"`
}

type ServerWithChannels struct {
	ServerResponse
	X int32 `json:"x"`
	Y int32 `json:"y"`
	// Roles    []db.Role         `json:"roles"`
	Channels    map[string]ChannelsWithMembers `json:"channels"`
	MemberCount int                            `json:"member_count"`
	Members     []db.GetMembersFromServersRow  `json:"members"`
}

type VoiceUser struct {
	UserId []string `json:"id"`
	Deafen bool     `json:"deafen"`
	Mute   bool     `json:"mute"`
}

type ChannelsWithMembers struct {
	db.Channel
	LastMessageSent string                `json:"last_message_sent"`
	LastMessageRead string                `json:"last_message_read"`
	MentionsIds     json.RawMessage       `json:"last_mentions"`
	Users           []db.GetUsersByIdsRow `json:"users"`
	VoiceUsers      []VoiceUser           `json:"voice_users"`
}

type UserResponse struct {
	ID          string          `json:"id"`
	Email       string          `json:"email"`
	Username    string          `json:"username"`
	DisplayName string          `json:"display_name"`
	Avatar      pgtype.Text     `json:"avatar"`
	Banner      pgtype.Text     `json:"banner"`
	Body        pgtype.Text     `json:"body"`
	MainColor   pgtype.Text     `json:"main_color"`
	About       json.RawMessage `json:"about"`
	Links       json.RawMessage `json:"links"`
	Facts       json.RawMessage `json:"facts"`
	CreatedAt   time.Time       `json:"created_at"`
}

type FriendResponse struct {
	ID           string          `json:"id"`
	FriendshipID string          `json:"friendship_id"`
	ChannelID    string          `json:"channel_id"`
	DisplayName  string          `json:"display_name"`
	Avatar       pgtype.Text     `json:"avatar"`
	Banner       pgtype.Text     `json:"banner"`
	About        json.RawMessage `json:"about"`
	Accepted     bool            `json:"accepted"`
	Sender       bool            `json:"sender"`
}

type SetupResponse struct {
	User    UserResponse                  `json:"user"`
	Emojis  []db.GetEmojisRow             `json:"emojis"`
	Friends []FriendResponse              `json:"friends"`
	Servers map[string]ServerWithChannels `json:"servers"`
}

func GetSetup(ctx context.Context) (*SetupResponse, error) {
	var res SetupResponse

	ctxUser := ctx.Value("user").(db.User)

	friends, err := db.Query.GetFriends(ctx, ctxUser.ID)
	if err != nil {
		return nil, err
	}

	servers, err := db.Query.GetServersFromUser(ctx, ctxUser.ID)
	if err != nil {
		return nil, err
	}

	emojis, err := db.Query.GetEmojis(ctx, ctxUser.ID)
	if err != nil {
		return nil, err
	}

	res.User = UserResponse{
		ID:          ctxUser.ID,
		Email:       ctxUser.Email,
		Username:    ctxUser.Username,
		DisplayName: ctxUser.DisplayName,
		Avatar:      ctxUser.Avatar,
		Banner:      ctxUser.Banner,
		Body:        ctxUser.Body,
		MainColor:   ctxUser.MainColor,
		About:       ctxUser.About,
		CreatedAt:   ctxUser.CreatedAt,
		Links:       ctxUser.Links,
		Facts:       ctxUser.Facts,
	}

	for _, f := range friends {
		res.Friends = append(res.Friends, FriendResponse{
			ID:           f.ID,
			FriendshipID: f.FriendshipID,
			ChannelID:    f.ChannelID.String,
			DisplayName:  f.DisplayName,
			Avatar:       f.Avatar,
			Banner:       f.Banner,
			About:        f.About,
			Accepted:     f.Accepted,
			Sender:       f.FriendshipSenderID == ctxUser.ID,
		})
	}

	res.Emojis = emojis

	res.Servers = make(map[string]ServerWithChannels)
	if len(servers) > 0 {
		serversMap, err := processServers(ctx, ctxUser.ID, servers)
		if err != nil {
			return nil, err
		}

		res.Servers = serversMap
	}

	return &res, nil
}

func processServers(ctx context.Context, userId string, servers []db.GetServersFromUserRow) (map[string]ServerWithChannels, error) {
	serverIDs := make([]string, 0, len(servers))
	for _, server := range servers {
		serverIDs = append(serverIDs, server.ID)
	}

	allChannels, err := db.Query.GetChannelsFromServers(ctx, serverIDs)
	if err != nil {
		return nil, err
	}

	channelIDs := make([]string, 0, len(allChannels))
	for _, channel := range allChannels {
		channelIDs = append(channelIDs, channel.ID)
	}

	allMessagesSent, err := db.Query.GetLatestMessagesSent(ctx, channelIDs)
	if err != nil {
		return nil, err
	}
	allMessagesSentSet := make(map[string]string)
	for _, mess := range allMessagesSent {
		allMessagesSentSet[mess.ChannelID] = mess.ID
	}

	allMessagesRead, err := db.Query.GetLatestMessagesRead(ctx, userId)
	if err != nil {
		return nil, err
	}
	allMessagesReadSet := make(map[string]string)
	allMessagesMentionsSet := make(map[string]json.RawMessage)
	for _, mess := range allMessagesRead {
		allMessagesReadSet[mess.ChannelID] = mess.LastReadMessageID.String
		allMessagesMentionsSet[mess.ChannelID] = mess.UnreadMentionIds
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
				allMessagesSentSet[channel.ID],
				allMessagesReadSet[channel.ID],
				allMessagesMentionsSet[channel.ID],
				channelUsers,
				[]VoiceUser{},
			}
		}

		result[server.ID] = ServerWithChannels{
			ServerResponse{
				ID:          server.ID,
				OwnerID:     server.OwnerID,
				Name:        server.Name,
				Avatar:      server.Avatar,
				Banner:      server.Banner,
				Description: json.RawMessage(server.Description),
				MainColor:   server.MainColor.String,
				Private:     server.Private,
				CreatedAt:   server.CreatedAt,
				UpdatedAt:   server.UpdatedAt,
			},
			server.X.Int32,
			server.Y.Int32,
			channelMap,
			int(server.MemberCount),
			membersByServer[server.ID],
		}
	}

	return result, nil
}

func SaveLastState(ctx context.Context, body BodySaveState) error {
	var validChannelIDs []string
	var validLastReadMessageIDs []string
	var validUnreadMentionIDs []json.RawMessage

	for i, messageID := range body.LastMessageIDs {
		if messageID != "" && messageID != "null" {
			validChannelIDs = append(validChannelIDs, body.ChannelIDs[i])
			validLastReadMessageIDs = append(validLastReadMessageIDs, messageID)
			validUnreadMentionIDs = append(validUnreadMentionIDs, body.MentionsIds[i])
		}
	}

	err := db.Query.SaveUnreadMessagesState(ctx, db.SaveUnreadMessagesStateParams{
		UserID:             body.UserID,
		ChannelIds:         validChannelIDs,
		LastReadMessageIds: validLastReadMessageIDs,
		UnreadMentionIds:   validUnreadMentionIDs,
	})
	if err != nil {
		slog.Error("failed to save last state", "err", err)
		return err
	}

	return nil
}
