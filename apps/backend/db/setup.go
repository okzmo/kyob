package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DbManager struct {
	conn *pgxpool.Pool
}

var Query *Queries

func Setup() *DbManager {
	dsn := os.Getenv("DATABASE_URL")
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	Query = New(conn)
	return &DbManager{conn: conn}
}

func (db *DbManager) Close() {
	db.conn.Close()
}
