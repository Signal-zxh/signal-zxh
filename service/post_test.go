package service

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Signal-zxh/signalzxh-blog/db"
	"github.com/Signal-zxh/signalzxh-blog/model"
)

type spyPostRepo struct {
	getPostByIDCalled    bool
	getPostByIDReturn    model.Post
	getPostByIDErr       error
	getPostsByPageCalled bool
	getPostsByPageReturn []model.Post
	getPostsByPageErr    error
	getPostsCountCalled  bool
	getPostsCountReturn  int
	getPostsCountErr     error
	createPostCalled     bool
	createPostReturn     int64
	createPostErr        error
	deletePostCalled     bool
	deletePostErr        error
	updatePostCalled     bool
	updatePostErr        error
	getPostsCalled       bool
	getPostsReturn       []model.Post
	getPostsErr          error
}

func (s *spyPostRepo) GetPosts() ([]model.Post, error) {
	s.getPostsCalled = true
	return s.getPostsReturn, s.getPostsErr
}

func (s *spyPostRepo) GetPostsByPage(page, pageSize int) ([]model.Post, error) {
	s.getPostsByPageCalled = true
	return s.getPostsByPageReturn, s.getPostsByPageErr
}

func (s *spyPostRepo) GetPostsCount() (int, error) {
	s.getPostsCountCalled = true
	return s.getPostsCountReturn, s.getPostsCountErr
}

func (s *spyPostRepo) CreatePost(post model.Post) (int64, error) {
	s.createPostCalled = true
	return s.createPostReturn, s.createPostErr
}

func (s *spyPostRepo) UpdatePost(post model.Post) error {
	s.updatePostCalled = true
	return s.updatePostErr
}

func (s *spyPostRepo) DeletePost(id int) error {
	s.deletePostCalled = true
	return s.deletePostErr
}

func (s *spyPostRepo) GetPostByID(id int) (model.Post, error) {
	s.getPostByIDCalled = true
	return s.getPostByIDReturn, s.getPostByIDErr
}

type spyPostCache struct {
	getPostByIDCalled     bool
	getPostByIDReturn     model.Post
	getPostByIDFound      bool
	getPostByIDErr        error
	setPostCalled         bool
	setPostArg            model.Post
	setNilPostCalled      bool
	setNilPostID          int
	getPostsByPageCalled  bool
	getPostsByPageReturn  []model.Post
	getPostsByPageFound   bool
	getPostsByPageErr     error
	setPostsByPageCalled  bool
	invalidatePostCalled  bool
	invalidatePostID      int
	invalidatePostsCalled bool
	getPostsCalled        bool
	getPostsReturn        []model.Post
	getPostsFound         bool
	getPostsErr           error
	setPostsCalled        bool
}

func (s *spyPostCache) GetPostByID(id int) (model.Post, bool, error) {
	s.getPostByIDCalled = true
	return s.getPostByIDReturn, s.getPostByIDFound, s.getPostByIDErr
}

func (s *spyPostCache) SetPost(post model.Post, ttl time.Duration) error {
	s.setPostCalled = true
	s.setPostArg = post
	return nil
}

func (s *spyPostCache) SetNilPost(id int, ttl time.Duration) error {
	s.setNilPostCalled = true
	s.setNilPostID = id
	return nil
}

func (s *spyPostCache) GetPosts() ([]model.Post, bool, error) {
	s.getPostsCalled = true
	return s.getPostsReturn, s.getPostsFound, s.getPostsErr
}

func (s *spyPostCache) SetPosts(posts []model.Post, ttl time.Duration) error {
	s.setPostsCalled = true
	return nil
}

func (s *spyPostCache) GetPostsByPage(page, pageSize int) ([]model.Post, bool, error) {
	s.getPostsByPageCalled = true
	return s.getPostsByPageReturn, s.getPostsByPageFound, s.getPostsByPageErr
}

func (s *spyPostCache) SetPostsByPage(posts []model.Post, page, pageSize int, ttl time.Duration) error {
	s.setPostsByPageCalled = true
	return nil
}

func (s *spyPostCache) InvalidatePost(id int) error {
	s.invalidatePostCalled = true
	s.invalidatePostID = id
	return nil
}

func (s *spyPostCache) InvalidatePosts() error {
	s.invalidatePostsCalled = true
	return nil
}

