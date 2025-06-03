package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
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

func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	var body services.UpdateAccountBody

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

	err = services.UpdateAccount(r.Context(), &body)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUsernameInUse):
			utils.RespondWithError(w, http.StatusForbidden, "Username already in use.", "ERR_USERNAME_IN_USE")
		case errors.Is(err, services.ErrEmailInUse):
			utils.RespondWithError(w, http.StatusForbidden, "Email already in use.", "ERR_EMAIL_IN_USE")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusContinue, DefaultResponse{Message: "success"})
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var body services.UpdateProfileBody

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

	err = services.UpdateProfile(r.Context(), &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusContinue, DefaultResponse{Message: "success"})
}

func UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	var body services.UpdateAvatarBody
	var cropAvatar, cropBanner services.Crop

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		slog.Error(err.Error())
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to parse given image.")
		return
	}

	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		slog.Error(err.Error())
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to get image.")
		return
	}
	defer file.Close()
	fmt.Println(fileHeader.Header)

	fileData, err := io.ReadAll(file)
	if err != nil {
		slog.Error("Failed to read file", "err", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to process file.")
		return
	}

	cropAvatarJSON := r.FormValue("crop_avatar")
	cropBannerJSON := r.FormValue("crop_banner")
	mainColor := r.FormValue("main_color")

	if err := json.Unmarshal([]byte(cropAvatarJSON), &cropAvatar); err != nil {
		slog.Error(err.Error())
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid crop data.")
		return
	}

	if err := json.Unmarshal([]byte(cropBannerJSON), &cropBanner); err != nil {
		slog.Error(err.Error())
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid crop data.")
		return
	}

	body.CropAvatar = cropAvatar
	body.CropBanner = cropBanner
	body.MainColor = mainColor

	err = validate.Struct(body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := services.UpdateAvatar(r.Context(), fileData, fileHeader, &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(res)

	utils.RespondWithJSON(w, http.StatusContinue, res)
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
			About:       user.About,
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
			About:       friend.About,
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
