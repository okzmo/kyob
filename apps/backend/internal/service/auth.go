package services

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/okzmo/kyob/db"
	queries "github.com/okzmo/kyob/db/gen_queries"
	"github.com/okzmo/kyob/internal/utils"
)

var (
	ErrInvalidHash  = errors.New("invalid hash")
	ErrUserNotFound = errors.New("user not found")
)

func SignIn(ctx context.Context, emailOrUsername string, password string) (*string, error) {
	if emailOrUsername == "admin" {
		return nil, ErrInvalidHash
	}

	user, err := db.Query.GetUser(ctx, queries.GetUserParams{
		Email:    emailOrUsername,
		Username: emailOrUsername,
	})
	if err != nil {
		return nil, ErrUserNotFound
	}

	match, err := utils.VerifyPassword(password, user.Password)
	if err != nil {
		slog.Error("error on hashing", "err", err)
		return nil, ErrInvalidHash
	} else if !match {
		return nil, ErrInvalidHash
	}

	token, err := utils.GenerateRandomBytes(64)
	if err != nil {
		slog.Error("failed generate token", "err", err)
		return nil, err
	}

	b64Token := base64.RawStdEncoding.EncodeToString(token)
	_, err = db.Query.CreateToken(ctx, queries.CreateTokenParams{
		ID:       utils.Node.Generate().String(),
		UserID:   user.ID,
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

func SignUp(ctx context.Context, email, username, displayName, password, bodyURL, rpmID, rpmToken string) (*string, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	avatarFileName := fmt.Sprintf("avatar_%d.webp", rand.Intn(4)+1)
	avatarURL := pgtype.Text{String: fmt.Sprintf("%s/%s", os.Getenv("CDN_URL"), avatarFileName), Valid: true}
	mainColor := pgtype.Text{String: "12,12,16", Valid: true}
	queriesUser, err := db.Query.CreateUser(ctx, queries.CreateUserParams{
		ID:          utils.Node.Generate().String(),
		Email:       email,
		DisplayName: displayName,
		Username:    username,
		Password:    hashedPassword,
		Avatar:      avatarURL,
		Banner:      avatarURL,
		MainColor:   mainColor,
		Body:        pgtype.Text{String: bodyURL, Valid: true},
		RpmID:       pgtype.Text{String: rpmID, Valid: true},
		RpmToken:    pgtype.Text{String: rpmToken, Valid: true},
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
	_, err = db.Query.CreateToken(ctx, queries.CreateTokenParams{
		ID:       utils.Node.Generate().String(),
		UserID:   queriesUser.ID,
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
