package service

import (
	"errors"
	"time"

	"github.com/Signal-zxh/signal-zxh/db"
	"github.com/Signal-zxh/signal-zxh/model"
	"github.com/Signal-zxh/signal-zxh/service/cache"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
)

func GetPostByID(id int) (model.Post, error) {
	post, found, err := cache.GetPostByID(id)
	if err == nil && found {
		if post.ID == 0 {
			return model.Post{}, ErrNotFound
		}
		return post, nil
	}

	post, err = db.GetPostByID(id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			cache.SetNilPost(id, 1*time.Minute)
			return model.Post{}, ErrNotFound
		}
		return model.Post{}, err
	}

	cache.SetPost(post, 10*time.Minute)
	return post, nil
}

func GetPosts() ([]model.Post, error) {
	posts, found, err := cache.GetPosts()
	if err == nil && found {
		return posts, nil
	}

	posts, err = db.GetPosts()
	if err != nil {
		return nil, err
	}

	cache.SetPosts(posts, 10*time.Minute)
	return posts, nil
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

	id, err := db.CreatePost(post)
	if err != nil {
		return 0, err
	}

	cache.InvalidatePosts()
	return id, nil
}

func UpdatePost(id int, title, content string) error {
	if title == "" || len(title) > 100 {
		return ErrInvalidInput
	}

	post := model.Post{
		ID:      id,
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

	cache.InvalidatePost(id)
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

	cache.InvalidatePost(id)
	return nil
}
