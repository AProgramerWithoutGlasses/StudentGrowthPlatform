package teacherManage

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/models/jrx_model"
	"studentGrow/pkg/response"
	"studentGrow/service"
)

type ResponseStruct struct {
	TeacherInfo     []jrx_model.QueryTeacherResStruct `json:"teacherInfo"`
	AllTeacherCount int                               `json:"allTeacherCount"`
}

func QueryTeacherControl(c *gin.Context) {
	// 接收
	var queryParama jrx_model.QueryTeacherParamStruct
	err := c.ShouldBindJSON(&queryParama)
	if err != nil {
		response.ResponseError(c, response.ParamFail)
		zap.L().Error("teacherManage.QueryTeacher() c.Bind() err : ", zap.Error(err))
		return
	}

	// 业务
	teacherResList, allTeacherCount, err := service.QueryTeacher(queryParama)
	if err != nil {
		response.ResponseError(c, response.ServerErrorCode)
		zap.L().Error("teacherManage.QueryTeacher() service.QueryTeacher() err : ", zap.Error(err))
		return
	}

	// 响应
	responseStruct := ResponseStruct{
		TeacherInfo:     teacherResList,
		AllTeacherCount: allTeacherCount,
	}

	response.ResponseSuccess(c, responseStruct)

}
