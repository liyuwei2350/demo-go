package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	bookStore := newBookStore()

	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// /auth is used by Nginx auth_request.
	r.GET("/auth", func(c *gin.Context) {
		isGrayUser := true
		if isGrayUser {
			c.Header("X-Gray-Env", "gray")
		}

		c.Status(http.StatusOK)
	})

	registerBookRoutes(r, bookStore)

	log.Println("Nginx Auth Service starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
