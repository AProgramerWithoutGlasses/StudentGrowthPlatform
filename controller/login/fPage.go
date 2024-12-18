package login

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/dao/mysql"
	"studentGrow/pkg/response"
	"studentGrow/service/permission"
	token2 "studentGrow/utils/token"
	"time"
)

// FPageClass 返回前端后台首页数据--班级管理员
func FPageClass(c *gin.Context) {
	//查询用户账号
	token := token2.NewToken(c)
	user, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
	}
	//把一个班的用户id查出来
	uidslice, err := mysql.SelUid(user.Class)

	//班级总帖数
	article_total, err := service.ArticleData(uidslice)
	if err != nil {
		fmt.Println("FPage ArticleDataClass err", err)
		response.ResponseError(c, 400)
		return
	}

	//今日班级贴子数
	today_article_total, err := service.NarticleDataClass(uidslice)
	if err != nil {
		fmt.Println("FPage NarticleDataClass err", err)
		response.ResponseError(c, 400)
		return
	}

	//昨日班级帖子数跟今日的比率
	article_ratio, err := service.ArticleDataClassRate(uidslice, today_article_total)
	if err != nil {
		fmt.Println("FPage ArticleDataClassRate err", err)
		response.ResponseError(c, 400)
		return
	}

	//今日访客人数
	today_visitor_total, err := service.TodayVictor()
	//今日跟昨日的对比
	visitor_ratio, err := service.VictorRate(today_visitor_total)
	if err != nil {
		fmt.Println("FPage VictorRate err", err)
		response.ResponseError(c, 400)
		return
	}

	//人员总数
	user_total := len(uidslice)

	//学生总数
	student_total, err := mysql.SelStudent()
	if err != nil {
		fmt.Println("FPage SelStudent err", err)
		response.ResponseError(c, 400)
		return
	}

	//教师总数
	teacher_total, err := mysql.SelTeacher()
	if err != nil {
		response.ResponseError(c, 400)
		return
	}

	//班级帖子总赞数
	upvote_amount, err := service.LikeAmount(uidslice)

	//班级帖子总阅读数
	article_read_total := service.ReadAmount(uidslice)

	//柱状图的数据
	tagname, count, err := service.PillarData()
	if err != nil {
		response.ResponseError(c, 400)
		return
	}
	chartOption := map[string]any{
		"xAxis": map[string]any{
			"type": "category",
			"data": tagname,
		},
		"yAxis": map[string]string{
			"type": "value",
		},
		"series": map[string]any{
			"data": count,
			"type": "bar",
		},
	}
	data := map[string]any{
		"article_total":       article_total,
		"today_article_total": today_article_total,
		"article_ratio":       article_ratio,
		"today_visitor_total": today_visitor_total,
		"visitor_ratio":       visitor_ratio,
		"user_total":          user_total,
		"student_total":       student_total,
		"teacher_total":       teacher_total,
		"upvote_amount":       upvote_amount,
		"article_read_total":  article_read_total,
		"chartOption":         chartOption,
	}
	response.ResponseSuccess(c, data)
}

// FPageGrade 返回前端后台首页数据--年级管理员
func FPageGrade(c *gin.Context) {
	var uidSlice []int
	var err error
	//拿到角色
	token := token2.NewToken(c)
	role, err := token.GetRole()
	nowdata := time.Now()
	switch role {
	case "grade1":
		uidSlice, _, err = mysql.SelGradeId(nowdata, -1)
		if err != nil {
			response.ResponseError(c, 400)
			return
		}
	case "grade2":
		uidSlice, _, err = mysql.SelGradeId(nowdata, -2)
		if err != nil {
			response.ResponseError(c, 400)
			return
		}
	case "grade3":
		uidSlice, _, err = mysql.SelGradeId(nowdata, -3)
		if err != nil {
			response.ResponseError(c, 400)
			return
		}
	case "grade4":
		uidSlice, _, err = mysql.SelGradeId(nowdata, -4)
		if err != nil {
			response.ResponseError(c, 400)
			return
		}
	}
	//帖子总数
	article_total, err := service.ArticleData(uidSlice)
	if err != nil {
		response.ResponseError(c, 400)
		return
	}

	//今日帖子总数
	today_article_total, err := service.NarticleDataClass(uidSlice)
	if err != nil {
		response.ResponseError(c, 400)
		return
	}

	//昨日新帖跟今天的比率
	article_ratio, err := service.ArticleDataClassRate(uidSlice, today_article_total)
	if err != nil {
		response.ResponseError(c, 400)
		return
	}

	//今日访客数
	today_visitor_totale, err := service.TodayVictor()

	//访客比率
	visitor_ratio, err := service.VictorRate(today_visitor_totale)
	if err != nil {
		response.ResponseError(c, 400)
		return
	}

	//人员总数
	user_total := len(uidSlice)

	//学生总数
	student_total, err := mysql.SelStudent()
	if err != nil {
		response.ResponseError(c, 400)
		return
	}

	//教师总数
	teacher_total, err := mysql.SelTeacher()
	if err != nil {
		response.ResponseError(c, 400)
		return
	}

	//总赞数
	upvote_amount, err := service.LikeAmount(uidSlice)

	//总阅读数
	article_read_total := service.ReadAmount(uidSlice)

	//柱状图的数据
	tagname, count, err := service.PillarData()
	if err != nil {
		response.ResponseError(c, 400)
	}
	chartOption := map[string]any{
		"xAxis": map[string]any{
			"type": "category",
			"data": tagname,
		},
		"yAxis": map[string]string{
			"type": "value",
		},
		"series": map[string]any{
			"data": count,
			"type": "bar",
		},
	}
	data := map[string]any{
		"article_total":       article_total,
		"today_article_total": today_article_total,
		"article_ratio":       article_ratio,
		"today_visitor_total": today_visitor_totale,
		"visitor_ratio":       visitor_ratio,
		"user_total":          user_total,
		"student_total":       student_total,
		"teacher_total":       teacher_total,
		"upvote_amount":       upvote_amount,
		"article_read_total":  article_read_total,
		"chartOption":         chartOption,
	}

	response.ResponseSuccess(c, data)
}

