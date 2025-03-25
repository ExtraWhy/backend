package main

import (
	"github.com/gin-gonic/gin"
	"login-service/handlers"
)

func main() {
	r := gin.Default()

	r.GET("/auth/google/login", handlers.GoogleLogin)
	r.GET("/auth/google/callback", handlers.GoogleCallback)

	r.Run(":8080")
}

