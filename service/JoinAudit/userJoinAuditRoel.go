package JoinAudit

import (
	"github.com/gin-gonic/gin"
	token2 "studentGrow/utils/token"
)

// 根据身份是否为空判断其是否有入团审核访问权限
func GetUserJoinAuditRoel(c *gin.Context) bool {
	token := token2.NewToken(c)
	user, _ := token.GetUser()
	if user.JobClass == "班长" && user.JobStuUnion == "学生会" {
		return true
	} else {
		return false
	}
}
