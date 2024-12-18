package growth

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/dao/mysql"
	"studentGrow/models"
	"studentGrow/models/constant"
	"studentGrow/pkg/response"
	"studentGrow/service/starService"
	token2 "studentGrow/utils/token"
	"time"
)

// Student 定义接收前端数据结构体
type Student struct {
	Username         string `json:"username"`
	Name             string `json:"name"`
	Userarticletotal int    `json:"user_article_total"`
	Userfans         int    `json:"userfans"`
	Score            int    `json:"score"`
	Hot              int    `json:"hot"`
	Frequency        int    `json:"frequency"`
}

// Search 搜索表格数据
func Search(c *gin.Context) {
	//返回前端限制人数
	var peopleLimit int
	var usernamesli []string
	var total int64
	//获取前端传来的数据
	var datas struct {
		Name  string `form:"search"`
		Page  int    `form:"page"`
		Limit int    `form:"pageCapacity"`
	}
	err := c.Bind(&datas)
	if err != nil {
		zap.L().Error("Search Bind err", zap.Error(err))
		response.ResponseError(c, response.ServerErrorCode)
		return
	}

	//得到登录者的角色和账号
	token := token2.NewToken(c)
	user, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
	}
	role, err := token.GetRole()
	if err != nil {
		zap.L().Error("Search err", zap.Error(err))
		response.ResponseError(c, response.ServerErrorCode)
		return
	}

	//根据角色分类
	switch role {
	case "class":
		if datas.Name == "" {
			usernamesli, total, err = mysql.SelUsername(user.Class, datas.Page, datas.Limit)
		} else {
			usernamesli, total, err = mysql.SelSearchUser(datas.Name, user.Class, datas.Page, datas.Limit)
		}
		peopleLimit = constant.PeopleLimitClass
		if err != nil {
			zap.L().Error("Search SelSearchUser err", zap.Error(err))
			response.ResponseError(c, response.ServerErrorCode)
			return
		}
	case "grade1":
		if datas.Name == "" {
			usernamesli, total, err = starService.StarGuidGrade(-1, datas.Page, datas.Limit)
		} else {
			usernamesli, total, err = starService.SearchGrade(datas.Name, -1, datas.Page, datas.Limit)
		}
		peopleLimit = constant.PeopleLimitGrade
		if err != nil {
			zap.L().Error("Search GetEnrollmentYear err", zap.Error(err))
			response.ResponseError(c, response.ServerErrorCode)
			return
		}
	case "grade2":
		if datas.Name == "" {
			usernamesli, total, err = starService.StarGuidGrade(-2, datas.Page, datas.Limit)
		} else {
			usernamesli, total, err = starService.SearchGrade(datas.Name, -2, datas.Page, datas.Limit)
		}
		peopleLimit = constant.PeopleLimitGrade
		if err != nil {
			zap.L().Error("Search GetEnrollmentYear err", zap.Error(err))
			response.ResponseError(c, response.ServerErrorCode)
			return
		}
	case "grade3":
		if datas.Name == "" {
			usernamesli, total, err = starService.StarGuidGrade(-3, datas.Page, datas.Limit)
		} else {
			usernamesli, total, err = starService.SearchGrade(datas.Name, -3, datas.Page, datas.Limit)
		}
		peopleLimit = constant.PeopleLimitGrade
		if err != nil {
			zap.L().Error("Search GetEnrollmentYear err", zap.Error(err))
			response.ResponseError(c, response.ServerErrorCode)
			return
		}
	case "grade4":
		if datas.Name == "" {
			usernamesli, total, err = starService.StarGuidGrade(-4, datas.Page, datas.Limit)
		} else {
			usernamesli, total, err = starService.SearchGrade(datas.Name, -4, datas.Page, datas.Limit)
		}
		peopleLimit = constant.PeopleLimitGrade
		if err != nil {
			zap.L().Error("Search GetEnrollmentYear err", zap.Error(err))
			response.ResponseError(c, response.ServerErrorCode)
			return
		}
	case "college":
		if datas.Name == "" {
			usernamesli, total, err = mysql.SelStarColl(datas.Page, datas.Limit)
		} else {
			usernamesli, total, err = mysql.SelSearchColl(datas.Name, datas.Page, datas.Limit)
		}
		peopleLimit = constant.PeopleLimitCollege
		if err != nil {
			zap.L().Error("Search SelSearchColl err", zap.Error(err))
			response.ResponseError(c, response.ServerErrorCode)
			return
		}
	case "superman":
		if datas.Name == "" {
			usernamesli, total, err = mysql.SelStarColl(datas.Page, datas.Limit)
		} else {
			usernamesli, total, err = mysql.SelSearchColl(datas.Name, datas.Page, datas.Limit)
		}
		peopleLimit = constant.PeopleLimitCollege
		if err != nil {
			zap.L().Error("Search SelSearchColl err", zap.Error(err))
			response.ResponseError(c, response.ServerErrorCode)
			return
		}
	}
	//表格所需所有数据
	starback, err := starService.StarGrid(usernamesli)
	if err != nil {
		zap.L().Error("Search SelSearchColl err", zap.Error(err))
		response.ResponseError(c, response.ServerErrorCode)
		return
	}

	//查询状态
	status, err := mysql.SelStatus(user.Username)
	if err != nil {
		response.ResponseError(c, 400)
		return
	}
	headline, err := mysql.SelMax()
	data := map[string]any{
		"tableData":   starback,
		"total":       int(total),
		"peopleLimit": peopleLimit,
		"isDisabled":  status,
		"headline":    headline + 1,
	}
	response.ResponseSuccess(c, data)
}

