package mysql

// SelRoleMessage 获取角色id,名称,状态码
func SelRoleMessage(id int) (string, string, error) {
	var role string
	var code string
	err := DB.Table("casbin_rule").Where("id = ?", id).Select("v1").Scan(&role).Error
	err = DB.Table("casbin_rule").Where("id = ?", id).Select("v0").Scan(&code).Error
	if err != nil {
		return "", "", err
	}
	return role, code, nil
}

// SelRoleId 查询代表角色的id切片
func SelRoleId() ([]int, error) {
	var id []int
	err := DB.Table("casbin_rule").Where("ptype = ?", "g").Select("id").Scan(&id).Error
	if err != nil {
		return nil, err
	}
	return id, nil
}

// SelMRole 模糊查询角色字段
func SelMRole(role string) (bool, error) {
	var count int64
	err := DB.Table("casbin_rule").Where("v1 LIKE ?", role+"%").Count(&count).Error

	// 检查是否有错误发生
	if err != nil {
		// 处理错误，例如记录日志
		return false, err // 如果发生错误，返回false
	}

	// 如果count大于0，表示找到了记录
	if count > 0 {
		return true, nil
	}

	// 如果没有找到记录，返回false
	return false, nil
}
