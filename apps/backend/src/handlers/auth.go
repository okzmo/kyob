package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/okzmo/kyob/db"
	"github.com/okzmo/kyob/src/services"
	"github.com/okzmo/kyob/src/utils"
)

type UserResponse struct {
	ID          int64       `json:"id"`
	Email       string      `json:"email"`
	Username    string      `json:"username"`
	DisplayName string      `json:"display_name"`
	Avatar      pgtype.Text `json:"avatar"`
	About       pgtype.Text `json:"about"`
}

type signInParams struct {
	EmailOrUsername string `validate:"required" json:"email_or_username"`
	Password        string `validate:"required" json:"password"`
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var body signInParams

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

	user, token, err := services.SignIn(r.Context(), body.EmailOrUsername, body.Password)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserNotFound):
			utils.RespondWithError(w, http.StatusNotFound, "No user exist under this email or username.")
		case errors.Is(err, services.ErrInvalidHash):
			utils.RespondWithError(w, http.StatusInternalServerError, "The informations are incorrect.")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	res := UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Avatar:      user.Avatar,
		About:       user.About,
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    *token,
		Path:     "/",
		Expires:  time.Now().Add(30 * (24 * time.Hour)),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	utils.RespondWithJSON(w, http.StatusContinue, res)
}

type signUpParams struct {
	Email       string `validate:"required,email" json:"email"`
	Username    string `validate:"required,max=50" json:"username"`
	DisplayName string `validate:"required,max=50" json:"display_name"`
	Password    string `validate:"required,min=8,max=254" json:"password"`
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var body signUpParams

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

	token, err := services.SignUp(r.Context(), body.Email, body.Username, body.DisplayName, body.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    *token,
		Path:     "/",
		MaxAge:   time.Now().Add(30 * (24 * time.Hour)).Second(),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	utils.RespondWithJSON(w, http.StatusCreated, &DefaultResponse{Message: "success"})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(db.User)
	err := db.Query.DeleteRememberMeToken(r.Context(), user.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	utils.RespondWithJSON(w, http.StatusContinue, &DefaultResponse{Message: "success"})
}
