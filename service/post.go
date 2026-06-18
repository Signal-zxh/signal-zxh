package service

import (
	"errors"

	"github.com/Signal-zxh/signal-zxh/db"
	"github.com/Signal-zxh/signal-zxh/model"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
)

func GetPostByID(id int) (model.Post, error) {
	post, err := db.GetPostByID(id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return model.Post{}, ErrNotFound
		}
		return model.Post{}, err
	}
	return post, nil
}

func GetPosts() ([]model.Post, error) {
	return db.GetPosts()
}

func CreatePost(title, content string, userID int) (int64, error) {
	if title == "" || len(title) > 100 {
		return 0, ErrInvalidInput
	}

	post := model.Post{
		Title:   title,
		Content: content,
		UserID:  userID,
	}

	return db.CreatePost(post)
}

func UpdatePost(id int, title, content string) error {
	if title == "" || len(title) > 100 {
		return ErrInvalidInput
	}
	post := model.Post{
		ID: 	 id,
		Title:   title,
		Content: content,
	}
	err := db.UpdatePost(post)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsAffected) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func DeletePost(id int) error {
	err := db.DeletePost(id)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsAffected) {
			return ErrNotFound
		}
		return err
	}
	return nil
}
