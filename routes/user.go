package routes

import (
	"context"
	"errors"

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

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Provider string `json:"provider"`
}

type UpdateUserRequest struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

type NormalResponse struct {
	Message string `json:"message"`
}

func CreateUser(pdb *pgxpool.Pool, rdb *redis.Client, ctx context.Context, data *CreateUserRequest) (*User, error) {
	exists, chkerr := checkUsername(pdb, ctx, data.Username)
	if chkerr != nil {
		return nil, chkerr
	}
	if exists {
		return nil, errors.New("username already exists")
	}
	var user User
	err := pdb.QueryRow(ctx, `INSERT INTO users (username, email, provider) VALUES ($1, $2, $3) RETURNING id, username, email, provider`, data.Username, data.Email, data.Provider).Scan(&user.ID, &user.Username, &user.Email, &user.Provider)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(pdb *pgxpool.Pool, rdb *redis.Client, ctx context.Context, data *UpdateUserRequest) (*User, error) {
	exists, chkerr := checkUsername(pdb, ctx, data.Username)
	if chkerr != nil {
		return nil, chkerr
	}
	if exists {
		return nil, errors.New("username already exists")
	}
	var user User
	err := pdb.QueryRow(ctx, `UPDATE users SET username = $1, WHERE id = $2 RETURNING id, username, email, provider`, data.Username, data.Id).Scan(&user.ID, &user.Username, &user.Email, &user.Provider)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func DeleteUser(pdb *pgxpool.Pool, rdb *redis.Client, ctx context.Context, id int) error {
	_, err := pdb.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(pdb *pgxpool.Pool, rdb *redis.Client, ctx context.Context, id int) (*User, error) {
	var user User
	err := pdb.QueryRow(ctx, `SELECT id, username, email, provider FROM users WHERE id = $1`, id).Scan(&user.ID, &user.Username, &user.Email, &user.Provider)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func checkUsername(pdb *pgxpool.Pool, ctx context.Context, username string) (bool, error) {
	var exists bool
	err := pdb.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
