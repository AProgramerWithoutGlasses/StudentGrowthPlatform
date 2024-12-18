package login

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/dao/mysql"
	"studentGrow/models"
	pkg "studentGrow/pkg/response"
	"studentGrow/service/userService"
	"studentGrow/utils/timeConverter"
	token2 "studentGrow/utils/token"
	"time"
)

// RidCode 图形验证码返回前端
func RidCode(c *gin.Context) {
	//获取验证码
	id, b64, hcode, _ := userService.CaptchaGenerate()
	if id == "" || b64 == " " {
		pkg.ResponseError(c, 400)
		fmt.Println("RidCode  Login() err")
		return
	}
	code := map[string]any{
		"Id":    id,
		"Hcode": hcode,
		"B64":   b64,
	}
	pkg.ResponseSuccess(c, code)
}

// HLogin 登录验证，成功返回token
func HLogin(c *gin.Context) {
	// 定义用户实例
	var user = new(models.Login)
	//获取前端返回的数据
	if err := c.BindJSON(&user); err != nil {
		pkg.ResponseErrorWithMsg(c, 400, "未获取到数据")
		fmt.Println("Hlogin BindJSON(&userService) err")
		return
	}
	//查询用户是否存在
	ok, err := mysql.SelExit(user.Username)
	if !ok || err != nil {
		pkg.ResponseErrorWithMsg(c, 500, "用户不存在")
		return
	}
	//验证密码
	if ok := userService.BVerify(user.Username, user.Password); !ok {
		pkg.ResponseErrorWithMsg(c, 400, "密码错误")
		return
	}
	//验证验证码
	if ok := userService.GetCodeAnswer(user.Id, user.Code); !ok {
		pkg.ResponseErrorWithMsg(c, 401, "验证码错误")
		return
	}
	//验证用户是否为管理员
	if ok := userService.BVerifyExit(user.Username); !ok {
		pkg.ResponseErrorWithMsg(c, 501, "身份验证失败")
		return
	}
	//验证用户是否被封禁
	ban, err := userService.BVerifyBan(user.Username)
	if err != nil {
		pkg.ResponseError(c, 400)
		return
	}
	if ban {
		data := time.Now().Format("2006-01-02")
		ok, err := userService.UpdateStatus(data, user.Username)
		if err != nil {
			pkg.ResponseError(c, 400)
			return
		}
		if !ok {
			pkg.ResponseErrorWithMsg(c, 400, "账号被封禁")
			return
		} else {
			//解禁
			err := mysql.UpdateBan(user.Username)
			if err != nil {
				pkg.ResponseError(c, 400)
				return
			}
		}
	}
	//查询整个user
	thisUser, err := mysql.SelUser(user.Username)
	//查询用户角色
	casbinId, _ := mysql.SelCasId(user.Username)
	if casbinId == "" {
		pkg.ResponseError(c, 400)
		return
	}
	role, err := mysql.SelRole(casbinId)
	//获取生成的token
	tokenString, err := token2.ReleaseToken(thisUser)
	if err != nil {
		fmt.Println("Hlogin的login.ReleaseToken()")
		return
	}
	if err != nil {
		fmt.Println("HLogin SelId err")
		return
	}
	slice := map[string]any{
		"username": user.Username,
		"token":    tokenString,
		"role":     role,
	}
	//发送给前端
	pkg.ResponseSuccess(c, slice)
}

// QLogin 登录验证，成功返回token
func QLogin(c *gin.Context) {
	// 定义用户实例
	var user = new(models.Login)
	var role string
	//记录班级
	var class string
	//记录年级
	var grade int
	//获取前端返回的数据
	if err := c.BindJSON(&user); err != nil {
		pkg.ResponseErrorWithMsg(c, 400, "Hlogin 获取数据失败")
		fmt.Println("Hlogin BindJSON(&userService) err")
		return
	}
	//查询用户是否存在
	ok, err := mysql.SelExit(user.Username)
	if !ok || err != nil {
		pkg.ResponseErrorWithMsg(c, 500, "用户不存在")
		return
	}
	//验证密码
	if ok := userService.BVerify(user.Username, user.Password); !ok {
		pkg.ResponseErrorWithMsg(c, 400, "密码错误")
		return
	}
	//验证验证码
	if ok := userService.GetCodeAnswer(user.Id, user.Code); !ok {
		pkg.ResponseErrorWithMsg(c, 401, "验证码错误")
		return
	}
	//验证用户是否被封禁
	ban, err := userService.BVerifyBan(user.Username)
	if err != nil {
		pkg.ResponseError(c, 400)
		return
	}
	if ban {
		data := time.Now().Format("2006-01-02")
		ok, err := userService.UpdateStatus(data, user.Username)
		if err != nil {
			pkg.ResponseError(c, 400)
			return
		}
		if !ok {
			pkg.ResponseErrorWithMsg(c, 400, "账号被封禁")
			return
		} else {
			//解禁
			err := mysql.UpdateBan(user.Username)
			if err != nil {
				pkg.ResponseError(c, 400)
				return
			}
		}
	}
	//验证用户是否是管理员
	newOk := userService.BVerifyExit(user.Username)
	if newOk {
		cId, err := mysql.SelCasId(user.Username)
		role, err = mysql.SelRole(cId)
		if err != nil {
			pkg.ResponseError(c, 400)
			return
		}
	} else {
		role = "user"
	}
	//查询整个user
	thisUser, err := mysql.SelUser(user.Username)
	//验证用户是否是老师
	ifTeacher, err := mysql.IfTeacher(user.Username)
	class, err = mysql.SelClass(user.Username)
	plusTime, err := mysql.SelPlus(user.Username)
	if !plusTime.Valid {
		grade = 0
	} else {
		grade = timeConverter.GetUserGrade(plusTime.Time)
	}
	if err != nil {
		fmt.Println("Hlogin的login.ReleaseToken()")
		return
	}
	//记录用户登录
	id, err := mysql.SelId(user.Username)
	err = mysql.CreateUser(user.Username, id)
	//获取生成的token
	tokenString, err := token2.ReleaseToken(thisUser)
	if err != nil {
		fmt.Println("Hlogin的login.ReleaseToken()")
		return
	}
	if err != nil {
		fmt.Println("HLogin SelId err")
		return
	}
	//前端转化class为数组
	str := []string{
		class,
	}
	slice := map[string]any{
		"username":  user.Username,
		"token":     tokenString,
		"role":      role,
		"class":     str,
		"grade":     grade,
		"ifTeacher": ifTeacher,
	}
	//发送给前端
	pkg.ResponseSuccess(c, slice)
}

// RegisterDay 前台获取注册天数
func RegisterDay(c *gin.Context) {
	var plus_time int
	token := token2.NewToken(c)
	user, exist := token.GetUser()
	if !exist {
		pkg.ResponseError(c, pkg.TokenError)
		zap.L().Error("token错误")
		return
	}
	plusTime, err := mysql.SelPlus(user.Username)
	if err != nil {
		pkg.ResponseError(c, 400)
		return
	}
	if !plusTime.Valid {
		plus_time = 0
	} else {
		plus_time = userService.IntervalInDays(plusTime.Time)
	}
	data := map[string]any{
		"plus_time": plus_time,
	}
	pkg.ResponseSuccess(c, data)
}
