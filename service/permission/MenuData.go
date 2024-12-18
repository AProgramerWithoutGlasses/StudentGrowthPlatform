package service

import (
	"studentGrow/dao/mysql"
	"studentGrow/models"
)

// MenuIdClass 管理员所有的菜单
func MenuIdClass(role string) ([]models.Sidebar, error) {
	//返回前端的切片
	var Menu []models.Sidebar
	//1.查询权限下所有的菜单和目录
	DId, err := mysql.SelFMenu(role)
	if err != nil {
		return nil, err
	}
	for _, id := range DId {
		//1.查询父id
		Fid, err := mysql.SelValueInt(id, "parent_id")
		//2.查询路由路径
		path, err := mysql.SelValueString(id, "path")
		//3.查询component
		component, err := mysql.SelValueString(id, "component")
		//4.查询跳转路径
		redirect, err := mysql.SelValueString(id, "redirect")
		//5.查询名字
		name, err := mysql.SelValueString(id, "name")
		//6.查询是否可见
		visible, err := mysql.SelValueInt(id, "visible")
		//7.查询图标
		icon, err := mysql.SelIcon(id)
		//8.查询路由名称
		routeName, err := mysql.SelValueString(id, "route_name")
		//9.查询菜单下所属参数的键和值
		params, err := mysql.SelParamKeyVal(id)
		if err != nil {
			return nil, err
		}
		mesa := models.Message{
			Name:    name,
			Visible: visible,
			Icon:    icon,
		}
		sidebar := models.Sidebar{
			Id:        id,
			ParentId:  Fid,
			Path:      path,
			RouteName: routeName,
			Component: component,
			Redirect:  redirect,
			Meta:      mesa,
			Params:    params,
		}
		Menu = append(Menu, sidebar)
	}
	return Menu, nil
}

// BuildMenuTree 递归查询所有的子菜单
func BuildMenuTree(parentID int) ([]models.Menu, error) {
	//定义返回的menu切片
	var backMenu []models.Menu
	//查询父id是参数的菜单
	menus, err := mysql.SelOneDad(parentID)
	if err != nil {
		return nil, err
	}
	//循环遍历这个切片
	for i := range menus {
		//查询父菜单的名字
		fatherName, err := mysql.SelValueString(parentID, "name")
		//查询这个菜单或者目录的参数切片
		params, err := mysql.SelParamKeyVal(int(menus[i].ID))
		//返回前端的切片中的一个对象
		menu := models.Menu{
			ID:            int(menus[i].ID),
			ParentId:      parentID,
			Name:          menus[i].Name,
			Type:          menus[i].Type,
			RouteName:     menus[i].RouteName,
			Path:          menus[i].Path,
			Perm:          menus[i].Perm,
			Redirect:      menus[i].Redirect,
			Visible:       menus[i].Visible,
			Sort:          menus[i].Sort,
			FatherMenu:    fatherName,
			RequestUrl:    menus[i].RequestUrl,
			RequestMethod: menus[i].RequestMethod,
			Params:        params,
		}
		children, err := BuildMenuTree(int(menus[i].ID))
		if err != nil {
			return nil, err
		}
		menu.Children = children
		backMenu = append(backMenu, menu)
	}
	return backMenu, nil
}

// BuildMenu 查询id对应的菜单及子菜单
func BuildMenu(name string) ([]models.Menu, error) {
	var backMenu []models.Menu

	menus, err := mysql.SelMenuIds(name)
	if err != nil {
		return nil, err
	}
	for i := range menus {
		//查询父菜单的名字
		fatherId, err := mysql.SelValueInt(i, "parent_id")
		fatherName, err := mysql.SelValueString(fatherId, "name")
		//查询这个菜单或者目录的参数切片
		params, err := mysql.SelParamKeyVal(i)
		//返回前端的切片中的一个对象
		menu := models.Menu{
			ID:            int(menus[i].ID),
			ParentId:      fatherId,
			Name:          menus[i].Name,
			Type:          menus[i].Type,
			RouteName:     menus[i].RouteName,
			Path:          menus[i].Path,
			Perm:          menus[i].Perm,
			Redirect:      menus[i].Redirect,
			Visible:       menus[i].Visible,
			Sort:          menus[i].Sort,
			FatherMenu:    fatherName,
			RequestUrl:    menus[i].RequestUrl,
			RequestMethod: menus[i].RequestMethod,
			Params:        params,
		}
		children, err := BuildMenuTree(int(menus[i].ID))
		if err != nil {
			return nil, err
		}
		menu.Children = children
		backMenu = append(backMenu, menu)
	}
	return backMenu, nil
}

// UpdateMenuData 更新菜单及菜单列表
func UpdateMenuData(menu models.Menu) error {
	fId, err := mysql.SelMenuFId(menu.FatherMenu)
	if err != nil {
		return err
	}
	//更新菜单
	err = mysql.UpdateMenus(menu, fId)
	//更新路由参数
	if err != nil {
		return err
	}
	return nil
}

// MenuList 菜单下拉列表
func MenuList(parentID int) ([]models.MenuList, error) {
	//定义返回的menu切片
	var backMenu []models.MenuList
	//查询父id是参数的菜单
	menus, err := mysql.SelOneDadMenu(parentID)
	if err != nil {
		return nil, err
	}
	//循环遍历这个切片
	for i := range menus {
		//返回前端的切片中的一个对象
		menulist := models.MenuList{
			Name:  menus[i].Name,
			Value: menus[i].Name,
		}
		children, err := MenuList(int(menus[i].ID))
		if err != nil {
			return nil, err
		}
		menulist.Children = children
		backMenu = append(backMenu, menulist)
	}
	return backMenu, nil
}

// RoleMenuTree 递归查询所有的子菜单并标注是否有权限
func RoleMenuTree(role string, parentID int) ([]models.Menu, error) {
	//定义返回的menu切片
	var backMenu []models.Menu
	//查询父id是参数的菜单
	menus, err := mysql.SelOneDad(parentID)
	if err != nil {
		return nil, err
	}
	//循环遍历这个切片
	for i := range menus {
		//查询父菜单的名字
		fatherName, err := mysql.SelValueString(parentID, "name")
		//查询这个菜单或者目录的参数切片
		params, err := mysql.SelParamKeyVal(int(menus[i].ID))
		//查询是否有权限
		status, err := mysql.SelRoleMenu(role, int(menus[i].ID))
		//返回前端的切片中的一个对象
		menu := models.Menu{
			ID:            int(menus[i].ID),
			ParentId:      parentID,
			Name:          menus[i].Name,
			Type:          menus[i].Type,
			RouteName:     menus[i].RouteName,
			Path:          menus[i].Path,
			Perm:          menus[i].Perm,
			Redirect:      menus[i].Redirect,
			Visible:       menus[i].Visible,
			Sort:          menus[i].Sort,
			FatherMenu:    fatherName,
			RequestUrl:    menus[i].RequestUrl,
			RequestMethod: menus[i].RequestMethod,
			Params:        params,
			Status:        status,
		}
		children, err := RoleMenuTree(role, int(menus[i].ID))
		if err != nil {
			return nil, err
		}
		menu.Children = children
		backMenu = append(backMenu, menu)
	}
	return backMenu, nil
}