// ElectClass  班级管理员推选数据
func ElectClass(c *gin.Context) {
	//获取管理员班级
	token := token2.NewToken(c)
	user, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
	}
	//查找表中存在几条已推选的数据
	Number, err := starService.SelNumClass(user.Class)
	if err != nil {
		response.ResponseError(c, 400)
		return
	}
	//接收前端数据
	var Responsedata struct {
		ElectedArr []Student `json:"electedArr"`
	}
	err = c.Bind(&Responsedata)
	if err != nil {
		zap.L().Error("Search Bind err", zap.Error(err))
		response.ResponseErrorWithMsg(c, response.ServerErrorCode, "未获取到数据")
		return
	}

	//判断是否超出权限范围并判断是否有重复数据
	length := len(Responsedata.ElectedArr)
	if Star := length + Number; Star > constant.PeopleLimitClass {
		response.ResponseErrorWithMsg(c, 200, "本次推选超出名额限制")
		return
	}
	for _, student := range Responsedata.ElectedArr {
		username := student.Username
		name := student.Name
		//防止有重复数据
		number, err := mysql.Selstarexit(username)
		if err != nil || number != 0 {
			response.ResponseErrorWithMsg(c, 200, "该人员已被推选")
			return
		}

		//添加数据
		err = mysql.CreatClass(username, name)
		if err != nil {
			response.ResponseErrorWithMsg(c, 400, "推选失败")
			return
		}
	}
	if err != nil {
		response.ResponseError(c, 400)
		return
	}
	//判断是否超出权限范围

	if Star := length + Number; Star == constant.PeopleLimitClass {
		data := map[string]any{
			"news": "No seats left",
		}
		response.ResponseSuccess(c, data)
		return
	}
	data := map[string]any{
		"news": "推选成功",
	}
	response.ResponseSuccess(c, data)
}

// ElectGrade 年级管理员推选数据
func ElectGrade(c *gin.Context) {
	//代表年级管理员已推选的个数
	var number int
	var err error
	date := time.Now()
	//拿到角色
	token := token2.NewToken(c)
	role, err := token.GetRole()
	//获取number的值
	switch role {
	case "grade1":
		number, err = starService.SelNumGrade(date, -1)
	case "grade2":
		number, err = starService.SelNumGrade(date, -2)
	case "grade3":
		number, err = starService.SelNumGrade(date, -3)
	case "grade4":
		number, err = starService.SelNumGrade(date, -4)
	}
	if err != nil {
		response.ResponseError(c, 400)
		return
	}
	//获取前端传来的数据
	var ResponseData struct {
		ElectedArr []Student `json:"electedArr"`
	}
	err = c.Bind(&ResponseData)
	if err != nil {
		zap.L().Error("Search Bind err", zap.Error(err))
		response.ResponseErrorWithMsg(c, response.ServerErrorCode, "未获取到数据")
		return
	}

	//查询这次需要推选的人
	length := len(ResponseData.ElectedArr)
	if star := number + length; star > constant.PeopleLimitGrade {
		response.ResponseErrorWithMsg(c, 200, "本次推选超出名额限制")
		return
	}

	//开始添加
	for _, user := range ResponseData.ElectedArr {
		username := user.Username
		//更新数据
		err = mysql.UpdateGrade(username)
		if err != nil {
			response.ResponseErrorWithMsg(c, 400, "推选失败")
			return
		}
	}

	//刚好推选够5个时
	if Star := length + number; Star == constant.PeopleLimitGrade {
		data := map[string]any{
			"news": "No seats left",
		}
		response.ResponseSuccess(c, data)
		return
	}
	data := map[string]any{
		"news": "推选成功",
	}
	response.ResponseSuccess(c, data)
}