func TestGetPostByID_CacheHit(t *testing.T) {
	expectedPost := model.Post{ID: 1, Title: "Cache Hit", Content: "Content"}

	spyRepo := &spyPostRepo{}
	spyCache := &spyPostCache{
		getPostByIDReturn: expectedPost,
		getPostByIDFound:  true,
	}

	service := NewPostService(spyRepo, spyCache)
	post, err := service.GetPostByID(1)

	if err != nil {
		t.Errorf("GetPostByID() error = %v", err)
	}

	if post.ID != expectedPost.ID || post.Title != expectedPost.Title {
		t.Errorf("GetPostByID() got = %v, want %v", post, expectedPost)
	}

	if !spyCache.getPostByIDCalled {
		t.Error("GetPostByID() should call cache.GetPostByID")
	}

	if spyRepo.getPostByIDCalled {
		t.Error("GetPostByID() should NOT call repo.GetPostByID when cache hits")
	}

	if spyCache.setPostCalled {
		t.Error("GetPostByID() should NOT call cache.SetPost when cache hits")
	}
}

func TestGetPostByID_CacheMiss_DBHit(t *testing.T) {
	expectedPost := model.Post{ID: 1, Title: "DB Hit", Content: "Content"}

	spyRepo := &spyPostRepo{
		getPostByIDReturn: expectedPost,
	}
	spyCache := &spyPostCache{
		getPostByIDFound: false,
	}

	service := NewPostService(spyRepo, spyCache)
	post, err := service.GetPostByID(1)

	if err != nil {
		t.Errorf("GetPostByID() error = %v", err)
	}

	if post.ID != expectedPost.ID || post.Title != expectedPost.Title {
		t.Errorf("GetPostByID() got = %v, want %v", post, expectedPost)
	}

	if !spyCache.getPostByIDCalled {
		t.Error("GetPostByID() should call cache.GetPostByID")
	}

	if !spyRepo.getPostByIDCalled {
		t.Error("GetPostByID() should call repo.GetPostByID when cache misses")
	}

	if !spyCache.setPostCalled {
		t.Error("GetPostByID() should call cache.SetPost after DB hit")
	}

	if spyCache.setPostArg.ID != expectedPost.ID {
		t.Errorf("cache.SetPost() got post ID = %v, want %v", spyCache.setPostArg.ID, expectedPost.ID)
	}
}

func TestGetPostByID_DBNotFound(t *testing.T) {
	spyRepo := &spyPostRepo{
		getPostByIDErr: db.ErrNotFound,
	}
	spyCache := &spyPostCache{
		getPostByIDFound: false,
	}

	service := NewPostService(spyRepo, spyCache)
	_, err := service.GetPostByID(999)

	if err == nil {
		t.Error("GetPostByID() should return error when DB not found")
	}

	if !spyCache.getPostByIDCalled {
		t.Error("GetPostByID() should call cache.GetPostByID")
	}

	if !spyRepo.getPostByIDCalled {
		t.Error("GetPostByID() should call repo.GetPostByID when cache misses")
	}

	if !spyCache.setNilPostCalled {
		t.Error("GetPostByID() should call cache.SetNilPost when DB not found")
	}

	if spyCache.setNilPostID != 999 {
		t.Errorf("cache.SetNilPost() got ID = %v, want %v", spyCache.setNilPostID, 999)
	}
}

func TestCreatePost_Success(t *testing.T) {
	spyRepo := &spyPostRepo{
		createPostReturn: 1,
	}
	spyCache := &spyPostCache{}

	service := NewPostService(spyRepo, spyCache)
	id, err := service.CreatePost("Test Title", "Test Content", 1)

	if err != nil {
		t.Errorf("CreatePost() error = %v", err)
	}

	if id != 1 {
		t.Errorf("CreatePost() got id = %v, want %v", id, 1)
	}

	if !spyRepo.createPostCalled {
		t.Error("CreatePost() should call repo.CreatePost")
	}

	if !spyCache.invalidatePostsCalled {
		t.Error("CreatePost() should call cache.InvalidatePosts after success")
	}
}

func TestCreatePost_Failure(t *testing.T) {
	spyRepo := &spyPostRepo{
		createPostErr: errors.New("DB error"),
	}
	spyCache := &spyPostCache{}

	service := NewPostService(spyRepo, spyCache)
	_, err := service.CreatePost("Test Title", "Test Content", 1)

	if err == nil {
		t.Error("CreatePost() should return error when DB fails")
	}

	if !spyRepo.createPostCalled {
		t.Error("CreatePost() should call repo.CreatePost")
	}

	if spyCache.invalidatePostsCalled {
		t.Error("CreatePost() should NOT call cache.InvalidatePosts when DB fails")
	}
}

