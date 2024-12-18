package teacherManage

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/dao/mysql"
	"studentGrow/pkg/response"
)

type DeleteTeacherStruct struct {
	Name            string `json:"name"`
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password"`
	UserGender      string `json:"user_gender"`
	UserManagerType bool   `json:"user_manager_type"`
	UserBan         bool   `json:"user_ban"`
}

// 删除多个老师
func DeleteMultipleTeacherControl(c *gin.Context) {
	// 接收请求
	var input struct {
		SelectedTeachers []DeleteTeacherStruct `json:"selected_teachers"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		response.ResponseErrorWithMsg(c, 500, "teacherManager.DeleteMultipleTeacherControl() c.Bind() failed : "+err.Error())
		zap.L().Error("teacherManager.DeleteMultipleTeacherControl() c.Bind() failed : ", zap.Error(err))
		return
	}

	// 删除选中的用户
	for _, v := range input.SelectedTeachers {
		id, err := mysql.GetIdByUsername(v.Username)
		if err != nil {
			response.ResponseErrorWithMsg(c, 500, "teacherManager.DeleteMultipleTeacherControl() mysql.GetIdByUsername() failed : "+err.Error())
			zap.L().Error("teacherManager.DeleteMultipleTeacherControl() mysql.GetIdByUsername() failed : ", zap.Error(err))
			return
		}

		err = mysql.DeleteSingleUser(id)
		if err != nil {
			response.ResponseErrorWithMsg(c, 500, "teacherManager.DeleteMultipleTeacherControl() mysql.GetIdByUsername() failed : "+err.Error())
			zap.L().Error("teacherManager.DeleteMultipleTeacherControl() mysql.DeleteSingleStudent() err : ", zap.Error(err))
			return
		}
	}

	// 响应成功
	response.ResponseSuccess(c, "删除成功!")
}
