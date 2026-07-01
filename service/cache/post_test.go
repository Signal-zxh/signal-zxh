package cache

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/Signal-zxh/signalzxh-blog/db"
	"github.com/Signal-zxh/signalzxh-blog/model"
	"github.com/redis/go-redis/v9"
)

func setupTestRedis(t *testing.T) *redis.Client {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skipf("redis not available: %v", err)
	}

	db.RDB = rdb
	return rdb
}

func TestGetPostByID(t *testing.T) {
	rdb := setupTestRedis(t)
	ctx := context.Background()

	key := "post:1"
	rdb.Del(ctx, key)

	post := model.Post{
		ID:      1,
		Title:   "test title",
		Content: "test content",
		UserID:  1,
	}

	b, _ := json.Marshal(post)
	rdb.Set(ctx, key, b, 0)

	cache := &postCache{}
	gotPost, found, err := cache.GetPostByID(1)
	if err != nil {
		t.Errorf("GetPostByID() error = %v", err)
	}
	if !found {
		t.Errorf("GetPostByID() expected found = true")
	}
	if gotPost.ID != post.ID {
		t.Errorf("GetPostByID() got ID = %d, want %d", gotPost.ID, post.ID)
	}

	rdb.Del(ctx, key)
}

func TestGetPostByID_NilPost(t *testing.T) {
	rdb := setupTestRedis(t)
	ctx := context.Background()

	key := "post:999"
	rdb.Del(ctx, key)

	rdb.Set(ctx, key, RedisNil, 0)

	cache := &postCache{}
	gotPost, found, err := cache.GetPostByID(999)
	if err != nil {
		t.Errorf("GetPostByID() error = %v", err)
	}
	if !found {
		t.Errorf("GetPostByID() expected found = true")
	}
	if gotPost.ID != 0 {
		t.Errorf("GetPostByID() nil post got ID = %d, want 0", gotPost.ID)
	}

	rdb.Del(ctx, key)
}

func TestGetPostByID_NotFound(t *testing.T) {
	rdb := setupTestRedis(t)
	ctx := context.Background()

	key := "post:999"
	rdb.Del(ctx, key)

	cache := &postCache{}
	_, found, err := cache.GetPostByID(999)
	if err != nil && err != redis.Nil {
		t.Errorf("GetPostByID() error = %v", err)
	}
	if found {
		t.Errorf("GetPostByID() expected found = false")
	}
}

func TestSetPost(t *testing.T) {
	rdb := setupTestRedis(t)
	ctx := context.Background()

	key := "post:2"
	rdb.Del(ctx, key)

	post := model.Post{
		ID:      2,
		Title:   "test title",
		Content: "test content",
		UserID:  1,
	}

	cache := &postCache{}
	err := cache.SetPost(post, 1*time.Minute)
	if err != nil {
		t.Errorf("SetPost() error = %v", err)
	}

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		t.Errorf("SetPost() failed to set key, error = %v", err)
	}

	var gotPost model.Post
	json.Unmarshal([]byte(val), &gotPost)
	if gotPost.ID != post.ID {
		t.Errorf("SetPost() got ID = %d, want %d", gotPost.ID, post.ID)
	}

	rdb.Del(ctx, key)
}

func TestGetPostsByPage(t *testing.T) {
	rdb := setupTestRedis(t)
	ctx := context.Background()

	key := "posts:list:page:1:size:10"
	rdb.Del(ctx, key)

	posts := []model.Post{
		{ID: 1, Title: "post 1", Content: "content 1", UserID: 1},
		{ID: 2, Title: "post 2", Content: "content 2", UserID: 1},
	}

	b, _ := json.Marshal(posts)
	rdb.Set(ctx, key, b, 0)

	cache := &postCache{}
	gotPosts, found, err := cache.GetPostsByPage(1, 10)
	if err != nil {
		t.Errorf("GetPostsByPage() error = %v", err)
	}
	if !found {
		t.Errorf("GetPostsByPage() expected found = true")
	}
	if len(gotPosts) != len(posts) {
		t.Errorf("GetPostsByPage() got length = %d, want %d", len(gotPosts), len(posts))
	}

	rdb.Del(ctx, key)
}

