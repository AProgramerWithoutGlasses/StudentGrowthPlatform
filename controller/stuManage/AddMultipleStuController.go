package stuManage

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"studentGrow/dao/mysql"
	"studentGrow/models/gorm_model"
	"studentGrow/pkg/response"
	"studentGrow/utils/studentUtils"
	token2 "studentGrow/utils/token"
	"time"
)

func AddMultipleStuControl(c *gin.Context) {
	token := token2.NewToken(c)
	user, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.ParamFail)
		zap.L().Error("token错误")
	}

	// 获取上传的Excel文件
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		response.ResponseErrorWithMsg(c, 500, "stuManage.AddMultipleStuControl() c.Request.FormFile() failed: "+err.Error())
		zap.L().Error("stuManage.AddMultipleStuControl() c.Request.FormFile() failed: " + err.Error())
		return
	}
	defer file.Close()

	// 解析Excel文件并获取数据
	f, err := excelize.OpenReader(file)
	if err != nil {
		response.ResponseErrorWithMsg(c, 500, "stuManage.AddMultipleStuControl() excelize.OpenReader() failed: "+err.Error())
		zap.L().Error("stuManage.AddMultipleStuControl() excelize.OpenReader() failed: " + err.Error())
		return
	}

	rows := f.GetRows("Sheet1")

	duplicatedUser := make([]string, 0)

	var existedNum int
	var unExistedNum int

	// 检查数据库中是否已经存在该用户
	for _, row := range rows[1:] { // 忽略表头行
		err = mysql.ExistedUsername(row[2])
		if err != nil {
			// 用户不存在
			if errors.Is(err, gorm.ErrRecordNotFound) {
				//err = nil
				//continue

				unExistedNum++

				// 导入

				// 去除所有字段中的空格
				row[0] = studentUtils.RemoveBlank(row[0])
				row[1] = studentUtils.RemoveBlank(row[1])
				row[2] = studentUtils.RemoveBlank(row[2])
				row[3] = studentUtils.RemoveBlank(row[3])
				row[4] = studentUtils.RemoveBlank(row[4])

				// 去除班级中的“班”字
				if len(row[0]) > 9 {
					row[0] = row[0][0:9]
				}

				yearInt, err := strconv.Atoi(row[0][6:8])
				if err != nil {
					response.ResponseError(c, response.ServerErrorCode)
					zap.L().Error(err.Error())
					return
				}

				user1 := gorm_model.User{
					Class:    row[0],
					Name:     row[1],
					Username: row[2],
					Password: row[3],
					Gender:   row[4],
					Identity: "学生",
					PlusTime: time.Date(yearInt+2000, 9, 1, 0, 0, 0, 0, time.Now().Location()),
					HeadShot: "https://student-grow.oss-cn-beijing.aliyuncs.com/image/user_headshot/user_headshot_1.png",
				}
				err = mysql.AddSingleStudent(&user1)
				if err != nil {
					response.ResponseErrorWithMsg(c, 500, "stuManage.AddMultipleStuControl() mysql.AddSingleStudent failed: "+err.Error())
					zap.L().Error("stuManage.AddMultipleStuControl() mysql.AddSingleStudent failed: " + err.Error())
					return
				}

				// 添加学生记录
				addUserRecord := gorm_model.UserAddRecord{
					Username:    user.Username,
					AddUsername: row[2],
				}
				err = mysql.AddSingleStudentRecord(&addUserRecord)
				if err != nil {
					response.ResponseError(c, response.ServerErrorCode)
					zap.L().Error("stuManage.AddMultipleStuControl() mysql.AddSingleStudentRecord() failed: " + err.Error())
					return
				}

			} else {
				response.ResponseErrorWithMsg(c, 500, "stuManage.AddMultipleStuControl() mysql.ExistedUsername() failed: "+err.Error())
				zap.L().Error("stuManage.AddMultipleStuControl() mysql.ExistedUsername() failed: " + err.Error())
				return
			}
		} else {
			// 用户存在
			existedNum++
			duplicatedUser = append(duplicatedUser, row[2])
		}

	}

	//var duplicatedUserStr string
	//if len(duplicatedUser) > 0 {
	//	for _, v := range duplicatedUser {
	//		duplicatedUserStr = duplicatedUserStr + v + ", "
	//	}
	//	duplicatedUserStr = duplicatedUserStr[:len(duplicatedUserStr)-2]
	//	response.ResponseErrorWithMsg(c, 500, "导入失败，请不要导入已存在的学生学号: "+duplicatedUserStr)
	//	zap.L().Error("导入失败，请不要导入已存在的学生学号: " + duplicatedUserStr)
	//	return
	//}

	//// 创建新的用户
	//for _, row := range rows[1:] {
	//	yearInt, err := strconv.Atoi(row[0][6:8])
	//	if err != nil {
	//		response.ResponseError(c, response.ServerErrorCode)
	//		zap.L().Error(err.Error())
	//		return
	//	}
	//
	//	if len(row[0]) > 9 {
	//		row[0] = row[0][0:9]
	//	}
	//	user1 := gorm_model.User{
	//		Class:    row[0],
	//		Name:     row[1],
	//		Username: row[2],
	//		Password: row[3],
	//		Gender:   row[4],
	//		Identity: "学生",
	//		PlusTime: time.Date(yearInt+2000, 9, 1, 0, 0, 0, 0, time.Now().Location()),
	//		HeadShot: "https://student-grow.oss-cn-beijing.aliyuncs.com/image/user_headshot/user_headshot_1.png",
	//	}
	//	err = mysql.AddSingleStudent(&user1)
	//	if err != nil {
	//		response.ResponseErrorWithMsg(c, 500, "stuManage.AddMultipleStuControl() mysql.AddSingleStudent failed: "+err.Error())
	//		zap.L().Error("stuManage.AddMultipleStuControl() mysql.AddSingleStudent failed: " + err.Error())
	//		return
	//	}
	//
	//	// 添加学生记录
	//	addUserRecord := gorm_model.UserAddRecord{
	//		Username:    user.Username,
	//		AddUsername: row[2],
	//	}
	//	err = mysql.AddSingleStudentRecord(&addUserRecord)
	//	if err != nil {
	//		response.ResponseError(c, response.ServerErrorCode)
	//		zap.L().Error("stuManage.AddMultipleStuControl() mysql.AddSingleStudentRecord() failed: " + err.Error())
	//		return
	//	}
	//
	//}

	//addStuNumber := strconv.Itoa(len(rows[1:]))
	var resStr string
	if existedNum == 0 {
		resStr = fmt.Sprintf("已成功导入%d条学生信息", unExistedNum)
	} else {
		resStr = fmt.Sprintf("已成功导入%d条学生信息，剩余%d条未导入，因为其已存在", unExistedNum, existedNum)
	}
	response.ResponseSuccessWithMsg(c, resStr, nil)

}
