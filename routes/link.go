package routes

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Link struct {
	Id     int    `json:"id"`
	UserId int    `json:"userId"`
	Url    string `json:"url"`
	Name   string `json:"name"`
}

type CreateLinkRequest struct {
	UserId int    `json:"userId"`
	Url    string `json:"url"`
	Name   string `json:"name"`
}

type UpdateLinkRequest struct {
	Id     int    `json:"id"`
	UserId int    `json:"userId"`
	Url    string `json:"url"`
	Name   string `json:"name"`
}

func CreateLink(pdb *pgxpool.Pool, rdb *redis.Client, ctx context.Context, data *CreateLinkRequest) (*Link, error) {
	var link Link
	err := pdb.QueryRow(ctx, `INSERT INTO links (user_id, url, name) VALUES ($1, $2, $3) RETURNING id, user_id, url, name`, data.UserId, data.Url, data.Name).Scan(&link.Id, &link.UserId, &link.Url, &link.Name)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func GetLink(pdb *pgxpool.Pool, rdb *redis.Client, ctx context.Context, id int) (*Link, error) {
	var link Link
	err := pdb.QueryRow(ctx, `SELECT id, user_id, url, name FROM links WHERE id = $1`, id).Scan(&link.Id, &link.UserId, &link.Url, &link.Name)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func UpdateLink(pdb *pgxpool.Pool, rdb *redis.Client, ctx context.Context, id int, data *UpdateLinkRequest) (*Link, error) {
	exists, chkErr := checkNameforLinks(pdb, ctx, data.Name, data.UserId)
	if chkErr != nil {
		return nil, chkErr
	}
	if exists {
		return nil, errors.New("name already exists")
	}
	var link Link
	err := pdb.QueryRow(ctx, `UPDATE links SET url = $1, name = $2 WHERE id = $3 RETURNING id, user_id, url, name`, data.Url, data.Name, data.Id).Scan(&link.Id, &link.UserId, &link.Url, &link.Name)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func DeleteLink(pdb *pgxpool.Pool, rdb *redis.Client, ctx context.Context, id int) error {
	_, err := pdb.Exec(ctx, `DELETE FROM links WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func GetLinksByUserId(pdb *pgxpool.Pool, rdb *redis.Client, ctx context.Context, userId int) ([]Link, error) {
	rows, err := pdb.Query(ctx, `SELECT id, user_id, url, name FROM links WHERE user_id = $1`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var links []Link
	for rows.Next() {
		var link Link
		err := rows.Scan(&link.Id, &link.UserId, &link.Url, &link.Name)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}

func checkNameforLinks(pdb *pgxpool.Pool, ctx context.Context, name string, userId int) (bool, error) {
	var count int
	err := pdb.QueryRow(ctx, `SELECT COUNT(*) FROM links WHERE name = $1 AND user_id = $2`, name, userId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
