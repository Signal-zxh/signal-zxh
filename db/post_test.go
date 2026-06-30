package db_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Signal-zxh/signalzxh-blog/db"
	"github.com/Signal-zxh/signalzxh-blog/model"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("创建 mock 失败：%v", err)
	}

	return mockDB, mock
}

type mockPostRepo struct {
	db *sql.DB
}

func (r *mockPostRepo) GetPosts() ([]model.Post, error) {
	rows, err := r.db.Query("SELECT id, title, content, user_id FROM posts ORDER BY id DESC")
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

func (r *mockPostRepo) GetPostsByPage(page, pageSize int) ([]model.Post, error) {
	offset := (page - 1) * pageSize
	rows, err := r.db.Query(
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

func (r *mockPostRepo) GetPostsCount() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM posts").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *mockPostRepo) CreatePost(post model.Post) (int64, error) {
	res, err := r.db.Exec(
		"INSERT INTO posts(title, content, user_id) VALUES(?, ?,?)",
		post.Title, post.Content, post.UserID,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *mockPostRepo) UpdatePost(post model.Post) error {
	res, err := r.db.Exec("UPDATE posts SET title = ?, content = ? WHERE id = ?", post.Title, post.Content, post.ID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return db.ErrNoRowsAffected
	}
	return nil
}

func (r *mockPostRepo) DeletePost(id int) error {
	res, err := r.db.Exec("DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return db.ErrNoRowsAffected
	}
	return nil
}

func (r *mockPostRepo) GetPostByID(id int) (model.Post, error) {
	row := r.db.QueryRow("SELECT id, title, content, user_id FROM posts WHERE id = ?", id)
	var post model.Post
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Post{}, db.ErrNotFound
		}
		return model.Post{}, err
	}
	return post, nil
}

func TestGetPosts_Success(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := &mockPostRepo{db: mockDB}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "user_id"}).
		AddRow(1, "Test Post 1", "Content 1", 1).
		AddRow(2, "Test Post 2", "Content 2", 1)

	mock.ExpectQuery("SELECT id, title, content, user_id FROM posts ORDER BY id DESC").
		WillReturnRows(rows)

	posts, err := repo.GetPosts()
	if err != nil {
		t.Errorf("GetPosts() 错误：%v", err)
	}

	if len(posts) != 2 {
		t.Errorf("GetPosts() 返回 %d 条记录，期望 2 条", len(posts))
	}

	if posts[0].Title != "Test Post 1" {
		t.Errorf("GetPosts()[0].Title = %s，期望 'Test Post 1'", posts[0].Title)
	}

	if posts[1].Title != "Test Post 2" {
		t.Errorf("GetPosts()[1].Title = %s，期望 'Test Post 2'", posts[1].Title)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("期望验证失败：%v", err)
	}
}

func TestGetPosts_Empty(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := &mockPostRepo{db: mockDB}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "user_id"})

	mock.ExpectQuery("SELECT id, title, content, user_id FROM posts ORDER BY id DESC").
		WillReturnRows(rows)

	posts, err := repo.GetPosts()
	if err != nil {
		t.Errorf("GetPosts() 错误：%v", err)
	}

	if len(posts) != 0 {
		t.Errorf("GetPosts() 返回 %d 条记录，期望 0 条", len(posts))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("期望验证失败：%v", err)
	}
}

func TestGetPostsByPage_Success(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := &mockPostRepo{db: mockDB}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "user_id"}).
		AddRow(1, "Test Post 1", "Content 1", 1).
		AddRow(2, "Test Post 2", "Content 2", 1)

	mock.ExpectQuery("SELECT id, title, content, user_id FROM posts ORDER BY id DESC LIMIT \\? OFFSET \\?").
		WithArgs(10, 0).
		WillReturnRows(rows)

	posts, err := repo.GetPostsByPage(1, 10)
	if err != nil {
		t.Errorf("GetPostsByPage() 错误：%v", err)
	}

	if len(posts) != 2 {
		t.Errorf("GetPostsByPage() 返回 %d 条记录，期望 2 条", len(posts))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("期望验证失败：%v", err)
	}
}

