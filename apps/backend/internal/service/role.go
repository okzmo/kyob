package services

import (
	"context"

	"github.com/okzmo/kyob/db"
	"github.com/okzmo/kyob/internal/utils"
)

type BodyRoleCreation struct {
	Name      string   `validate:"required,max=20" json:"name"`
	Color     string   `validate:"required" json:"color"`
	Abilities []string `json:"abilities"`
	Index     int      `json:"index"`
}

type BodyMoveRole struct {
	RoleID string `json:"role_id"`
	From   int    `json:"from"`
	To     int    `json:"to"`
}

func CreateRole(ctx context.Context, serverID string, body *BodyRoleCreation) (*db.Role, error) {
	role, err := db.Query.CreateRole(ctx, db.CreateRoleParams{
		ID:        utils.Node.Generate().String(),
		ServerID:  serverID,
		Name:      body.Name,
		Color:     body.Color,
		Abilities: body.Abilities,
		Idx:       int32(body.Index),
	})
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func GetRoles(ctx context.Context, serverID string) ([]db.Role, error) {
	roles, err := db.Query.GetRoles(ctx, serverID)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func MoveRole(ctx context.Context, body *BodyMoveRole) error {
	err := db.Query.UpdateRolePositions(ctx, db.UpdateRolePositionsParams{
		Idx:   int32(body.To),
		Idx_2: int32(body.From),
	})
	if err != nil {
		return err
	}

	err = db.Query.MoveRole(ctx, db.MoveRoleParams{
		ID:  body.RoleID,
		Idx: int32(body.To),
	})
	if err != nil {
		return err
	}

	return nil
}

func DeleteRole(ctx context.Context, roleID string) error {
	err := db.Query.DeleteRole(ctx, roleID)
	if err != nil {
		return err
	}

	return nil
}