func TestSetPostsByPage(t *testing.T) {
	rdb := setupTestRedis(t)
	ctx := context.Background()

	key := "posts:list:page:2:size:5"
	rdb.Del(ctx, key)

	posts := []model.Post{
		{ID: 3, Title: "post 3", Content: "content 3", UserID: 1},
	}

	cache := &postCache{}
	err := cache.SetPostsByPage(posts, 2, 5, 1*time.Minute)
	if err != nil {
		t.Errorf("SetPostsByPage() error = %v", err)
	}

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		t.Errorf("SetPostsByPage() failed to set key, error = %v", err)
	}

	var gotPosts []model.Post
	json.Unmarshal([]byte(val), &gotPosts)
	if len(gotPosts) != len(posts) {
		t.Errorf("SetPostsByPage() got length = %d, want %d", len(gotPosts), len(posts))
	}

	rdb.Del(ctx, key)
}

func TestInvalidatePost(t *testing.T) {
	rdb := setupTestRedis(t)
	ctx := context.Background()

	rdb.Del(ctx, "post:1")
	rdb.Set(ctx, "post:1", "test value", 0)

	cache := &postCache{}
	err := cache.InvalidatePost(1)
	if err != nil {
		t.Errorf("InvalidatePost() error = %v", err)
	}

	_, err = rdb.Get(ctx, "post:1").Result()
	if err != redis.Nil {
		t.Errorf("InvalidatePost() failed to delete post:1")
	}
}

func TestGetPosts_Success(t *testing.T) {
	rdb := setupTestRedis(t)
	ctx := context.Background()

	key := "posts:list"
	rdb.Del(ctx, key)

	posts := []model.Post{
		{ID: 1, Title: "post 1", Content: "content 1", UserID: 1},
		{ID: 2, Title: "post 2", Content: "content 2", UserID: 1},
	}

	b, _ := json.Marshal(posts)
	rdb.Set(ctx, key, b, 0)

	cache := &postCache{}
	gotPosts, found, err := cache.GetPosts()
	if err != nil {
		t.Errorf("GetPosts() error = %v", err)
	}
	if !found {
		t.Errorf("GetPosts() expected found = true")
	}
	if len(gotPosts) != len(posts) {
		t.Errorf("GetPosts() got length = %d, want %d", len(gotPosts), len(posts))
	}
	if gotPosts[0].Title != posts[0].Title {
		t.Errorf("GetPosts()[0].Title = %s, want %s", gotPosts[0].Title, posts[0].Title)
	}

	rdb.Del(ctx, key)
}

func TestGetPosts_NotFound(t *testing.T) {
	rdb := setupTestRedis(t)
	ctx := context.Background()

	key := "posts:list"
	rdb.Del(ctx, key)

	cache := &postCache{}
	_, found, err := cache.GetPosts()
	if err != nil && err != redis.Nil {
		t.Errorf("GetPosts() error = %v", err)
	}
	if found {
		t.Errorf("GetPosts() expected found = false")
	}
}

func TestGetPosts_InvalidJSON(t *testing.T) {
	rdb := setupTestRedis(t)
	ctx := context.Background()

	key := "posts:list"
	rdb.Del(ctx, key)

	rdb.Set(ctx, key, "invalid json data", 0)

	cache := &postCache{}
	_, found, err := cache.GetPosts()
	if err == nil {
		t.Errorf("GetPosts() expected error for invalid JSON")
	}
	if found {
		t.Errorf("GetPosts() expected found = false for invalid JSON")
	}

	rdb.Del(ctx, key)
}

func TestSetPosts(t *testing.T) {
	rdb := setupTestRedis(t)
	ctx := context.Background()

	key := "posts:list"
	rdb.Del(ctx, key)

	posts := []model.Post{
		{ID: 1, Title: "post 1", Content: "content 1", UserID: 1},
		{ID: 2, Title: "post 2", Content: "content 2", UserID: 1},
	}

	cache := &postCache{}
	err := cache.SetPosts(posts, 1*time.Minute)
	if err != nil {
		t.Errorf("SetPosts() error = %v", err)
	}

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		t.Errorf("SetPosts() failed to set key, error = %v", err)
	}

	var gotPosts []model.Post
	json.Unmarshal([]byte(val), &gotPosts)
	if len(gotPosts) != len(posts) {
		t.Errorf("SetPosts() got length = %d, want %d", len(gotPosts), len(posts))
	}

	rdb.Del(ctx, key)
}

