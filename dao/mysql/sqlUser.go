package mysql

import (
	"database/sql"
	"fmt"
	"studentGrow/models/gorm_model"
	"time"
)

// SelPassword 根据用户名和密码查询用户是否存在
func SelPassword(username, password string) (int64, error) {
	var number int64
	err := DB.Model(&gorm_model.User{}).Select("password").Where("deleted_at IS NULL").Where("username = ?", username).Where("password = ?", password).Count(&number).Error
	return number, err
}

// SelCasId 根据用户id查询对应角色id
func SelCasId(username string) (string, error) {
	var code string
	err := DB.Model(&gorm_model.UserCasbinRules{}).Select("casbin_cid").Where("c_username = ?", username).Scan(&code).Error
	return code, err
}

// SelRole 根据角色id查询角色
func SelRole(id string) (string, error) {
	var role string
	err := DB.Table("casbin_rule").Select("v1").Where("v0 = ?", id).Scan(&role).Error
	return role, err
}

// SelClass 根据学号获取班级
func SelClass(username string) (string, error) {
	var class sql.NullString
	err := DB.Model(&gorm_model.User{}).Select("class").Where("username = ?", username).Scan(&class).Error
	if err != nil {
		return "", err
	}
	return class.String, nil
}

// SelIfexit 查找用户是否是管理员
func SelIfexit(username string) (int64, error) {
	var number int64
	err := DB.Model(&gorm_model.UserCasbinRules{}).Where("deleted_at IS NULL").Where("c_username = ?", username).Count(&number).Error
	if err != nil {
		return 0, err
	}
	return number, nil
}

// SelHead 查找用户头像
func SelHead(username string) (string, error) {
	var headshot sql.NullString
	result := DB.Model(&gorm_model.User{}).Where("username = ? ", username).Select("head_shot").Scan(&headshot)
	if result.Error != nil {
		return "", result.Error
	}
	if result.RowsAffected == 0 {
		// 没有找到记录，返回空字符串
		return "", nil
	}
	return headshot.String, nil
}

// SelBan 查询是否被ban
func SelBan(username string) (bool, error) {
	var ban bool
	err := DB.Model(&gorm_model.User{}).Where("username = ?", username).Select("ban").Scan(&ban).Error
	if err != nil {
		return false, err
	}
	return ban, err
}

// CreateUser 记录用户登录
func CreateUser(username string, id int) error {
	var userLogin gorm_model.UserLoginRecord
	userLogin.Username = username
	userLogin.UserID = uint(id)
	err := DB.Create(&userLogin).Error
	if err != nil {
		return err
	}
	return nil
}

// SelEndTime 查询解禁时间
func SelEndTime(username string) (time.Time, error) {
	var data time.Time
	err := DB.Model(&gorm_model.User{}).Where("username = ?", username).Select("user_ban_end_time").Scan(&data).Error
	return data, err
}

// UpdateBan 解禁
func UpdateBan(username string) error {
	var user gorm_model.User
	user.Ban = false
	err := DB.Model(&gorm_model.User{}).Where("username = ?", username).Updates(map[string]interface{}{"ban": false}).Error
	return err
}

// SelExit 查看用户是否存在
func SelExit(username string) (bool, error) {
	var number int64
	err := DB.Model(&gorm_model.User{}).Where("username = ?", username).Count(&number).Error
	if err != nil || number != 1 {
		return false, err
	}
	return true, nil
}

// IfTeacher 查看用户是不是老师
func IfTeacher(username string) (bool, error) {
	var number int64
	err := DB.Model(&gorm_model.User{}).Where("username = ?", username).Where("identity = ?", "老师").Count(&number).Error
	if err != nil {
		return false, err
	} else if number != 1 {
		return false, nil
	}
	return true, nil
}

func SelUser(username string) (gorm_model.User, error) {
	var user gorm_model.User
	err := DB.First(&user, "username = ?", username).Error
	if err != nil {
		fmt.Println(err)
		return gorm_model.User{}, err
	}
	return user, nil
}
