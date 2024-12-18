package userService

import (
	"studentGrow/dao/mysql"
)

// BVerify 后台登录验证
func BVerify(username, password string) bool {
	number, err := mysql.SelPassword(username, password)
	if err != nil || number != 1 {
		return false
	}
	return true
}

// BVerifyExit 验证用户是否为管理员
func BVerifyExit(username string) bool {
	number, err := mysql.SelIfexit(username)
	if err != nil || number != 1 {
		return false
	}
	return true
}

// BVerifyBan 验证用户是否被封
func BVerifyBan(username string) (bool, error) {
	ban, err := mysql.SelBan(username)
	if err != nil {
		return false, err
	}
	return ban, nil
}
