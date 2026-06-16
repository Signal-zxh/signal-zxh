package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 静态页面（知识图谱前端放这里）
	r.Static("/static", "./static")

	// 首页
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	r.Run(":8080") // 监听 8080 端口
}
