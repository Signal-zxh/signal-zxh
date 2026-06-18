package db

import (
	"database/sql"
	"errors"

	"github.com/Signal-zxh/signal-zxh/model"
)

var ErrNoRowsAffected = errors.New("no rows affected")
var ErrNotFound = errors.New("not found")

func GetPosts() ([]model.Post, error) {
	rows, err := DB.Query("SELECT id, title, content, user_id FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post

	for rows.Next() {
		var post model.Post

		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func CreatePost(post model.Post) (int64, error) {
	res, err := DB.Exec(
		"INSERT INTO posts(title, content, user_id) VALUES(?, ?,?)",
		post.Title, post.Content, post.UserID,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func UpdatePost(post model.Post) error {
	res, err := DB.Exec("UPDATE posts SET title = ?, content = ? WHERE id = ?", post.Title, post.Content, post.ID)
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
	row := DB.QueryRow("SELECT id, title, content, user_id FROM posts WHERE id = ?", id)

	var post model.Post

	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Post{}, ErrNotFound
		}
		return model.Post{}, err
	}

	return post, nil
}
