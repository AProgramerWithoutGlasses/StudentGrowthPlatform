package teacherManage

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"studentGrow/models/gorm_model"
	"studentGrow/pkg/response"
	"studentGrow/service"
)

func AddSingleTeacherControl(c *gin.Context) {
	var addSingleTeacherReqStruct gorm_model.User
	err := c.ShouldBindJSON(&addSingleTeacherReqStruct)
	if err != nil {
		response.ResponseError(c, response.ParamFail)
		zap.L().Error(err.Error())
		return
	}

	err = service.AddSingleTeacherService(addSingleTeacherReqStruct)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			err = errors.New("添加失败, 该用户已存在")
		}
		response.ResponseError(c, response.ServerErrorCode)
		zap.L().Error(err.Error())
		return
	}

	response.ResponseSuccess(c, addSingleTeacherReqStruct.Name+" 添加成功")

}
