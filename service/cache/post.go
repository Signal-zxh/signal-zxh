package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Signal-zxh/signal-zxh/db"
	"github.com/Signal-zxh/signal-zxh/model"
)

const (
	RedisNil = "__nil__"
)

func GetPostByID(id int) (model.Post, bool, error) {
	key := fmt.Sprintf("post:%d", id)
	val, err := db.RDB.Get(context.Background(), key).Result()
	if err != nil {
		return model.Post{}, false, err
	}

	if val == RedisNil {
		return model.Post{}, true, nil
	}

	var post model.Post
	if err := json.Unmarshal([]byte(val), &post); err != nil {
		return model.Post{}, false, err
	}
	return post, true, nil
}

func SetPost(post model.Post, ttl time.Duration) error {
	key := fmt.Sprintf("post:%d", post.ID)
	b, err := json.Marshal(post)
	if err != nil {
		return err
	}
	return db.RDB.Set(context.Background(), key, b, ttl).Err()
}

func SetNilPost(id int, ttl time.Duration) error {
	key := fmt.Sprintf("post:%d", id)
	return db.RDB.Set(context.Background(), key, RedisNil, ttl).Err()
}

func GetPosts() ([]model.Post, bool, error) {
	key := "posts:list"
	val, err := db.RDB.Get(context.Background(), key).Result()
	if err != nil {
		return nil, false, err
	}

	var posts []model.Post
	if err := json.Unmarshal([]byte(val), &posts); err != nil {
		return nil, false, err
	}
	return posts, true, nil
}

func SetPosts(posts []model.Post, ttl time.Duration) error {
	key := "posts:list"
	b, err := json.Marshal(posts)
	if err != nil {
		return err
	}
	return db.RDB.Set(context.Background(), key, b, ttl).Err()
}

func InvalidatePost(id int) error {
	ctx := context.Background()
	db.RDB.Del(ctx, fmt.Sprintf("post:%d", id))
	db.RDB.Del(ctx, "posts:list")
	return nil
}

func InvalidatePosts() error {
	ctx := context.Background()
	db.RDB.Del(ctx, "posts:list")
	return nil
}
