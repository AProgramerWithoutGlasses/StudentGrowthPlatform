package routes

import (
	"github.com/gin-gonic/gin"
	"studentGrow/controller/teacherManage"
)

// 后台老师管理路由
func routesTeacherManage(r *gin.Engine) {
	rt := r.Group("/teacherManage")

	rt.POST("/queryTeacher", teacherManage.QueryTeacherControl)
	rt.POST("/addSingleTeacher", teacherManage.AddSingleTeacherControl)
	rt.POST("/addMultipleTeacher", teacherManage.AddMultipleTeacherControl)
	rt.POST("/deleteTeacher", teacherManage.DeleteTeacherControl)
	rt.POST("/deleteMultipleTeacher", teacherManage.DeleteMultipleTeacherControl)
	rt.POST("/setTeacherManager", teacherManage.SetTeacherManagerControl)
	rt.POST("/editTeacher", teacherManage.EditTeacherControl)
	rt.POST("/banTeacher", teacherManage.BanTeacherControl)
	rt.POST("/outputMultipleTeacher", teacherManage.OutputMultipleTeacherControl)
	rt.POST("/getUsername", teacherManage.GetUsername)
}
