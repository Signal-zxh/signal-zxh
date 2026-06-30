package db

import (
	"database/sql"
	"errors"

	"github.com/Signal-zxh/signalzxh-blog/model"
)

var ErrNoRowsAffected = errors.New("no rows affected")
var ErrNotFound = errors.New("not found")

type PostRepo interface {
	GetPosts() ([]model.Post, error)
	GetPostsByPage(page, pageSize int) ([]model.Post, error)
	GetPostsCount() (int, error)
	CreatePost(post model.Post) (int64, error)
	UpdatePost(post model.Post) error
	DeletePost(id int) error
	GetPostByID(id int) (model.Post, error)
}

type postRepo struct{}

var PostRepoImpl = &postRepo{}

func (r *postRepo) GetPosts() ([]model.Post, error) {
	rows, err := DB.Query("SELECT id, title, content, user_id FROM posts ORDER BY id DESC")
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

func (r *postRepo) GetPostsByPage(page, pageSize int) ([]model.Post, error) {
	offset := (page - 1) * pageSize
	rows, err := DB.Query(
		"SELECT id, title, content, user_id FROM posts ORDER BY id DESC LIMIT ? OFFSET ?",
		pageSize, offset,
	)
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

func (r *postRepo) GetPostsCount() (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM posts").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *postRepo) CreatePost(post model.Post) (int64, error) {
	res, err := DB.Exec(
		"INSERT INTO posts(title, content, user_id) VALUES(?, ?,?)",
		post.Title, post.Content, post.UserID,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (r *postRepo) UpdatePost(post model.Post) error {
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

func (r *postRepo) DeletePost(id int) error {
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

func (r *postRepo) GetPostByID(id int) (model.Post, error) {
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
