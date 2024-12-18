package service

import (
	"studentGrow/dao/mysql"
	"studentGrow/models"
)

func RoleData() ([]models.RoleList, error) {
	var RoleList []models.RoleList
	ids, err := mysql.SelRoleId()
	if err != nil {
		return nil, err
	}
	for _, id := range ids {
		role, code, err := mysql.SelRoleMessage(id)
		if err != nil {
			return nil, err
		}
		roleList := models.RoleList{
			Id:       id,
			RoleName: role,
			RoleCode: code,
		}
		RoleList = append(RoleList, roleList)
	}
	return RoleList, nil
}
