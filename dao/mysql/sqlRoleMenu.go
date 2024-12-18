package mysql

func SelOldRole(role string) ([]string, error) {
	var menuId []string
	err := DB.Table("casbin_rule").Where("v0 = ?", role).Select("v1").Scan(&menuId).Error
	return menuId, err
}
