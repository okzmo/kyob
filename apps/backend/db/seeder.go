package db

import (
	"context"
	"log"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
	queries "github.com/okzmo/kyob/db/gen_queries"
	"github.com/okzmo/kyob/internal/utils"
)

func RunSeeder() {
	createUsers()
}

func createUsers() {
	var servers []queries.CreateServerParams

	for range 100 {
		servers = append(servers, queries.CreateServerParams{
			ID:        utils.Node.Generate().String(),
			OwnerID:   "1936567174520901632",
			Name:      gofakeit.Word(),
			Avatar:    pgtype.Text{String: "https://d2vxg81co3irvv.cloudfront.net/avatar-server-b0sK74ea.webp", Valid: true},
			MainColor: pgtype.Text{String: "12,12,15", Valid: true},
		})
	}

	for _, server := range servers {
		Query.CreateServer(context.TODO(), server)

		err := Query.JoinServer(context.TODO(), queries.JoinServerParams{
			ID:       utils.Node.Generate().String(),
			ServerID: server.ID,
			UserID:   "1936567174520901632",
			X:        int32(gofakeit.IntRange(0, 3000)),
			Y:        int32(gofakeit.IntRange(0, 3000)),
		})
		if err != nil {
			log.Printf("Failed to join server %s: %v", server.ID, err)
		}

		for range 20 {
			Query.CreateChannel(context.TODO(), queries.CreateChannelParams{
				ID:       utils.Node.Generate().String(),
				ServerID: server.ID,
				Name:     gofakeit.Word(),
				Type:     queries.ChannelType("textual"),
				X:        int32(gofakeit.IntRange(0, 3000)),
				Y:        int32(gofakeit.IntRange(0, 3000)),
			})
		}

		for range 100 {
			user, _ := Query.CreateUser(context.TODO(), queries.CreateUserParams{
				ID:          utils.Node.Generate().String(),
				Email:       gofakeit.Email(),
				Username:    gofakeit.Username(),
				DisplayName: gofakeit.Word(),
				MainColor:   pgtype.Text{String: "12,12,15", Valid: true},
				Password:    gofakeit.Password(true, true, true, true, false, 32),
				Avatar:      pgtype.Text{String: "https://d2vxg81co3irvv.cloudfront.net/avatar-1936492940591370240-Xw8gdSeX.webp", Valid: true},
				Banner:      pgtype.Text{String: "https://d2vxg81co3irvv.cloudfront.net/banner-1936492940591370240-Xw8gdSeX.webp", Valid: true},
			})

			Query.JoinServer(context.TODO(), queries.JoinServerParams{
				ID:       utils.Node.Generate().String(),
				ServerID: server.ID,
				UserID:   user.ID,
				X:        int32(gofakeit.IntRange(0, 3000)),
				Y:        int32(gofakeit.IntRange(0, 3000)),
			})
		}
	}
}
