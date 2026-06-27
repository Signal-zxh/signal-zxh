package agent

import (
	"strings"
)

func RouteTool(query string) string {
	if strings.Contains(query, "posts") {
		return GetPosts()
	}

	if strings.Contains(query, "post") {
		return GetPostByID("26")
	}
	return "没找到工具"
}