// FPageCollege 学院管理员和超级管理员
func FPageCollege(c *gin.Context) {
	//查询所有用户的id及username
	uidslice, _, err := mysql.SelCollegeId()
	if err != nil {
		response.ResponseError(c, 400)
		return
	}
	//总帖数
	article_total, err := service.ArticleData(uidslice)
	if err != nil {
		response.ResponseError(c, 400)
		return
	}

	//今日总帖数
	today_article_total, err := service.NarticleDataClass(uidslice)
	if err != nil {
		response.ResponseError(c, 400)
		return
	}

	//昨日帖子跟今天的比率
	article_ratio, err := service.ArticleDataClassRate(uidslice, today_article_total)
	if err != nil {
		response.ResponseError(c, 400)
		return
	}

	//查询今日访客数
	today_visitor_total, err := service.TodayVictor()
	if err != nil {
		response.ResponseError(c, 400)
		return
	}

	//今天和昨天访客比率
	visitor_ratio, err := service.VictorRate(today_visitor_total)
	if err != nil {
		response.ResponseError(c, 400)
		return
	}

	//人员总数
	user_total := len(uidslice)

	//学生总数
	student_total, err := mysql.SelStudent()
	if err != nil {
		response.ResponseError(c, 400)
		return
	}

	//教师总数
	teacher_total, err := mysql.SelTeacher()
	if err != nil {
		response.ResponseError(c, 400)
		return
	}

	//总赞数
	upvote_amount, err := service.LikeAmount(uidslice)
	if err != nil {
		response.ResponseError(c, 400)
		return
	}

	//总阅读数
	article_read_total := service.ReadAmount(uidslice)

	//柱状图的数据
	tagname, count, err := service.PillarData()
	if err != nil {
		response.ResponseError(c, 400)
		return
	}
	chartOption := map[string]any{
		"xAxis": map[string]any{
			"type": "category",
			"data": tagname,
		},
		"yAxis": map[string]string{
			"type": "value",
		},
		"series": map[string]any{
			"data": count,
			"type": "bar",
		},
	}
	data := map[string]any{
		"article_total":       article_total,
		"today_article_total": today_article_total,
		"article_ratio":       article_ratio,
		"today_visitor_total": today_visitor_total,
		"visitor_ratio":       visitor_ratio,
		"user_total":          user_total,
		"student_total":       student_total,
		"teacher_total":       teacher_total,
		"upvote_amount":       upvote_amount,
		"article_read_total":  article_read_total,
		"chartOption":         chartOption,
	}

	response.ResponseSuccess(c, data)
}

// Pillar 柱状图(首页)
func Pillar(c *gin.Context) {
	//接收前端传来的data值
	var tagname []string
	var count []int
	var Num struct {
		Date string `form:"date"`
	}
	err := c.Bind(&Num)
	if err != nil {
		response.ResponseError(c, 400)
		return
	}
	if Num.Date == "" {
		tagname, count, _ = service.PillarData()
	} else {
		tagname, count, _ = service.PillarDataTime(Num.Date)
	}

	chartOption := map[string]any{
		"xAxis": map[string]any{
			"type": "category",
			"data": tagname,
		},
		"yAxis": map[string]string{
			"type": "value",
		},
		"series": map[string]any{
			"data": count,
			"type": "bar",
		},
	}
	data := map[string]any{
		"chartOption": chartOption,
	}
	response.ResponseSuccess(c, data)
}
