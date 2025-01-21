package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	conn *pgxpool.Pool
}

func NewDB() *DB {
	connString := os.Getenv("GOOSE_DBSTRING")
	if connString == "" {
		log.Fatal("GOOSE_DBSTRING environment variable is not set")
	}
	conn, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		panic(err)
	}

	query := `SELECT 1 FROM pg_database WHERE datname = $1;`
	var exists int
	err = conn.QueryRow(context.Background(), query, "simpleurl").Scan(&exists)
	log.Printf("Database simpleurl exists: %v", exists)
	if err != nil && err != pgx.ErrNoRows {
		panic(err)
	}

	if exists == 0 {
		query = `CREATE DATABASE simpleurl;`
		_, err = conn.Exec(context.Background(), query)
		if err != nil {
			panic(err)
		}
		log.Printf("Database simpleurl created")
	}

	return &DB{
		conn: conn,
	}
}
