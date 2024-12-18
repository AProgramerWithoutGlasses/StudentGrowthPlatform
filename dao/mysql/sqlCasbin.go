package mysql

import (
	"regexp"
	"strings"
	"studentGrow/models/gorm_model"
)

func SelMenuId(requestUrl, requestMethod string) (string, error) {
	var menuId string
	//获取最后一个/后的信息
	re := regexp.MustCompile("/([^/]+)$")
	matches := re.FindStringSubmatch(requestUrl)
	//判断是不是角色
	ok, err := SelMRole(matches[1])
	if ok {
		requestUrl = strings.Replace(requestUrl, matches[0], "", 1)
	}
	// 使用模糊查询，例如：查找所有以requestUrl开头的记录
	err = DB.Model(&gorm_model.Menus{}).Select("id").Where("request_url = ?", requestUrl).Where("request_method = ? ", requestMethod).Scan(&menuId).Error
	return menuId, err
}
