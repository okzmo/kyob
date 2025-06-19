package handlers

import (
	"encoding/json"
	"errors"
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
	userID := chi.URLParam(r, "user_id")

	user, err := services.GetUser(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserNotFound):
			utils.RespondWithError(w, http.StatusNotFound, "No user found.")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, user)
}

func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	var body services.UpdateAccountBody

	err := utils.ParseAndValidate(r, validate, &body)
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

	if body.Username != "" {
		user := r.Context().Value("user").(db.User)
		userPID := actors.UsersEngine.Registry.GetPID("user", user.ID)
		actors.UsersEngine.Send(userPID, &proto.UserChangedInformations{
			UserId: user.ID,
			UserInformations: &proto.UserInformations{
				Username: &body.Username,
			},
		})
	}

	utils.RespondWithJSON(w, http.StatusOK, DefaultResponse{Message: "success"})
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var body services.UpdateProfileBody

	err := utils.ParseAndValidate(r, validate, &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	facts, links, err := services.UpdateProfile(r.Context(), &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user := r.Context().Value("user").(db.User)
	userPID := actors.UsersEngine.Registry.GetPID("user", user.ID)
	messageToSend := &proto.UserChangedInformations{
		UserId: user.ID,
		UserInformations: &proto.UserInformations{
			DisplayName: &body.DisplayName,
			About:       body.About,
		},
	}

	if len(links) > 0 {
		messageToSend.UserInformations.Links = links
	}

	if len(facts) > 0 {
		messageToSend.UserInformations.Facts = facts
	}

	actors.UsersEngine.Send(userPID, messageToSend)

	utils.RespondWithJSON(w, http.StatusOK, DefaultResponse{Message: "success"})
}

func UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	var body services.UpdateAvatarBody
	var cropAvatar, cropBanner services.Crop
	config := utils.ImageValidationConfig{
		MaxSize: 10 << 20, // 10 MB
		AllowedMimeTypes: []string{
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/webp",
		},
		RequireValidHeader: true,
	}

	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		slog.Error(err.Error())
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to get image.")
		return
	}
	defer file.Close()

	if err := utils.ParseAndValidateImage(fileHeader, config); err != nil {
		slog.Error("Avatar image validation failed", "error", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Avatar's invalid.")
		return
	}

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

	user := r.Context().Value("user").(db.User)
	userPID := actors.UsersEngine.Registry.GetPID("user", user.ID)
	actors.UsersEngine.Send(userPID, &proto.UserChangedInformations{
		UserId: user.ID,
		UserInformations: &proto.UserInformations{
			Avatar:    &res.Avatar,
			Banner:    &res.Banner,
			MainColor: &body.MainColor,
		},
	})

	utils.RespondWithJSON(w, http.StatusOK, res)
}

func UploadEmojis(w http.ResponseWriter, r *http.Request) {
	var body services.UploadEmojiBody
	config := utils.ImageValidationConfig{
		MaxSize: 1 << 20, // 1 MB
		AllowedMimeTypes: []string{
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/webp",
		},
		RequireValidHeader: true,
	}

	err := r.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		slog.Error("Failed to parse multipart form", "error", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid form data.", "ERR_INVALID_FORM")
		return
	}

	if r.MultipartForm == nil {
		slog.Error("MultipartForm is nil after parsing")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid form data.")
		return
	}

	emojis, exists := r.MultipartForm.File["emojis[]"]
	if !exists || len(emojis) == 0 {
		slog.Error("No emojis sent")
		utils.RespondWithError(w, http.StatusBadRequest, "No emojis sent.", "ERR_MISSING_EMOJIS")
		return
	}

	for _, emoji := range emojis {
		if err := utils.ParseAndValidateImage(emoji, config); err != nil {
			slog.Error("Emoji image validation failed", "error", err)
			utils.RespondWithError(w, http.StatusBadRequest, "Emoji's invalid.", "ERR_EMOJIS_INVALID")
			return
		}
	}

	shortcodes := r.Form["shortcodes[]"]
	if len(shortcodes) != len(emojis) {
		slog.Error("Missing shortcodes")
		utils.RespondWithError(w, http.StatusBadRequest, "Missing shortcodes", "ERR_MISSING_SHORTCODES")
		return
	}

	body.Shortcodes = shortcodes
	err = validate.Struct(body)
	if err != nil {
		slog.Error("Emoji body validation failed", "error", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Shortcode is invalid.", "ERR_SHORTCODES_INVALID")
		return
	}

	res, err := services.UploadEmojis(r.Context(), emojis, &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, res)
}

func UpdateEmoji(w http.ResponseWriter, r *http.Request) {
	var body services.UpdateEmojiBody
	emojiID := chi.URLParam(r, "emoji_id")

	if err := utils.ParseAndValidate(r, validate, &body); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := services.UpdateEmoji(r.Context(), emojiID, &body)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorizedEmojiDeletion):
			utils.RespondWithError(w, http.StatusForbidden, "You can't update this emoji.")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, DefaultResponse{Message: "ok"})
}

