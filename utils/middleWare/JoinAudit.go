package middleWare

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"studentGrow/service/JoinAudit"
	"studentGrow/service/userService"
	token2 "studentGrow/utils/token"
)

// 入团申请权限判断中间件
func JoinAuditMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("入团中间件")
		token := token2.NewToken(c)
		user, _ := token.GetUser()
		isAdmin := userService.BVerifyExit(user.Username)
		isJoinAudit := JoinAudit.GetUserJoinAuditRoel(c)
		fmt.Println(isAdmin, isJoinAudit)
		if !isJoinAudit && !isAdmin {
			fmt.Println("驳回")
			c.JSON(500, gin.H{"code": 500, "msg": "驳回"})
			c.Abort()
			return
		}
		fmt.Println("继续")
		c.Next()
	}
}
