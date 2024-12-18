package teacherManage

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
	"studentGrow/dao/mysql"
	"studentGrow/models/gorm_model"
	"studentGrow/pkg/response"
)

func AddMultipleTeacherControl(c *gin.Context) {
	// 获取上传的Excel文件
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		response.ResponseErrorWithMsg(c, 500, "teacherManage.AddMultipleTeacherControl() c.Request.FormFile() failed: "+err.Error())
		zap.L().Error("teacherManage.AddMultipleTeacherControl() c.Request.FormFile() failed: " + err.Error())
		return
	}
	defer file.Close()

	// 解析Excel文件并获取数据
	f, err := excelize.OpenReader(file)
	if err != nil {
		response.ResponseErrorWithMsg(c, 500, "teacherManage.AddMultipleTeacherControl() excelize.OpenReader() failed: "+err.Error())
		zap.L().Error("teacherManage.AddMultipleTeacherControl() excelize.OpenReader() failed: " + err.Error())
		return
	}

	rows := f.GetRows("Sheet1")

	duplicatedUser := make([]string, 0)

	// 检查数据库中是否已经存在该用户
	for _, row := range rows[1:] { // 忽略表头行
		err = mysql.ExistedUsername(row[1])
		if err != nil {
			if err == gorm.ErrRecordNotFound {

			} else {
				response.ResponseErrorWithMsg(c, 500, "teacherManage.AddMultipleTeacherControl() mysql.ExistedUsername() failed: "+err.Error())
				zap.L().Error("teacherManage.AddMultipleTeacherControl() mysql.ExistedUsername() failed: " + err.Error())
				return
			}

		} else { // 用户存在
			duplicatedUser = append(duplicatedUser, row[2])
		}
	}

	var duplicatedUserStr string
	if len(duplicatedUser) > 0 {
		for _, v := range duplicatedUser {
			duplicatedUserStr = duplicatedUserStr + v + ", "
		}
		duplicatedUserStr = duplicatedUserStr[:len(duplicatedUserStr)-2]
		response.ResponseErrorWithMsg(c, 500, "导入失败，请不要导入已存在的老师账户: "+duplicatedUserStr)
		zap.L().Error("导入失败，请不要导入已存在的老师账户: " + duplicatedUserStr)
		return
	}

	// 创建新的用户
	for _, row := range rows[1:] {
		user := gorm_model.User{
			Name:     row[0],
			Username: row[1],
			Password: row[2],
			Gender:   row[3],
			Identity: "老师",
			HeadShot: "https://student-grow.oss-cn-beijing.aliyuncs.com/image/user_headshot/user_headshot_5.png",
		}
		err = mysql.AddSingleTeacher(&user)
		if err != nil {
			if strings.Contains(err.Error(), "Duplicate entry") {
				err = nil
			} else {
				response.ResponseErrorWithMsg(c, 500, "teacherManage.AddMultipleTeacherControl() mysql.AddSingleStudent failed: "+err.Error())
				zap.L().Error("teacherManage.AddMultipleTeacherControl() mysql.AddSingleStudent failed: " + err.Error())
				return
			}

		}
	}

	response.ResponseSuccess(c, "导入成功!")

}