func DeleteEmoji(w http.ResponseWriter, r *http.Request) {
	emojiID := chi.URLParam(r, "emoji_id")

	err := services.DeleteEmoji(r.Context(), emojiID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorizedEmojiDeletion):
			utils.RespondWithError(w, http.StatusForbidden, "You can't delete this emoji.")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, DefaultResponse{Message: "ok"})
}

func AddFriend(w http.ResponseWriter, r *http.Request) {
	var body services.AddFriendBody

	err := utils.ParseAndValidate(r, validate, &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user := r.Context().Value("user").(db.User)
	inviteID, friendID, err := services.AddFriend(r.Context(), &body)
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
		InviteId: inviteID,
		User: &proto.User{
			Id:          user.ID,
			DisplayName: user.DisplayName,
			Avatar:      &user.Avatar.String,
			About:       user.About,
		},
	}

	friendPid := actors.UsersEngine.Registry.GetPID("user", friendID)
	actors.UsersEngine.Send(friendPid, inviteMessage)

	utils.RespondWithJSON(w, http.StatusOK, DefaultResponse{Message: "success"})
}

func AcceptFriend(w http.ResponseWriter, r *http.Request) {
	var body services.AcceptFriendBody

	err := utils.ParseAndValidate(r, validate, &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	friend, existingChannel, err := services.AcceptFriend(r.Context(), &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var channelID string

	friendPid := actors.UsersEngine.Registry.GetPID("user", body.FriendID)
	userPid := actors.UsersEngine.Registry.GetPID("user", body.UserID)
	globalServerPid := actors.ServersEngine.Registry.GetPID("server", "global")

	if existingChannel != nil {
		channelID = existingChannel.ID
		actors.ServersEngine.Send(globalServerPid, &proto.StartChannel{
			ChannelId: channelID,
			Users:     []string{body.FriendID, body.UserID},
		})
	} else {
		channelID = utils.Node.Generate().String()
		newChannelMessage := &proto.BodyChannelCreation{
			ServerId:  "global",
			CreatorId: "global",
			Name:      "friends",
			Type:      "dm",
			Users:     []string{body.FriendID, body.UserID},
			X:         0,
			Y:         0,
			Id:        channelID,
		}

		actors.ServersEngine.Send(globalServerPid, newChannelMessage)
	}

	friendMessage := &proto.AcceptFriendInvite{
		InviteId:  body.FriendshipID,
		ChannelId: channelID,
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
		ChannelId: channelID,
		Sender:    false,
	}
	actors.UsersEngine.Send(userPid, receiverMessage)

	utils.RespondWithJSON(w, http.StatusOK, DefaultResponse{Message: "success"})
}

func DeleteFriend(w http.ResponseWriter, r *http.Request) {
	var body services.RemoveFriendBody

	err := utils.ParseAndValidate(r, validate, &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	channelID, err := services.DeleteFriend(r.Context(), &body)
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
		ChannelId: channelID,
		Users:     []string{body.FriendID, body.UserID},
	})

	utils.RespondWithJSON(w, http.StatusOK, DefaultResponse{Message: "success"})
}
