package service

// 对用户权限进行修改
//func UpdateRolePermissions(userId int, role string, newFunctions []string) (bool, error) {
//	//对比权限大小,cId1为当前用户的，cId2是用户修改的，默认权限搞得等级低
//	cId1, err := mysql.SelCasId(userId)
//	if err != nil {
//		fmt.Println("UpdateRolePermissions SelCasId(userId) err:", err)
//	}
//	cId2, err := mysql.SelCode(role)
//	if err != nil {
//		fmt.Println("UpdateRolePermissions SelCode(role) err:", err)
//	}
//	// 如果需要按数值比较，先转换字符串为整数
//	num1, err := strconv.Atoi(cId1)
//	if err != nil {
//		return false, err
//	}
//
//	num2, err := strconv.Atoi(cId2)
//	if err != nil {
//		return false, err
//	}
//
//	// 现在按数值比较
//	if num1 > num2 {
//		// 如果 cId2 的权限级别大于 cId1，返回错误
//		return false, errors.New("角色权限不够无法更改")
//	} else {
//		// 开始事务
//		tx := mysql.DB.Begin()
//		defer func() {
//			if r := recover(); r != nil {
//				tx.Rollback()
//			}
//		}()
//
//		// 删除旧权限
//		if err := tx.Where("v0 = ?", role).Delete(gorm_model.CasbinRule{}).Error; err != nil {
//			tx.Rollback()
//			return false, err
//		}
//
//		// 创建新权限
//		for _, functionName := range newFunctions {
//			var menuId int
//			if err := tx.Table("menus").Select("id").Where("name = ?", functionName).Scan(&menuId).Error; err != nil {
//				tx.Rollback()
//				return false, err
//			}
//			casbinRule := gorm_model.CasbinRule{V0: role, V1: strconv.Itoa(menuId)}
//			if err := tx.Create(&casbinRule).Error; err != nil {
//				tx.Rollback()
//				return false, err
//			}
//		}
//
//		// 提交事务
//		if err := tx.Commit().Error; err != nil {
//			return false, err
//		}
//
//		return true, nil
//	}
//}

// 添加管理员
func AddAdmin() {

}
