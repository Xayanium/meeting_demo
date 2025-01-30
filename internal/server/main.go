package main

import (
	"meeting_demo/internal/models"
	"meeting_demo/internal/server/router"
)

func main() {
	models.NewDB()
	engine := router.Router()

	_ = engine.Run("127.0.0.1:8080")
}
