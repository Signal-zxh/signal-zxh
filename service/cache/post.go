package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Signal-zxh/signalzxh-blog/db"
	"github.com/Signal-zxh/signalzxh-blog/model"
)

const (
	RedisNil = "__nil__"
)

type PostCache interface {
	GetPostByID(id int) (model.Post, bool, error)
	SetPost(post model.Post, ttl time.Duration) error
	SetNilPost(id int, ttl time.Duration) error
	GetPosts() ([]model.Post, bool, error)
	SetPosts(posts []model.Post, ttl time.Duration) error
	GetPostsByPage(page, pageSize int) ([]model.Post, bool, error)
	SetPostsByPage(posts []model.Post, page, pageSize int, ttl time.Duration) error
	InvalidatePost(id int) error
	InvalidatePosts() error
}

type postCache struct{}

var PostCacheImpl = &postCache{}

func (c *postCache) GetPostByID(id int) (model.Post, bool, error) {
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

func (c *postCache) SetPost(post model.Post, ttl time.Duration) error {
	key := fmt.Sprintf("post:%d", post.ID)
	b, err := json.Marshal(post)
	if err != nil {
		return err
	}
	return db.RDB.Set(context.Background(), key, b, ttl).Err()
}

func (c *postCache) SetNilPost(id int, ttl time.Duration) error {
	key := fmt.Sprintf("post:%d", id)
	return db.RDB.Set(context.Background(), key, RedisNil, ttl).Err()
}

func (c *postCache) GetPosts() ([]model.Post, bool, error) {
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

func (c *postCache) SetPosts(posts []model.Post, ttl time.Duration) error {
	key := "posts:list"
	b, err := json.Marshal(posts)
	if err != nil {
		return err
	}
	return db.RDB.Set(context.Background(), key, b, ttl).Err()
}

func (c *postCache) GetPostsByPage(page, pageSize int) ([]model.Post, bool, error) {
	key := fmt.Sprintf("posts:list:page:%d:size:%d", page, pageSize)
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

func (c *postCache) SetPostsByPage(posts []model.Post, page, pageSize int, ttl time.Duration) error {
	key := fmt.Sprintf("posts:list:page:%d:size:%d", page, pageSize)
	b, err := json.Marshal(posts)
	if err != nil {
		return err
	}
	return db.RDB.Set(context.Background(), key, b, ttl).Err()
}

func (c *postCache) InvalidatePost(id int) error {
	ctx := context.Background()
	db.RDB.Del(ctx, fmt.Sprintf("post:%d", id))
	keys, _ := db.RDB.Keys(ctx, "posts:list:*").Result() //nolint:errcheck
	if len(keys) > 0 {
		_ = db.RDB.Del(ctx, keys...).Err() //nolint:errcheck
	}
	return nil
}

func (c *postCache) InvalidatePosts() error {
	ctx := context.Background()
	keys, _ := db.RDB.Keys(ctx, "posts:list:*").Result() //nolint:errcheck
	if len(keys) > 0 {
		_ = db.RDB.Del(ctx, keys...).Err() //nolint:errcheck
	}
	return nil
}
