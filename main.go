package main

import (
	"Tikky/db"
	"Tikky/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDatabase()

	r := gin.Default()
	r.Use(handler.AuthMiddleware)

	initRouter(r)

	err := r.Run()
	if err != nil {
		panic("Failed to start server")
	}
}
