package menuController

import (
	"github.com/gin-gonic/gin"
	"studentGrow/dao/mysql"
	"studentGrow/models"
	"studentGrow/models/gorm_model"
	"studentGrow/pkg/response"
	service "studentGrow/service/permission"
)

// MenuMangerInit 初始化菜单
func MenuMangerInit(c *gin.Context) {
	menu, err := service.BuildMenuTree(0)
	if err != nil {
		response.ResponseError(c, 400)
		return
	}
	response.ResponseSuccess(c, menu)
}

// AddMenu 新增菜单
func AddMenu(c *gin.Context) {
	var fId int
	var backMenu models.Menu
	err := c.Bind(&backMenu)
	if err != nil || backMenu.Name == "" {
		response.ResponseErrorWithMsg(c, 400, "获取数据失败")
		return
	}
	//1.查找父id
	if backMenu.FatherMenu == "" {
		fId = 0
	} else {
		fId, err = mysql.SelMenuFId(backMenu.FatherMenu)
	}
	//2.增添数据模型
	menu := gorm_model.Menus{
		ParentId:      fId,
		TreePath:      "",
		Name:          backMenu.Name,
		Type:          backMenu.Type,
		RouteName:     backMenu.RouteName,
		Path:          backMenu.Path,
		Component:     backMenu.Component,
		Perm:          backMenu.Perm,
		Visible:       backMenu.Visible,
		Sort:          backMenu.Sort,
		Icon:          backMenu.Icon,
		Redirect:      backMenu.Redirect,
		Roles:         "",
		RequestUrl:    backMenu.RequestUrl,
		RequestMethod: backMenu.RequestMethod,
	}
	//查询是否存在
	ok, err := mysql.SelMenuExit(menu.Name)
	if !ok {
		response.ResponseSuccess(c, "菜单已存在")
		return
	}
	//3.新增菜单
	err = mysql.AddMenu(menu)
	if err != nil {
		response.ResponseError(c, 400)
		return
	}
	//4.新增路由参数
	//查询是否有路由参数
	if backMenu.Params != nil {
		//查询此菜单的id
		newid, err := mysql.SelMenuFId(backMenu.Name)
		for _, parm := range backMenu.Params {
			parm.MenuId = newid
			err = mysql.AddParam(parm)
			if err != nil {
				response.ResponseSuccess(c, "")
				return
			}
		}
	}
	response.ResponseSuccess(c, "")
}

// DeleteMenu 删除菜单
func DeleteMenu(c *gin.Context) {
	//接收前端数据
	var fromdata struct {
		MenuID int `json:"MenuId"`
	}
	err := c.Bind(&fromdata)
	if err != nil || fromdata.MenuID == 0 {
		response.ResponseError(c, 400)
		return
	}
	//删除menu中的数据
	err = mysql.DeleteMenu(fromdata.MenuID)
	if err != nil {
		response.ResponseError(c, 400)
		return
	}
	//删除路由参数
	err = mysql.DeleteParam(fromdata.MenuID)
	if err != nil {
		response.ResponseError(c, 400)
		return
	}
	response.ResponseSuccess(c, "")
}

// SearchMenu 搜索菜单
func SearchMenu(c *gin.Context) {
	//返回前端数据
	var menus []models.Menu
	//接收前端数据
	var inputdata struct {
		Input string `form:"input"`
	}
	err := c.Bind(&inputdata)
	if err != nil {
		response.ResponseError(c, 401)
		return
	}
	if inputdata.Input == "" {
		menus, err = service.BuildMenuTree(0)
	} else {
		//返回前端的数据
		menus, err = service.BuildMenu(inputdata.Input)
		if err != nil {
			response.ResponseError(c, 400)
			return
		}
	}

	response.ResponseSuccess(c, menus)
}

// UpdateMenu 编辑菜单
func UpdateMenu(c *gin.Context) {
	//接收前端数据
	var menu models.Menu
	err := c.Bind(&menu)
	if err != nil {
		response.ResponseErrorWithMsg(c, 401, "没有接收到数据")
		return
	}
	err = service.UpdateMenuData(menu)
	//删除之前的所有路由参数
	err = mysql.DeleteParam(menu.ID)
	//创建新的路由参数
	for _, param := range menu.Params {
		param.MenuId = menu.ID
		err := mysql.UpdateParam(param)
		if err != nil {
			response.ResponseError(c, 400)
			return
		}
	}
	if err != nil {
		response.ResponseError(c, 400)
		return
	}
	response.ResponseSuccess(c, "")
}

// MenuList 菜单下拉列表
func MenuList(c *gin.Context) {
	menulist, err := service.MenuList(0)
	if err != nil {
		response.ResponseError(c, 400)
		return
	}
	response.ResponseSuccess(c, menulist)
}
