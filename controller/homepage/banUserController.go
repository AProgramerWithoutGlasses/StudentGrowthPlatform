package homepage

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/pkg/response"
	"studentGrow/service"
	token2 "studentGrow/utils/token"
)

func BanUserControl(c *gin.Context) {
	// 接收数据
	input := struct {
		BanTime     int    `json:"ban_time" binding:"required"`
		BanUsername string `json:"ban_username" binding:"required"`
	}{}
	err := c.BindJSON(&input)
	if err != nil {
		fmt.Println("stuManage.BanStuControl() c.Bind() err : ", err)
		response.ResponseErrorWithMsg(c, 500, "stuManage.BanStuControl() c.Bind() failed : "+err.Error())
		return
	}
	//
	//token := c.GetHeader("token")
	//username, err := token2.GetUsername(token)
	//if err != nil {
	//	response.ResponseError(c, response.ParamFail)
	//	zap.L().Error(err.Error())
	//	return
	//}

	token := token2.NewToken(c)
	user, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.ParamFail)
		zap.L().Error("token错误")
	}

	err = service.BanHomepageUserService(input.BanUsername, input.BanTime, user.Username)
	if err != nil {
		response.ResponseErrorWithMsg(c, 500, "stuManage.BanStuControl() service.BanUserService() failed : "+err.Error())
		return
	}

	response.ResponseSuccess(c, struct{}{})
}
