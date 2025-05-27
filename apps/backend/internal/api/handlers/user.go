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

	err = services.AcceptFriend(r.Context(), &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// inviteMessage := &proto.SendFriendInvite{
	// 	InviteId: inviteId,
	// 	User: &proto.User{
	// 		Id:          user.ID,
	// 		DisplayName: user.DisplayName,
	// 		Avatar:      &user.Avatar.String,
	// 	},
	// }
	//
	// friendPid := actors.UsersEngine.Registry.GetPID("user", friendId)
	// actors.UsersEngine.Send(friendPid, inviteMessage)

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

	err = services.DeleteFriend(r.Context(), &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// inviteMessage := &proto.SendFriendInvite{
	// 	InviteId: inviteId,
	// 	User: &proto.User{
	// 		Id:          user.ID,
	// 		DisplayName: user.DisplayName,
	// 		Avatar:      &user.Avatar.String,
	// 	},
	// }
	//
	// friendPid := actors.UsersEngine.Registry.GetPID("user", friendId)
	// actors.UsersEngine.Send(friendPid, inviteMessage)

	utils.RespondWithJSON(w, http.StatusContinue, DefaultResponse{Message: "success"})
}
