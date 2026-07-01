package model

// Post 文章结构
type Post struct {
	ID      int    `json:"id"`        // 文章ID
	Title   string `json:"title"`     // 文章标题
	Content string `json:"content"`   // 文章内容
	UserID  int    `json:"user_id"`   // 作者ID
}