// ElectCollege 院级管理员推选
func ElectCollege(c *gin.Context) {
	var Responsedata struct {
		ElectedArr []Student `json:"electedArr"`
	}
	err := c.Bind(&Responsedata)
	if err != nil {
		zap.L().Error("Search Bind err", zap.Error(err))
		response.ResponseErrorWithMsg(c, response.ServerErrorCode, "未获取到数据")
		return
	}
	//已推选的个数
	BefName, err := mysql.SelDataCollege()
	if len(BefName)+len(Responsedata.ElectedArr) > constant.PeopleLimitCollege {
		response.ResponseErrorWithMsg(c, 200, "超出名额限制")
		return
	}
	for _, user := range Responsedata.ElectedArr {
		username := user.Username
		err = mysql.UpdateCollege(username)
		if err != nil {
			response.ResponseErrorWithMsg(c, 400, "推选失败")
			return
		}
	}
	if len(BefName)+len(Responsedata.ElectedArr) == constant.PeopleLimitCollege {
		data := map[string]any{
			"news": "No seats left",
		}
		response.ResponseSuccess(c, data)
		return
	}
	data := map[string]any{
		"news": "推选成功",
	}
	response.ResponseSuccess(c, data)
}

// PublicStar 公布成长之星
func PublicStar(c *gin.Context) {
	//获取session字段的最大值
	nowSession, err := mysql.SelMax()
	if err != nil {
		response.ResponseErrorWithMsg(c, 400, "获取值失败")
		return
	}
	session := nowSession + 1
	//更新字段
	err = mysql.UpdateSession(session)
	if err != nil {
		response.ResponseErrorWithMsg(c, 400, "公布失败")
		return
	}

	//更新所有管理员状态字段
	err = mysql.UpdateStatus()
	if err != nil {
		response.ResponseError(c, 400)
		return
	}

	//展示最新一期
	session, err = mysql.SelMax()
	if err != nil {
		response.ResponseErrorWithMsg(c, 400, "获取最新数据失败")
		return
	}

	//返回数据
	//班级成长之星
	classData, err := starService.StarClass(session)
	if err != nil {
		response.ResponseErrorWithMsg(c, 400, "班级之星查找失败")
		return
	}

	//年级之星
	gradeData, err := starService.StarGrade(session)
	if err != nil {
		response.ResponseErrorWithMsg(c, 400, "年级之星查找失败")
		return
	}

	//院级之星
	hospitalData, err := starService.StarCollege(session)
	if err != nil {
		response.ResponseErrorWithMsg(c, 400, "院级之星查找失败")
		return
	}

	data := map[string]any{
		"classData":    classData,
		"gradeData":    gradeData,
		"hospitalData": hospitalData,
	}
	response.ResponseSuccess(c, data)
}

// StarPub 搜索第几届成长之星
func StarPub(c *gin.Context) {
	//定义届数
	var session int
	//接收前端数据
	var term struct {
		TermNumber int `form:"termNumber"`
	}
	err := c.Bind(&term)
	if err != nil {
		zap.L().Error("Search Bind err", zap.Error(err))
		response.ResponseError(c, response.ServerErrorCode)
		return
	}
	//设置session的值
	if term.TermNumber == 0 {
		//找到最大的session最新一期进行展示
		session, err = mysql.SelMax()
		if session == 0 {
			response.ResponseSuccess(c, "")
			return
		}
		if err != nil {
			response.ResponseErrorWithMsg(c, 400, "获取最新数据失败")
			return
		}
	} else {
		session = term.TermNumber
	}
	//班级成长之星
	classData, err := starService.StarClass(session)
	if err != nil {
		response.ResponseErrorWithMsg(c, 400, "班级之星查找失败")
		return
	}

	//年级之星
	gradeData, err := starService.StarGrade(session)
	if err != nil {
		response.ResponseErrorWithMsg(c, 400, "年级之星查找失败")
		return
	}

	//院级之星
	hospitalData, err := starService.StarCollege(session)
	if err != nil {
		response.ResponseErrorWithMsg(c, 400, "院级之星查找失败")
		return
	}

	data := map[string]any{
		"classData":    classData,
		"gradeData":    gradeData,
		"hospitalData": hospitalData,
	}
	response.ResponseSuccess(c, data)
}

