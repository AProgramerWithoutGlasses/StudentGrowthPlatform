package routes

import (
	"github.com/gin-gonic/gin"
	_ "studentGrow/controller/article"
	"studentGrow/logger"
	"studentGrow/utils/middleWare"
)

func Setup() *gin.Engine {
	r := gin.New()

	r.Use(logger.GinLogger(), logger.GinRecovery(true), middleWare.CORSMiddleware())

	// 星
	routesArticle(r)
	routesTopic(r)
	routesMsg(r)
	routesComment(r)
	routesNotification(r)

	// 勋
	routesHomepage(r)
	routesStudentManage(r)
	routesTeacherManage(r)
	routesClass(r)

	// 雪
	RoutesXue(r)

	RoutsJoinAudit(r)
	//RoutsJoinAuditManager(r)
	return r
}
