package services

import (
	"context"
	"encoding/base64"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/okzmo/kyob/db"
	"github.com/okzmo/kyob/src/utils"
)

var (
	ErrInvalidHash  = errors.New("invalid hash")
	ErrUserNotFound = errors.New("user not found")
)

func SignIn(ctx context.Context, emailOrUsername string, password string) (*db.User, *string, error) {
	user, err := db.Query.GetUser(ctx, db.GetUserParams{
		Email:    emailOrUsername,
		Username: emailOrUsername,
	})
	if err != nil {
		return nil, nil, ErrUserNotFound
	}

	match, err := utils.VerifyPassword(password, user.Password)
	if err != nil {
		slog.Error("error on hashing", "err", err)
		return nil, nil, ErrInvalidHash
	} else if !match {
		return nil, nil, ErrInvalidHash
	}

	token, err := utils.GenerateRandomBytes(64)
	if err != nil {
		slog.Error("failed generate token", "err", err)
		return nil, nil, err
	}

	b64Token := base64.RawStdEncoding.EncodeToString(token)
	_, err = db.Query.CreateToken(ctx, db.CreateTokenParams{
		UserID:   user.ID,
		Token:    b64Token,
		ExpireAt: time.Now().Add(30 * (24 * time.Hour)),
		Type:     "REMEMBER_ME_TOKEN",
	})
	if err != nil {
		slog.Error("failed create token", "err", err)
		return nil, nil, err
	}

	return &user, &b64Token, nil
}

func SignUp(ctx context.Context, email string, username string, displayName string, password string) (*string, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	dbUser, err := db.Query.CreateUser(ctx, db.CreateUserParams{
		Email:       email,
		DisplayName: displayName,
		Username:    username,
		Password:    hashedPassword,
		Avatar:      pgtype.Text{String: "test", Valid: true},
	})
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateRandomBytes(64)
	if err != nil {
		slog.Error("failed generate token", "err", err)
		return nil, err
	}

	b64Token := base64.RawStdEncoding.EncodeToString(token)
	_, err = db.Query.CreateToken(ctx, db.CreateTokenParams{
		UserID:   dbUser.ID,
		Token:    b64Token,
		ExpireAt: time.Now().Add(30 * (24 * time.Hour)),
		Type:     "REMEMBER_ME_TOKEN",
	})
	if err != nil {
		slog.Error("failed create token", "err", err)
		return nil, err
	}

	return &b64Token, nil
}
