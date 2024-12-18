package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/controller/stuManage"
	"studentGrow/dao/mysql"
	"studentGrow/models/casbinModels"
	"studentGrow/utils/middleWare"
	"studentGrow/utils/token"
)

// 后台学生管理路由
func routesStudentManage(r *gin.Engine) {
	casbinService, err := casbinModels.NewCasbinService(mysql.DB)
	if err != nil {
		zap.L().Error("routesArticle() routes.routesArticle.NewCasbinService err=", zap.Error(err))
		return
	}
	rs := r.Group("/stuManage")
	rs.Use(token.AuthMiddleware(), middleWare.NewCasbinAuth(casbinService))

	{
		// rs.POST("/queryPageStudent", stuManage.QueryPageStuContro)
		rs.POST("/addSingleStudent", stuManage.AddSingleStuContro)
		rs.POST("/addMultipleStudent", stuManage.AddMultipleStuControl)

		rs.POST("/queryStudent", stuManage.QueryStuContro) // 学号、角色（7种）

		rs.POST("/deleteStudent", stuManage.DeleteStuControl)
		rs.POST("/setStudentManager", stuManage.SetStuManagerControl)
		rs.POST("/editStudent", stuManage.EditStuControl)
		rs.POST("/banStudent", stuManage.BanUserControl)
		rs.POST("/outputMultipleStudent", stuManage.OutputMultipleStuControl)
		rs.POST("/reSetPassword", stuManage.ReSetPasswordControl)
	}

}