func TestGetPostsByPage_SecondPage(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := &mockPostRepo{db: mockDB}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "user_id"}).
		AddRow(3, "Test Post 3", "Content 3", 1)

	mock.ExpectQuery("SELECT id, title, content, user_id FROM posts ORDER BY id DESC LIMIT \\? OFFSET \\?").
		WithArgs(10, 10).
		WillReturnRows(rows)

	posts, err := repo.GetPostsByPage(2, 10)
	if err != nil {
		t.Errorf("GetPostsByPage() 错误：%v", err)
	}

	if len(posts) != 1 {
		t.Errorf("GetPostsByPage() 返回 %d 条记录，期望 1 条", len(posts))
	}

	if posts[0].ID != 3 {
		t.Errorf("GetPostsByPage()[0].ID = %d，期望 3", posts[0].ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("期望验证失败：%v", err)
	}
}

func TestGetPostsByPage_QueryError(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := &mockPostRepo{db: mockDB}

	mock.ExpectQuery("SELECT id, title, content, user_id FROM posts ORDER BY id DESC LIMIT \\? OFFSET \\?").
		WithArgs(10, 0).
		WillReturnError(sql.ErrConnDone)

	_, err := repo.GetPostsByPage(1, 10)
	if err == nil {
		t.Error("GetPostsByPage() 期望返回错误")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("期望验证失败：%v", err)
	}
}

func TestGetPostsCount_Success(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := &mockPostRepo{db: mockDB}

	rows := sqlmock.NewRows([]string{"count"}).AddRow(100)

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM posts").
		WillReturnRows(rows)

	count, err := repo.GetPostsCount()
	if err != nil {
		t.Errorf("GetPostsCount() 错误：%v", err)
	}

	if count != 100 {
		t.Errorf("GetPostsCount() = %d，期望 100", count)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("期望验证失败：%v", err)
	}
}

func TestGetPostsCount_Zero(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := &mockPostRepo{db: mockDB}

	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM posts").
		WillReturnRows(rows)

	count, err := repo.GetPostsCount()
	if err != nil {
		t.Errorf("GetPostsCount() 错误：%v", err)
	}

	if count != 0 {
		t.Errorf("GetPostsCount() = %d，期望 0", count)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("期望验证失败：%v", err)
	}
}

func TestGetPostsCount_QueryError(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := &mockPostRepo{db: mockDB}

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM posts").
		WillReturnError(sql.ErrConnDone)

	_, err := repo.GetPostsCount()
	if err == nil {
		t.Error("GetPostsCount() 期望返回错误")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("期望验证失败：%v", err)
	}
}