func TestDeletePost_Success(t *testing.T) {
	spyRepo := &spyPostRepo{}
	spyCache := &spyPostCache{}

	service := NewPostService(spyRepo, spyCache)
	err := service.DeletePost(1)

	if err != nil {
		t.Errorf("DeletePost() error = %v", err)
	}

	if !spyRepo.deletePostCalled {
		t.Error("DeletePost() should call repo.DeletePost")
	}

	if !spyCache.invalidatePostCalled {
		t.Error("DeletePost() should call cache.InvalidatePost")
	}

	if spyCache.invalidatePostID != 1 {
		t.Errorf("cache.InvalidatePost() got ID = %v, want %v", spyCache.invalidatePostID, 1)
	}

	if !spyCache.invalidatePostsCalled {
		t.Error("DeletePost() should call cache.InvalidatePosts")
	}
}

func TestDeletePost_Failure(t *testing.T) {
	spyRepo := &spyPostRepo{
		deletePostErr: errors.New("DB error"),
	}
	spyCache := &spyPostCache{}

	service := NewPostService(spyRepo, spyCache)
	err := service.DeletePost(1)

	if err == nil {
		t.Error("DeletePost() should return error when DB fails")
	}

	if !spyRepo.deletePostCalled {
		t.Error("DeletePost() should call repo.DeletePost")
	}

	if spyCache.invalidatePostCalled {
		t.Error("DeletePost() should NOT call cache.InvalidatePost when DB fails")
	}

	if spyCache.invalidatePostsCalled {
		t.Error("DeletePost() should NOT call cache.InvalidatePosts when DB fails")
	}
}

func TestGetPostsByPage_CacheHit(t *testing.T) {
	expectedPosts := []model.Post{
		{ID: 1, Title: "Post 1"},
		{ID: 2, Title: "Post 2"},
	}
	expectedTotal := 100

	spyRepo := &spyPostRepo{
		getPostsCountReturn: expectedTotal,
	}
	spyCache := &spyPostCache{
		getPostsByPageReturn: expectedPosts,
		getPostsByPageFound:  true,
	}

	service := NewPostService(spyRepo, spyCache)
	posts, total, err := service.GetPostsByPage(1, 10)

	if err != nil {
		t.Errorf("GetPostsByPage() error = %v", err)
	}

	if len(posts) != len(expectedPosts) {
		t.Errorf("GetPostsByPage() got %d posts, want %d", len(posts), len(expectedPosts))
	}

	if total != expectedTotal {
		t.Errorf("GetPostsByPage() got total = %v, want %v", total, expectedTotal)
	}

	if !spyCache.getPostsByPageCalled {
		t.Error("GetPostsByPage() should call cache.GetPostsByPage")
	}

	if spyRepo.getPostsByPageCalled {
		t.Error("GetPostsByPage() should NOT call repo.GetPostsByPage when cache hits")
	}

	if !spyRepo.getPostsCountCalled {
		t.Error("GetPostsByPage() should call repo.GetPostsCount when cache hits")
	}

	if spyCache.setPostsByPageCalled {
		t.Error("GetPostsByPage() should NOT call cache.SetPostsByPage when cache hits")
	}
}

func TestGetPostsByPage_CacheMiss_DBHit(t *testing.T) {
	expectedPosts := []model.Post{
		{ID: 1, Title: "Post 1"},
		{ID: 2, Title: "Post 2"},
	}
	expectedTotal := 100

	spyRepo := &spyPostRepo{
		getPostsByPageReturn: expectedPosts,
		getPostsCountReturn:  expectedTotal,
	}
	spyCache := &spyPostCache{
		getPostsByPageFound: false,
	}

	service := NewPostService(spyRepo, spyCache)
	posts, total, err := service.GetPostsByPage(1, 10)

	if err != nil {
		t.Errorf("GetPostsByPage() error = %v", err)
	}

	if len(posts) != len(expectedPosts) {
		t.Errorf("GetPostsByPage() got %d posts, want %d", len(posts), len(expectedPosts))
	}

	if total != expectedTotal {
		t.Errorf("GetPostsByPage() got total = %v, want %v", total, expectedTotal)
	}

	if !spyCache.getPostsByPageCalled {
		t.Error("GetPostsByPage() should call cache.GetPostsByPage")
	}

	if !spyRepo.getPostsByPageCalled {
		t.Error("GetPostsByPage() should call repo.GetPostsByPage when cache misses")
	}

	if !spyRepo.getPostsCountCalled {
		t.Error("GetPostsByPage() should call repo.GetPostsCount when cache misses")
	}

	if !spyCache.setPostsByPageCalled {
		t.Error("GetPostsByPage() should call cache.SetPostsByPage after DB hit")
	}
}

