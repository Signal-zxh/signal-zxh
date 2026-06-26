package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Signal-zxh/signal-zxh/db"
	"github.com/Signal-zxh/signal-zxh/model"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
)

func GetPostByID(id int) (model.Post, error) {
	// 从redis中查询帖子
	key := fmt.Sprintf("post:%d", id)
	val, err := db.RDB.Get(context.Background(), key).Result()
	if err == nil {
		fmt.Println("hit redis")
		var post model.Post
		_ = json.Unmarshal([]byte(val), &post)
		return post, nil
	}
	// 从数据库中查询帖子
	post, err := db.GetPostByID(id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return model.Post{}, ErrNotFound
		}
		return model.Post{}, err
	}
	// 回写redis缓存
	b, _ := json.Marshal(post)
	db.RDB.Set(
		context.Background(),
		key,
		b,
		10*time.Minute,
	)

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
	// 删除redis缓存，确保下次查询获取最新数据
	db.RDB.Del(context.Background(), fmt.Sprintf("post:%d", id))
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
	// 删除redis缓存，确保下次查询获取最新数据
	db.RDB.Del(context.Background(), fmt.Sprintf("post:%d", id))
	return nil
}
