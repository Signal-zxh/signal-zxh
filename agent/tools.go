package agent

func GetPosts() string {
	return "返回 posts 列表（这里以后查数据库）"
}

func GetPostByID(id string) string {
	return "post详情: " + id
}
