package middleWare

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"studentGrow/dao/mysql"
	"studentGrow/models/casbinModels"
	pkg "studentGrow/pkg/response"
	token2 "studentGrow/utils/token"
)

func NewCasbinAuth(srv *casbinModels.CasbinService) gin.HandlerFunc {
	return func(c *gin.Context) {
		//加载策略文件
		err := srv.Enforcer.LoadPolicy()
		if err != nil {
			pkg.ResponseError(c, 400)
			c.Abort()
			return
		}
		//拿到角色
		token := token2.NewToken(c)
		role, err := token.GetRole()
		if err != nil || role == "" {
			fmt.Println("NewCasbinAuth myToken.GetRole(token) err:")
			return
		}
		menuId, err := mysql.SelMenuId(c.Request.URL.Path, c.Request.Method)
		if err != nil || menuId == "" {
			fmt.Println("NewCasbinAuth() mysql.SelMenuId err", err)
			pkg.ResponseErrorWithMsg(c, 500, "没有找到菜单ID")
			c.Abort()
			return
		}
		fmt.Println(menuId, err)
		ok, err := srv.Enforcer.Enforce(role, menuId)
		if err != nil {
			c.JSON(400, gin.H{"code": 400, "msg": "出错"})
			c.Abort()
			return
		} else if !ok {
			c.JSON(500, gin.H{"code": 500, "msg": "驳回"})
			c.Abort()
			return
		}
		c.Next()
	}
}
