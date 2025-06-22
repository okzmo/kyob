package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/okzmo/kyob/db/gen_queries"
)

type DBManager struct {
	conn *pgxpool.Pool
}

var Query *db.Queries

func Setup() *DBManager {
	dsn := os.Getenv("DATABASE_URL")
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	Query = db.New(conn)
	return &DBManager{conn: conn}
}

func (db *DBManager) Close() {
	db.conn.Close()
}