// BackStarClass 返回前台班级成长之星
func BackStarClass(c *gin.Context) {
	//返回前端数据
	var starList []models.StarStu
	//接收前端数据
	var backData struct {
		StartTime string `form:"startTime"`
		EndTime   string `form:"endTime"`
		Page      int    `form:"page"`
		Limit     int    `form:"limit"`
	}
	err := c.Bind(&backData)
	if err != nil {
		response.ResponseError(c, 400)
		return
	}
	//如果传来的数值为空
	if backData.StartTime == "" && backData.EndTime == "" {
		starList, err = starService.QStarClass(1, backData.Page, backData.Limit)
		if err != nil {
			response.ResponseErrorWithMsg(c, 400, "未找到班级之星")
			return
		}
	} else {
		starList, err = starService.SelTimeStar(backData.StartTime, backData.EndTime, 1, backData.Page, backData.Limit)
		if err != nil {
			response.ResponseErrorWithMsg(c, 400, "未找到班级之星")
			return
		}
	}
	if starList == nil {
		starList = []models.StarStu{}
	}
	data := map[string]any{
		"starlist": starList,
	}
	response.ResponseSuccess(c, data)
}

// BackStarGrade 返回前台年级成长之星
func BackStarGrade(c *gin.Context) {
	//返回前端数据
	var starList []models.StarStu
	//接收前端数据
	var backData struct {
		StartTime string `form:"startTime"`
		EndTime   string `form:"endTime"`
		Page      int    `form:"page"`
		Limit     int    `form:"limit"`
	}
	err := c.Bind(&backData)
	if err != nil {
		response.ResponseError(c, 400)
		return
	}
	//如果传来的数值为空
	if backData.StartTime == "" && backData.EndTime == "" {
		starList, err = starService.QStarClass(2, backData.Page, backData.Limit)
		if err != nil {
			response.ResponseErrorWithMsg(c, 400, "未找到年级之星")
			return
		}
	} else {
		starList, err = starService.SelTimeStar(backData.StartTime, backData.EndTime, 2, backData.Page, backData.Limit)
		if err != nil {
			response.ResponseErrorWithMsg(c, 400, "未找到年级之星")
			return
		}
	}
	if starList == nil {
		starList = []models.StarStu{}
	}
	data := map[string]any{
		"starlist": starList,
	}
	response.ResponseSuccess(c, data)
}

// BackStarCollege 返回前台院级成长之星
func BackStarCollege(c *gin.Context) {
	//返回前端数据
	var starList []models.StarStu
	//接收前端数据
	var backData struct {
		StartTime string `form:"startTime"`
		EndTime   string `form:"endTime"`
		Page      int    `form:"page"`
		Limit     int    `form:"limit"`
	}
	err := c.Bind(&backData)
	if err != nil {
		response.ResponseError(c, 400)
		return
	}
	//如果传来的数值为空
	if backData.StartTime == "" && backData.EndTime == "" {
		starList, err = starService.QStarClass(3, backData.Page, backData.Limit)
		if err != nil {
			response.ResponseErrorWithMsg(c, 400, "未找到院级之星")
			return
		}
	} else {
		starList, err = starService.SelTimeStar(backData.StartTime, backData.EndTime, 3, backData.Page, backData.Limit)
		if err != nil {
			response.ResponseErrorWithMsg(c, 400, "未找到院级之星")
			return
		}
	}
	if starList == nil {
		starList = []models.StarStu{}
	}
	data := map[string]any{
		"starlist": starList,
	}
	response.ResponseSuccess(c, data)
}

// ChangeStatus 修改是否可以再次推选
func ChangeStatus(c *gin.Context) {
	token := token2.NewToken(c)
	user, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
	}
	err := mysql.UpdateOne(user.Username)
	if err != nil {
		response.ResponseErrorWithMsg(c, 400, "修改状态失败")
		return
	}
}

// BackTime 返回前端最大年份及最小年份
func BackTime(c *gin.Context) {
	maxtime, mintime, err := mysql.SelTime()
	if err != nil {
		response.ResponseError(c, 400)
		return
	}
	data := map[string]string{
		"maxDate": maxtime,
		"minDate": mintime,
	}
	response.ResponseSuccess(c, data)
}

// BackName 返回前端已经推选过的名字
func BackName(c *gin.Context) {
	token := token2.NewToken(c)
	user, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
	}
	role, err := token.GetRole()
	stuName, err := starService.BackNameData(user.Username, role)
	if err != nil {
		response.ResponseError(c, 401)
		return
	}
	data := map[string]any{
		"list": stuName,
	}
	response.ResponseSuccess(c, data)
}
