package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type DbManager struct {
	conn *pgx.Conn
}

var Query *Queries

func Setup() *DbManager {
	dsn := os.Getenv("DATABASE_URL")
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	Query = New(conn)
	return &DbManager{conn: conn}
}

func (db *DbManager) Close() error {
	return db.conn.Close(context.Background())
}