func TestGetPostsByPage_ParameterValidation(t *testing.T) {
	service := NewPostService(&spyPostRepo{}, &spyPostCache{})

	tests := []struct {
		name     string
		page     int
		pageSize int
		wantPage int
		wantSize int
	}{
		{"valid_parameters", 1, 10, 1, 10},
		{"page_less_than_1", 0, 10, 1, 10},
		{"page_equals_0", -1, 10, 1, 10},
		{"pageSize_less_than_1", 1, 0, 1, 10},
		{"pageSize_equals_0", 1, -1, 1, 10},
		{"pageSize_greater_than_100", 1, 200, 1, 100},
		{"pageSize_equals_100", 1, 100, 1, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := service.GetPostsByPage(tt.page, tt.pageSize)
			if err != nil {
				t.Errorf("GetPostsByPage() error = %v", err)
			}
		})
	}
}

func TestCreatePost_Validation(t *testing.T) {
	service := NewPostService(&spyPostRepo{}, &spyPostCache{})

	tests := []struct {
		name    string
		title   string
		content string
		userID  int
		wantErr bool
	}{
		{"valid_title", "Test Title", "Test Content", 1, false},
		{"empty_title", "", "Test Content", 1, true},
		{"title_exactly_255_chars", strings.Repeat("a", 255), "content", 1, false},
		{"title_256_chars", strings.Repeat("a", 256), "content", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.CreatePost(tt.title, tt.content, tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetPosts_CacheHit(t *testing.T) {
	expectedPosts := []model.Post{
		{ID: 1, Title: "Post 1"},
		{ID: 2, Title: "Post 2"},
	}

	spyRepo := &spyPostRepo{}
	spyCache := &spyPostCache{
		getPostsReturn: expectedPosts,
		getPostsFound:  true,
	}

	service := NewPostService(spyRepo, spyCache)
	posts, err := service.GetPosts()

	if err != nil {
		t.Errorf("GetPosts() error = %v", err)
	}

	if len(posts) != len(expectedPosts) {
		t.Errorf("GetPosts() got %d posts, want %d", len(posts), len(expectedPosts))
	}

	if !spyCache.getPostsCalled {
		t.Error("GetPosts() should call cache.GetPosts")
	}

	if spyRepo.getPostsCalled {
		t.Error("GetPosts() should NOT call repo.GetPosts when cache hits")
	}
}

func TestGetPosts_CacheMiss_DBHit(t *testing.T) {
	expectedPosts := []model.Post{
		{ID: 1, Title: "Post 1"},
		{ID: 2, Title: "Post 2"},
	}

	spyRepo := &spyPostRepo{
		getPostsReturn: expectedPosts,
	}
	spyCache := &spyPostCache{
		getPostsFound: false,
	}

	service := NewPostService(spyRepo, spyCache)
	posts, err := service.GetPosts()

	if err != nil {
		t.Errorf("GetPosts() error = %v", err)
	}

	if len(posts) != len(expectedPosts) {
		t.Errorf("GetPosts() got %d posts, want %d", len(posts), len(expectedPosts))
	}

	if !spyCache.getPostsCalled {
		t.Error("GetPosts() should call cache.GetPosts")
	}

	if !spyRepo.getPostsCalled {
		t.Error("GetPosts() should call repo.GetPosts when cache misses")
	}

	if !spyCache.setPostsCalled {
		t.Error("GetPosts() should call cache.SetPosts after DB hit")
	}
}

func TestGetPosts_CacheMiss_DBError(t *testing.T) {
	spyRepo := &spyPostRepo{
		getPostsErr: errors.New("DB error"),
	}
	spyCache := &spyPostCache{
		getPostsFound: false,
	}

	service := NewPostService(spyRepo, spyCache)
	_, err := service.GetPosts()

	if err == nil {
		t.Error("GetPosts() should return error when DB fails")
	}

	if !spyCache.getPostsCalled {
		t.Error("GetPosts() should call cache.GetPosts")
	}

	if !spyRepo.getPostsCalled {
		t.Error("GetPosts() should call repo.GetPosts when cache misses")
	}

	if spyCache.setPostsCalled {
		t.Error("GetPosts() should NOT call cache.SetPosts when DB fails")
	}
}

func TestUpdatePost_Success(t *testing.T) {
	spyRepo := &spyPostRepo{}
	spyCache := &spyPostCache{}

	service := NewPostService(spyRepo, spyCache)
	err := service.UpdatePost(1, "Updated Title", "Updated Content")

	if err != nil {
		t.Errorf("UpdatePost() error = %v", err)
	}

	if !spyRepo.updatePostCalled {
		t.Error("UpdatePost() should call repo.UpdatePost")
	}

	if !spyCache.invalidatePostCalled {
		t.Error("UpdatePost() should call cache.InvalidatePost after success")
	}

	if spyCache.invalidatePostID != 1 {
		t.Errorf("cache.InvalidatePost() got ID = %v, want %v", spyCache.invalidatePostID, 1)
	}
}

func TestUpdatePost_ValidationFailed(t *testing.T) {
	spyRepo := &spyPostRepo{}
	spyCache := &spyPostCache{}

	service := NewPostService(spyRepo, spyCache)

	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{"empty_title", "", true},
		{"title_256_chars", strings.Repeat("a", 256), true},
		{"valid_title", "Valid Title", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.UpdatePost(1, tt.title, "content")
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePost() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && spyRepo.updatePostCalled {
				t.Error("UpdatePost() should NOT call repo.UpdatePost when validation fails")
			}
		})
	}
}

func TestUpdatePost_NotFound(t *testing.T) {
	spyRepo := &spyPostRepo{
		updatePostErr: db.ErrNoRowsAffected,
	}
	spyCache := &spyPostCache{}

	service := NewPostService(spyRepo, spyCache)
	err := service.UpdatePost(999, "Updated Title", "Updated Content")

	if err != ErrNotFound {
		t.Errorf("UpdatePost() error = %v, want ErrNotFound", err)
	}

	if !spyRepo.updatePostCalled {
		t.Error("UpdatePost() should call repo.UpdatePost")
	}

	if spyCache.invalidatePostCalled {
		t.Error("UpdatePost() should NOT call cache.InvalidatePost when not found")
	}
}

func TestUpdatePost_DBError(t *testing.T) {
	spyRepo := &spyPostRepo{
		updatePostErr: errors.New("DB error"),
	}
	spyCache := &spyPostCache{}

	service := NewPostService(spyRepo, spyCache)
	err := service.UpdatePost(1, "Updated Title", "Updated Content")

	if err == nil {
		t.Error("UpdatePost() should return error when DB fails")
	}

	if !spyRepo.updatePostCalled {
		t.Error("UpdatePost() should call repo.UpdatePost")
	}

	if spyCache.invalidatePostCalled {
		t.Error("UpdatePost() should NOT call cache.InvalidatePost when DB fails")
	}
}

func TestDeletePost_NotFound(t *testing.T) {
	spyRepo := &spyPostRepo{
		deletePostErr: db.ErrNoRowsAffected,
	}
	spyCache := &spyPostCache{}

	service := NewPostService(spyRepo, spyCache)
	err := service.DeletePost(999)

	if err != ErrNotFound {
		t.Errorf("DeletePost() error = %v, want ErrNotFound", err)
	}

	if !spyRepo.deletePostCalled {
		t.Error("DeletePost() should call repo.DeletePost")
	}

	if spyCache.invalidatePostCalled {
		t.Error("DeletePost() should NOT call cache.InvalidatePost when not found")
	}
}

func TestGetPostsByPage_DBError(t *testing.T) {
	spyRepo := &spyPostRepo{
		getPostsByPageErr: errors.New("DB error"),
	}
	spyCache := &spyPostCache{
		getPostsByPageFound: false,
	}

	service := NewPostService(spyRepo, spyCache)
	_, _, err := service.GetPostsByPage(1, 10)

	if err == nil {
		t.Error("GetPostsByPage() should return error when DB fails")
	}

	if !spyRepo.getPostsByPageCalled {
		t.Error("GetPostsByPage() should call repo.GetPostsByPage")
	}
}

func TestGetPostsByPage_CountError(t *testing.T) {
	spyRepo := &spyPostRepo{
		getPostsCountErr: errors.New("count error"),
	}
	spyCache := &spyPostCache{
		getPostsByPageFound: true,
	}

	service := NewPostService(spyRepo, spyCache)
	_, _, err := service.GetPostsByPage(1, 10)

	if err == nil {
		t.Error("GetPostsByPage() should return error when count fails")
	}

	if !spyRepo.getPostsCountCalled {
		t.Error("GetPostsByPage() should call repo.GetPostsCount")
	}
}
