package router

import (
	"github.com/gin-gonic/gin"
	"meeting_demo/internal/server/service"
	"meeting_demo/middlewares"
)

func Router() *gin.Engine {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	// cors
	router.Use(middlewares.Cors())

	// user login
	router.POST("/user/login", service.UserLogin)

	// websocket
	router.GET("/ws/p2p/:room_id/:user_id", service.WsP2PConnection)

	auth := router.Group("/auth", middlewares.Auth())

	// get meeting list
	auth.GET("/meeting/list", service.MeetingList)

	// create meeting
	auth.POST("/meeting/create", service.MeetingCreate)

	// edit meeting
	auth.PUT("/meeting/edit", service.MeetingEdit)

	// delete meeting
	auth.DELETE("/meeting/delete", service.MeetingDelete)

	return router
}
