package routes

import (
	"github.com/gin-gonic/gin"
	"studentGrow/controller/notification"
)

func routesNotification(r *gin.Engine) {
	rN := r.Group("/notification")

	rN.GET("/socket_connection", notification.SocketConnectionController)
}
