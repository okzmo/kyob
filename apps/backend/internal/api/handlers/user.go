package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/okzmo/kyob/db"
	"github.com/okzmo/kyob/internal/api/actors"
	services "github.com/okzmo/kyob/internal/service"
	"github.com/okzmo/kyob/internal/utils"
	proto "github.com/okzmo/kyob/types"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "user_id")

	user, err := services.GetUser(r.Context(), userId)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserNotFound):
			utils.RespondWithError(w, http.StatusNotFound, "No user found.")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusContinue, user)
}

func AddFriend(w http.ResponseWriter, r *http.Request) {
	var body services.AddFriendBody
	user := r.Context().Value("user").(db.User)

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = validate.Struct(body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	inviteId, friendId, err := services.AddFriend(r.Context(), &body)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserNotFound):
			utils.RespondWithError(w, http.StatusNotFound, "User not found.")
		case errors.Is(err, services.ErrAddingItself):
			utils.RespondWithError(w, http.StatusForbidden, "You can't add yourself.")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	inviteMessage := &proto.SendFriendInvite{
		InviteId: inviteId,
		User: &proto.User{
			Id:          user.ID,
			DisplayName: user.DisplayName,
			Avatar:      &user.Avatar.String,
			About:       &user.About.String,
		},
	}

	friendPid := actors.UsersEngine.Registry.GetPID("user", friendId)
	actors.UsersEngine.Send(friendPid, inviteMessage)

	utils.RespondWithJSON(w, http.StatusContinue, DefaultResponse{Message: "success"})
}

func AcceptFriend(w http.ResponseWriter, r *http.Request) {
	var body services.AcceptFriendBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = validate.Struct(body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	friend, existingChannel, err := services.AcceptFriend(r.Context(), &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var channelId string

	friendPid := actors.UsersEngine.Registry.GetPID("user", body.FriendID)
	userPid := actors.UsersEngine.Registry.GetPID("user", body.UserID)
	globalServerPid := actors.ServersEngine.Registry.GetPID("server", "global")

	if existingChannel != nil {
		channelId = existingChannel.ID
		actors.ServersEngine.Send(globalServerPid, &proto.StartChannel{
			ChannelId: channelId,
			Users:     []string{body.FriendID, body.UserID},
		})
	} else {
		channelId = utils.Node.Generate().String()
		newChannelMessage := &proto.BodyChannelCreation{
			ServerId:  "global",
			CreatorId: "global",
			Name:      "friends",
			Type:      "dm",
			Users:     []string{body.FriendID, body.UserID},
			X:         0,
			Y:         0,
			Id:        channelId,
		}

		actors.ServersEngine.Send(globalServerPid, newChannelMessage)
	}

	friendMessage := &proto.AcceptFriendInvite{
		InviteId:  body.FriendshipID,
		ChannelId: channelId,
		User: &proto.User{
			Id:          friend.ID,
			DisplayName: friend.DisplayName,
			Avatar:      &friend.Avatar.String,
			About:       &friend.About.String,
		},
		Sender: true,
	}
	actors.UsersEngine.Send(friendPid, friendMessage)

	receiverMessage := &proto.AcceptFriendInvite{
		InviteId:  body.FriendshipID,
		ChannelId: channelId,
		Sender:    false,
	}
	actors.UsersEngine.Send(userPid, receiverMessage)

	utils.RespondWithJSON(w, http.StatusContinue, DefaultResponse{Message: "success"})
}

func DeleteFriend(w http.ResponseWriter, r *http.Request) {
	var body services.RemoveFriendBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = validate.Struct(body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	channelId, err := services.DeleteFriend(r.Context(), &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	deleteFriendMessage := &proto.DeleteFriend{
		InviteId: body.FriendshipID,
		UserId:   body.FriendID,
	}

	friendPid := actors.UsersEngine.Registry.GetPID("user", body.FriendID)
	globalServerPid := actors.ServersEngine.Registry.GetPID("server", "global")

	actors.UsersEngine.Send(friendPid, deleteFriendMessage)
	actors.ServersEngine.Send(globalServerPid, &proto.KillChannel{
		ChannelId: channelId,
		Users:     []string{body.FriendID, body.UserID},
	})

	utils.RespondWithJSON(w, http.StatusContinue, DefaultResponse{Message: "success"})
}
