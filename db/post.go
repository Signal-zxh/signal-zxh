package db

import (
	"database/sql"
	"errors"

	"github.com/Signal-zxh/signal-zxh/model"
)

var ErrNoRowsAffected = errors.New("no rows affected")
var ErrNotFound = errors.New("not found")

func GetPosts() ([]model.Post, error) {
	rows, err := DB.Query("SELECT id, title FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post

	for rows.Next() {
		var post model.Post

		if err := rows.Scan(&post.ID, &post.Title); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func CreatePost(title string) (int64, error) {
	res, err := DB.Exec(
		"INSERT INTO posts(title) VALUES(?)",
		title,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func UpdatePost(id int, title string) error {
	res, err := DB.Exec("UPDATE posts SET title = ? WHERE id = ?", title, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNoRowsAffected
	}

	return nil
}

func DeletePost(id int) error {
	res, err := DB.Exec("DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNoRowsAffected
	}

	return nil
}

func GetPostByID(id int) (model.Post, error) {
	row := DB.QueryRow("SELECT id, title FROM posts WHERE id = ?", id)

	var post model.Post

	err := row.Scan(&post.ID, &post.Title)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Post{}, ErrNotFound
		}
		return model.Post{}, err
	}

	return post, nil
}
