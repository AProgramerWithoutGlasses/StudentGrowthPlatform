package mysql

import (
	"studentGrow/models"
	"studentGrow/models/gorm_model"
)

// SelFMenu 查询权限下父ID是0的目录及菜单ID
func SelFMenu(role string) ([]int, error) {
	var DId []int
	err := DB.Model(&gorm_model.Menus{}).Where("type <> ?", 4).Where("type <> ?", 2).Where("roles LIKE ?", "%"+role+"%").Select("id").Scan(&DId).Error
	if err != nil {
		return nil, err
	}
	return DId, nil
}

// SelValueInt 查询menus菜单所有int类型的数据
func SelValueInt(id int, column string) (int, error) {
	var fId int
	err := DB.Model(&gorm_model.Menus{}).Where("id = ?", id).Select(column).Scan(&fId).Error
	if err != nil {
		return 0, err
	}
	return fId, nil
}

// SelValueString 查询menus菜单所有string类型的数据
func SelValueString(id int, column string) (string, error) {
	var value string
	err := DB.Model(&gorm_model.Menus{}).Where("id = ?", id).Select(column).Scan(&value).Error
	if err != nil {
		return "", err
	}
	return value, nil
}

// SelParamKeyVal 查询路由参数
func SelParamKeyVal(id int) ([]gorm_model.Param, error) {
	var param []gorm_model.Param
	err := DB.Model(&gorm_model.Param{}).Where("menu_id = ?", id).Scan(&param).Error
	if err != nil {
		return []gorm_model.Param{}, err
	}
	return param, nil
}

func SelIcon(id int) (string, error) {
	var icon string
	err := DB.Model(&gorm_model.Menus{}).Where("id = ?", id).Select("icon").Scan(&icon).Error
	if err != nil {
		return "", err
	}
	return icon, nil
}

// SelPerms 查询权限标识符
func SelPerms(role string) ([]string, error) {
	var perms []string
	var menuId []int
	err := DB.Table("casbin_rule").Where("v0 = ?", role).Select("v1").Scan(&menuId).Error
	err = DB.Model(&gorm_model.Menus{}).Where("type = ?", 2).Where("perm IS NOT NULL").Where("id IN (?)", menuId).Select("perm").Scan(&perms).Error
	if err != nil {
		return nil, err
	}
	return perms, nil
}

// SelOneDad 查询父Id为同一个值的所有数据
func SelOneDad(fid int) ([]gorm_model.Menus, error) {
	var menu []gorm_model.Menus
	err := DB.Model(&gorm_model.Menus{}).Where("parent_id = ?", fid).Scan(&menu).Error
	if err != nil {
		return nil, err
	}
	return menu, nil
}

// SelOneDadMenu 查询同一个父亲的目录和菜单
func SelOneDadMenu(fid int) ([]gorm_model.Menus, error) {
	var menu []gorm_model.Menus
	err := DB.Model(&gorm_model.Menus{}).Where("type <> ?", 4).Where("type <> ?", 2).Where("parent_id = ?", fid).Scan(&menu).Error
	if err != nil {
		return nil, err
	}
	return menu, nil
}

// SelMenuFId 根据名字查找id
func SelMenuFId(name string) (int, error) {
	var fId int
	err := DB.Model(&gorm_model.Menus{}).Where("name = ?", name).Select("id").Scan(&fId).Error
	if err != nil {
		return 0, err
	}
	return fId, nil
}

// SelMenuExit 查询菜单是否存在
func SelMenuExit(name string) (bool, error) {
	var count int64
	err := DB.Model(&gorm_model.Menus{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count != 0 {
		return false, nil
	}
	return true, nil
}

// AddMenu 新增菜单
func AddMenu(menu gorm_model.Menus) error {
	err := DB.Create(&menu).Error
	if err != nil {
		return err
	}
	return nil
}

// AddParam 新增路由参数
func AddParam(param gorm_model.Param) error {
	err := DB.Create(&param).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteMenu 删除菜单
func DeleteMenu(id int) error {
	err := DB.Where("id = ?", id).Delete(&gorm_model.Menus{}).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteParam 删除路由参数
func DeleteParam(id int) error {
	err := DB.Where("menu_id = ?", id).Delete(&gorm_model.Param{}).Error
	if err != nil {
		return err
	}
	return nil
}

// SelMenuIds 模糊查询菜单id
func SelMenuIds(name string) ([]gorm_model.Menus, error) {
	var menus []gorm_model.Menus
	err := DB.Model(&gorm_model.Menus{}).Where("name LIKE ?", "%"+name+"%").Where("deleted_at IS NULL").Scan(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

// UpdateMenus 更新菜单
func UpdateMenus(newMenu models.Menu, parentId int) error {
	var menu gorm_model.Menus
	err := DB.Table("menus").Where("id = ?", newMenu.ID).Scan(&menu).Error
	if err != nil {
		return err
	}
	menu.ParentId = parentId
	menu.Name = newMenu.Name
	menu.Type = newMenu.Type
	menu.RouteName = newMenu.RouteName
	menu.Path = newMenu.Path
	menu.Component = newMenu.Component
	menu.Perm = newMenu.Perm
	menu.Visible = newMenu.Visible
	menu.Sort = newMenu.Sort
	menu.Icon = newMenu.Icon
	menu.Redirect = newMenu.Redirect
	menu.RequestUrl = newMenu.RequestUrl
	menu.RequestMethod = newMenu.RequestMethod
	err = DB.Save(&menu).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateParam 更新路由
func UpdateParam(newParam gorm_model.Param) error {
	err := DB.Table("params").Create(&newParam).Error
	if err != nil {
		return err
	}
	return nil
}

// SelRoleMenu 模糊查询权限角色是否有权限
func SelRoleMenu(role string, id int) (bool, error) {
	var number int64
	err := DB.Model(&gorm_model.Menus{}).Where("id = ?", id).Where("roles LIKE ?", "%"+role+"%").Count(&number).Error
	if err != nil || number == 0 {
		return false, err
	}
	return true, nil
}
