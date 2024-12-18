package stuManage

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"studentGrow/dao/mysql"
	"studentGrow/models/jrx_model"
	"studentGrow/pkg/response"
	"studentGrow/service"
	"studentGrow/utils/readMessage"
	token2 "studentGrow/utils/token"
)

// 用于存储查询参数

// QueryStuContro 查询学生信息
func QueryStuContro(c *gin.Context) {
	var (
		queryParmaStruct  jrx_model.QueryParmaStruct
		querySql          string
		queryAllStuNumber int
		ranges            string
	)

	token := token2.NewToken(c)
	user, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
		return
	}

	role, err := token.GetRole()

	// 加上日志利于核验role
	zap.L().Info("角色信息记录 role：",
		zap.String("role", role), // 添加一个名为custom_field的字符串字段
	)

	fmt.Println("token解析为：" + role)
	if err != nil {
		response.ResponseError(c, response.TokenError)
		zap.L().Error(err.Error())
		return
	}

	id, err := mysql.GetIdByUsername(user.Username)
	if err != nil {
		zap.L().Error(err.Error())
		response.ResponseError(c, response.ServerErrorCode)
		return
	}

	class, err := mysql.GetClassById(id)
	if err != nil {
		zap.L().Error(err.Error())
		response.ResponseError(c, response.ServerErrorCode)
		return
	}

	// 权限管理
	switch role {
	case "class":
		ranges = " and class = '" + class + "'"

	case "grade1":
		ranges = " and YEAR(plus_time) = 2024"

	case "grade2":
		ranges = " and YEAR(plus_time) = 2023"

	case "grade3":
		ranges = " and YEAR(plus_time) = 2022"

	case "grade4":
		ranges = " and YEAR(plus_time) = 2021"

	case "college":

	case "superman":

	default:
	}

	// 接收请求数据
	stuMessage, err := readMessage.GetJsonvalue(c)
	if err != nil {
		fmt.Println("stuManage.QueryStuContro() readMessage.GetJsonvalue() err :", err)
	}

	offsetValue, err := stuMessage.GetInt("page")
	if err != nil {
		fmt.Println("page GetInt() err", err)
	}

	limitValue, err := stuMessage.GetInt("limit")
	if err != nil {
		fmt.Println("limit GetInt() err", err)
	}

	// 将请求数据整理到结构体
	queryParmaStruct = service.GetReqMes(stuMessage)
	// 获取sql语句
	querySql = service.CreateQuerySql(stuMessage, queryParmaStruct)
	querySql = querySql + ranges

	querySql = querySql + " ORDER BY username ASC"

	// 获取符合条件的所有学生，用于计算长度
	stuInfo, err := service.GetStuMesList(querySql) // 所有学生数据
	if err != nil {
		zap.L().Error("mysql.GetStuMesList(querySql) failed", zap.Error(err))
		response.ResponseError(c, response.ServerErrorCode)
		return
	}

	// 获取所有符合条件的学生数量
	queryAllStuNumber = len(stuInfo)

	// 重置sql语句中的分页部分
	whereSqlIndex := strings.Index(querySql, "limit")
	if whereSqlIndex != -1 {
		afterWhere := querySql[:whereSqlIndex]
		querySql = afterWhere
	}

	// limit 分页查询语句的拼接
	querySql = querySql + " limit " + strconv.Itoa(limitValue) + " offset " + strconv.Itoa((offsetValue-1)*limitValue)

	// 获取符合条件的当页学生
	stuPageInfo, err := service.GetStuMesList(querySql) // 当页学生数据
	if err != nil {
		zap.L().Error("mysql.GetStuMesList(queryPageSql) failed", zap.Error(err))
		response.ResponseError(c, response.ServerErrorCode)
		return
	}
	yearStructSlice := service.GetYearStructSlice()
	classStructSlice, err := service.GetClassStructSlice()
	if err != nil {
		response.ResponseErrorWithMsg(c, response.ServerErrorCode, err.Error())
		zap.L().Error("stuManage.QueryStuContro() mysql.GetDiffClass() failed", zap.Error(err))
	}

	// 响应结构体的初始化
	responseStruct := jrx_model.ResponseStruct{
		Role:            role,
		Year:            yearStructSlice,
		Class:           classStructSlice,
		StuInfo:         stuPageInfo,
		AllStudentCount: queryAllStuNumber,
	}

	// 响应数据
	response.ResponseSuccess(c, responseStruct)
}