func TestCreatePost_Success(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := &mockPostRepo{db: mockDB}

	post := model.Post{
		Title:   "New Post",
		Content: "New Content",
		UserID:  1,
	}

	mock.ExpectExec("INSERT INTO posts").
		WithArgs(post.Title, post.Content, post.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	id, err := repo.CreatePost(post)
	if err != nil {
		t.Errorf("CreatePost() 错误：%v", err)
	}

	if id != 1 {
		t.Errorf("CreatePost() 返回 id = %d，期望 1", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("期望验证失败：%v", err)
	}
}

func TestCreatePost_ExecError(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := &mockPostRepo{db: mockDB}

	post := model.Post{
		Title:   "New Post",
		Content: "New Content",
		UserID:  1,
	}

	mock.ExpectExec("INSERT INTO posts").
		WithArgs(post.Title, post.Content, post.UserID).
		WillReturnError(sql.ErrConnDone)

	_, err := repo.CreatePost(post)
	if err == nil {
		t.Error("CreatePost() 期望返回错误")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("期望验证失败：%v", err)
	}
}

func TestUpdatePost_Success(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := &mockPostRepo{db: mockDB}

	post := model.Post{
		ID:      1,
		Title:   "Updated Title",
		Content: "Updated Content",
	}

	mock.ExpectExec("UPDATE posts SET").
		WithArgs(post.Title, post.Content, post.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdatePost(post)
	if err != nil {
		t.Errorf("UpdatePost() 错误：%v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("期望验证失败：%v", err)
	}
}

func TestUpdatePost_NotFound(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := &mockPostRepo{db: mockDB}

	post := model.Post{
		ID:      999,
		Title:   "Updated Title",
		Content: "Updated Content",
	}

	mock.ExpectExec("UPDATE posts SET").
		WithArgs(post.Title, post.Content, post.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.UpdatePost(post)
	if err != db.ErrNoRowsAffected {
		t.Errorf("UpdatePost() 错误 = %v，期望 ErrNoRowsAffected", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("期望验证失败：%v", err)
	}
}

func TestUpdatePost_ExecError(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := &mockPostRepo{db: mockDB}

	post := model.Post{
		ID:      1,
		Title:   "Updated Title",
		Content: "Updated Content",
	}

	mock.ExpectExec("UPDATE posts SET").
		WithArgs(post.Title, post.Content, post.ID).
		WillReturnError(sql.ErrConnDone)

	err := repo.UpdatePost(post)
	if err == nil {
		t.Error("UpdatePost() 期望返回错误")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("期望验证失败：%v", err)
	}
}

func TestDeletePost_Success(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := &mockPostRepo{db: mockDB}

	mock.ExpectExec("DELETE FROM posts WHERE id = \\?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeletePost(1)
	if err != nil {
		t.Errorf("DeletePost() 错误：%v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("期望验证失败：%v", err)
	}
}

func TestDeletePost_NotFound(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := &mockPostRepo{db: mockDB}

	mock.ExpectExec("DELETE FROM posts WHERE id = \\?").
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.DeletePost(999)
	if err != db.ErrNoRowsAffected {
		t.Errorf("DeletePost() 错误 = %v，期望 ErrNoRowsAffected", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("期望验证失败：%v", err)
	}
}

func TestDeletePost_ExecError(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := &mockPostRepo{db: mockDB}

	mock.ExpectExec("DELETE FROM posts WHERE id = \\?").
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	err := repo.DeletePost(1)
	if err == nil {
		t.Error("DeletePost() 期望返回错误")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("期望验证失败：%v", err)
	}
}

func TestGetPostByID_Success(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := &mockPostRepo{db: mockDB}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "user_id"}).
		AddRow(1, "Test Post", "Test Content", 1)

	mock.ExpectQuery("SELECT id, title, content, user_id FROM posts WHERE id = \\?").
		WithArgs(1).
		WillReturnRows(rows)

	post, err := repo.GetPostByID(1)
	if err != nil {
		t.Errorf("GetPostByID() 错误：%v", err)
	}

	if post.ID != 1 {
		t.Errorf("GetPostByID().ID = %d，期望 1", post.ID)
	}

	if post.Title != "Test Post" {
		t.Errorf("GetPostByID().Title = %s，期望 'Test Post'", post.Title)
	}

	if post.Content != "Test Content" {
		t.Errorf("GetPostByID().Content = %s，期望 'Test Content'", post.Content)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("期望验证失败：%v", err)
	}
}

func TestGetPostByID_NotFound(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := &mockPostRepo{db: mockDB}

	mock.ExpectQuery("SELECT id, title, content, user_id FROM posts WHERE id = \\?").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetPostByID(999)
	if err != db.ErrNotFound {
		t.Errorf("GetPostByID() 错误 = %v，期望 ErrNotFound", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("期望验证失败：%v", err)
	}
}

func TestGetPostByID_QueryError(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := &mockPostRepo{db: mockDB}

	mock.ExpectQuery("SELECT id, title, content, user_id FROM posts WHERE id = \\?").
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	_, err := repo.GetPostByID(1)
	if err == nil {
		t.Error("GetPostByID() 期望返回错误")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("期望验证失败：%v", err)
	}
}
