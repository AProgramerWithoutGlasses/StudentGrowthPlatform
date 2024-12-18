package routes

import (
	"github.com/gin-gonic/gin"
	"studentGrow/controller/routesJoinAudit"
	"studentGrow/utils/token"
)

func RoutsJoinAudit(router *gin.Engine) {
	r := router.Group("/routesJoinAudit")
	r.Use(token.AuthMiddleware())
	r.POST("/isOpen", routesJoinAudit.OpenMsg)
	r.GET("/StuForm", routesJoinAudit.GetStudForm)
	r.POST("/StuForm", routesJoinAudit.SaveStudForm)
	r.GET("/StudFile", routesJoinAudit.GetStuFile)
	r.POST("/StudFile", routesJoinAudit.SaveStuFile)
	r.POST("/DelStudFile", routesJoinAudit.DelStuFile)
	//r.Use(token.AuthMiddleware(), middleWare.JoinAuditMiddle())
	r.POST("/activity", routesJoinAudit.GetActivityList)
	r.POST("/activityCreat", routesJoinAudit.SaveActivityMsg)
	r.POST("/activityDel", routesJoinAudit.DelActivityMsg)
	r.POST("/activityClass", routesJoinAudit.ClassApplicationList)
	r.POST("/activityClassAudit", routesJoinAudit.ClassApplicationManager)
	r.POST("/activityRuler", routesJoinAudit.ActivityRulerList)
	r.POST("/activityRulerAudit", routesJoinAudit.ActivityRulerManager)
	r.POST("/activityMaterial", routesJoinAudit.ActivityOrganizerList)
	r.POST("/activityMaterialAudit", routesJoinAudit.ActivityOrganizerMaterialManager)
	r.POST("/activityTrain", routesJoinAudit.ActivityOrganizerTrainList)
	r.POST("/activityTrainAudit", routesJoinAudit.ActivityOrganizerTrainManager)
	r.POST("/saveTrainScore", routesJoinAudit.SaveTrainScore)
	r.POST("/exportFormMsg", routesJoinAudit.ExportFormMsg)
	r.POST("/exportFormFile", routesJoinAudit.ExportFormFile)
}
