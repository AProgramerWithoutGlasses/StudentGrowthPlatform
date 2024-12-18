package notification

import (
	"github.com/gin-gonic/gin"
	res "studentGrow/pkg/response"
	"studentGrow/pkg/sse"
)

// SocketConnectionController 向客户端建立sse连接
func SocketConnectionController(c *gin.Context) {
	username := c.Query("username")
	err := sse.BuildNotificationChannel(username, c)
	if err != nil {
		return
	}
	res.ResponseSuccess(c, "断开sse连接")
}
