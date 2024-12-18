package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"studentGrow/controller/RoleController"
	"studentGrow/controller/growth"
	"studentGrow/controller/login"
	"studentGrow/controller/menuController"
	"studentGrow/dao/mysql"
	"studentGrow/models/casbinModels"
	"studentGrow/utils/middleWare"
	"studentGrow/utils/token"
)

func RoutesXue(router *gin.Engine) {
	casbinService, err := casbinModels.NewCasbinService(mysql.DB)
	if err != nil {
		fmt.Println("Setup models.NewCasbinService()  err")
	}
	//验证码登录不鉴权
	user := router.Group("user")
	{
		//1.像前端返回验证码
		user.POST("/code", login.RidCode)
		//2.后台登录
		user.POST("/hlogin", login.HLogin)
		//前台登录
		user.POST("/qlogin", login.QLogin)
	}
	//获取注册天数
	tokenUser := router.Group("user")
	tokenUser.Use(token.AuthMiddleware())
	{
		tokenUser.POST("/register/day", login.RegisterDay)
	}

	//casbin鉴权
	userLoginAfter := router.Group("user")
	userLoginAfter.Use(token.AuthMiddleware(), middleWare.NewCasbinAuth(casbinService))
	{
		//班级管理员首页
		userLoginAfter.POST("/fpage/class", login.FPageClass)
		//年级管理员首页
		userLoginAfter.POST("/fpage/grade", login.FPageGrade)
		//学院管理员首页
		userLoginAfter.POST("/fpage/college", login.FPageCollege)
		//超级管理员首页
		userLoginAfter.POST("/fpage/superman", login.FPageCollege)
		//首页柱状图
		userLoginAfter.GET("/fpage/pillar", login.Pillar)
		//获取登陆者的全部信息
		userLoginAfter.GET("/message", menuController.HeadRoute)
	}
	//前台展示成长之星
	showStar := router.Group("star")
	//showStar.Use(token.AuthMiddleware())
	{
		showStar.GET("/class_star", growth.BackStarClass)
		showStar.GET("/grade_star", growth.BackStarGrade)
		showStar.GET("/college_star", growth.BackStarCollege)
		showStar.GET("/time", growth.BackTime)
	}
	//后台成长之星
	elected := router.Group("star")
	elected.Use(token.AuthMiddleware(), middleWare.NewCasbinAuth(casbinService))
	{
		//成长之星退选时展示的表格
		elected.GET("/select", growth.Search)
		//班级管理员推选
		elected.POST("/elected/class", growth.ElectClass)
		//年级管理员推选
		elected.POST("/elected/grade", growth.ElectGrade)
		//学院管理员推选
		elected.POST("/elected/college", growth.ElectCollege)
		//院级管理员公布
		elected.POST("/public/college", growth.PublicStar)
		//搜索第几届成长之星的接口
		elected.GET("/termStar", growth.StarPub)
		//修改角色状态判断是否可以继续推选
		elected.POST("/change_disabled", growth.ChangeStatus)
		//展示已经推选过的成长之星
		elected.GET("/is_elected", growth.BackName)
	}
	//前端侧边栏(鉴权)
	sidebar := router.Group("sidebar")
	sidebar.Use(token.AuthMiddleware(), middleWare.NewCasbinAuth(casbinService))
	{
		sidebar.GET("/message", menuController.MenuSide)
	}

	//菜单管理 casbin鉴权
	menu := router.Group("menuManage")
	menu.Use(token.AuthMiddleware(), middleWare.NewCasbinAuth(casbinService))
	{
		//菜单初始化
		menu.GET("/init", menuController.MenuMangerInit)
		//添加菜单
		menu.POST("/newelyBuilt", menuController.AddMenu)
		//删除菜单
		menu.POST("/delete", menuController.DeleteMenu)
		//搜索菜单
		menu.GET("/selectInfo", menuController.SearchMenu)
		//编辑菜单
		menu.POST("/edit", menuController.UpdateMenu)
		//菜单列表
		menu.GET("/list", menuController.MenuList)
	}
	//角色管理
	role := router.Group("role")
	role.Use(token.AuthMiddleware(), middleWare.NewCasbinAuth(casbinService))
	{
		role.GET("/list", RoleController.RoleList)
		role.GET("/permission", RoleController.ShowMenu)
		role.POST("/update", RoleController.UpdateRoleMenu)
	}
	router.POST("role/addRole", RoleController.AddRole)
}
