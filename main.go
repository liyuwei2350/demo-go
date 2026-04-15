package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// /auth 路由供 Nginx auth_request 使用
	r.GET("/auth", func(c *gin.Context) {

		// --- 灰度标识逻辑 ---
		// 根据用户信息、特定 Header 或其他条件来决定是否开启灰度
		// 例如：如果客户端请求带了特定的灰度测试 Header X-Test-Gray，则判定为灰度用户
		isGrayUser := true

		if isGrayUser {

			grayHeaderValue := "gray"
			c.Header("X-Gray-Env", grayHeaderValue)

		}

		// 鉴权成功，必须返回 200 或其他 2xx 状态码，Nginx 才会放行请求
		c.Status(http.StatusOK)
	})

	// 启动微服务，监听 8080 端口
	log.Println("Nginx Auth Service starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