func TestSetNilPost(t *testing.T) {
	rdb := setupTestRedis(t)
	ctx := context.Background()

	key := "post:999"
	rdb.Del(ctx, key)

	cache := &postCache{}
	err := cache.SetNilPost(999, 1*time.Minute)
	if err != nil {
		t.Errorf("SetNilPost() error = %v", err)
	}

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		t.Errorf("SetNilPost() failed to set key, error = %v", err)
	}
	if val != RedisNil {
		t.Errorf("SetNilPost() got value = %s, want %s", val, RedisNil)
	}

	rdb.Del(ctx, key)
}

func TestInvalidatePosts(t *testing.T) {
	rdb := setupTestRedis(t)
	ctx := context.Background()

	key1 := "posts:list:page:1:size:10"
	key2 := "posts:list:page:2:size:10"

	rdb.Del(ctx, key1, key2)
	rdb.Set(ctx, key1, "test value 1", 0)
	rdb.Set(ctx, key2, "test value 2", 0)

	cache := &postCache{}
	err := cache.InvalidatePosts()
	if err != nil {
		t.Errorf("InvalidatePosts() error = %v", err)
	}

	_, err1 := rdb.Get(ctx, key1).Result()
	_, err2 := rdb.Get(ctx, key2).Result()

	if err1 != redis.Nil {
		t.Errorf("InvalidatePosts() failed to delete %s", key1)
	}
	if err2 != redis.Nil {
		t.Errorf("InvalidatePosts() failed to delete %s", key2)
	}
}

func TestGetPostByID_InvalidJSON(t *testing.T) {
	rdb := setupTestRedis(t)
	ctx := context.Background()

	key := "post:3"
	rdb.Del(ctx, key)

	rdb.Set(ctx, key, "invalid json", 0)

	cache := &postCache{}
	_, found, err := cache.GetPostByID(3)
	if err == nil {
		t.Errorf("GetPostByID() expected error for invalid JSON")
	}
	if found {
		t.Errorf("GetPostByID() expected found = false for invalid JSON")
	}

	rdb.Del(ctx, key)
}

func TestGetPostsByPage_NotFound(t *testing.T) {
	rdb := setupTestRedis(t)
	ctx := context.Background()

	key := "posts:list:page:99:size:10"
	rdb.Del(ctx, key)

	cache := &postCache{}
	_, found, err := cache.GetPostsByPage(99, 10)
	if err != nil && err != redis.Nil {
		t.Errorf("GetPostsByPage() error = %v", err)
	}
	if found {
		t.Errorf("GetPostsByPage() expected found = false")
	}
}

func TestGetPostsByPage_InvalidJSON(t *testing.T) {
	rdb := setupTestRedis(t)
	ctx := context.Background()

	key := "posts:list:page:3:size:10"
	rdb.Del(ctx, key)

	rdb.Set(ctx, key, "invalid json", 0)

	cache := &postCache{}
	_, found, err := cache.GetPostsByPage(3, 10)
	if err == nil {
		t.Errorf("GetPostsByPage() expected error for invalid JSON")
	}
	if found {
		t.Errorf("GetPostsByPage() expected found = false for invalid JSON")
	}

	rdb.Del(ctx, key)
}

func TestInvalidatePost_WithPageCache(t *testing.T) {
	rdb := setupTestRedis(t)
	ctx := context.Background()

	postKey := "post:1"
	pageKey := "posts:list:page:1:size:10"

	rdb.Del(ctx, postKey, pageKey)
	rdb.Set(ctx, postKey, "post value", 0)
	rdb.Set(ctx, pageKey, "page value", 0)

	cache := &postCache{}
	err := cache.InvalidatePost(1)
	if err != nil {
		t.Errorf("InvalidatePost() error = %v", err)
	}

	_, err1 := rdb.Get(ctx, postKey).Result()
	if err1 != redis.Nil {
		t.Errorf("InvalidatePost() failed to delete post key")
	}

	_, err2 := rdb.Get(ctx, pageKey).Result()
	if err2 != redis.Nil {
		t.Errorf("InvalidatePost() failed to delete page cache")
	}
}
