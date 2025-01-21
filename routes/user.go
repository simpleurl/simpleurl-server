package routes

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Provider string `json:"provider"`
	Links    []Link `json:"links"`
}

func CreateUser(pdb *pgxpool.Pool, rdb *redis.Client, ctx context.Context, data map[string]string) (*User, error) {
	var user User
	err := pdb.QueryRow(ctx, `INSERT INTO users (username, email, provider) VALUES ($1, $2, $3) RETURNING id, username, email, provider`, data["username"], data["email"], data["provider"]).Scan(&user.ID, &user.Username, &user.Email, &user.Provider)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUser(pdb *pgxpool.Pool, rdb *redis.Client, ctx context.Context, id int) (*User, error) {
	return &User{
		ID:       1,
		Username: "John Doe",
		Email:    "",
		Provider: "google",
		Links:    []Link{},
	}, nil
}
