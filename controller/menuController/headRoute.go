package menuController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/dao/mysql"
	"studentGrow/pkg/response"
	token2 "studentGrow/utils/token"
)

func HeadRoute(c *gin.Context) {
	token := token2.NewToken(c)
	user, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
	}
	role, err := token.GetRole()
	//1.查找用户姓名
	name := user.Name
	//2.查找用户头像
	avatar, err := mysql.SelHead(user.Username)
	//3.查找权限下按钮的所有权限标识
	perms, err := mysql.SelPerms(role)
	if err != nil {
		response.ResponseErrorWithMsg(c, 400, "侧边栏获取信息失败")
		return
	}
	data := map[string]any{
		"name":   name,
		"avatar": avatar,
		"perms":  perms,
	}
	response.ResponseSuccess(c, data)
}
