package mysql

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	model "studentGrow/models/gorm_model"
	"studentGrow/utils/timeConverter"
	"time"
)

// QueryClassByUsername 通过用户名查找班级
func QueryClassByUsername(username string) (string, error) {
	var class string
	if err := DB.Model(&model.User{}).Select("class").Where("username = ?", username).First(&class).Error; err != nil {
		zap.L().Error("QueryClassByUsername() dao.mysql.sql_user_nzx err=", zap.Error(err))
		return "", err
	}
	return class, nil
}

func QueryUserIdByUsername(username string) (int, error) {
	var user model.User
	if err := DB.Select("id").Where("username = ?", username).First(&user).Error; err != nil {
		zap.L().Error("QueryClassByUsername() dao.mysql.sql_user_nzx err=", zap.Error(err))
		return -1, err
	}
	return int(user.ID), nil
}

// QueryPlusTimeByUsername 通过用户名查找入学时间
func QueryPlusTimeByUsername(username string) (*time.Time, error) {
	var user model.User
	if err := DB.Select("plus_time").Where("username = ?", username).First(&user).Error; err != nil {
		zap.L().Error("QueryGradeByUsername() dao.mysql.sql_user_nzx err=", zap.Error(err))
		return nil, err
	}
	return &user.PlusTime, nil
}

// QueryUserByAdvancedFilter 高级筛选用户(年级、班级、姓名)
func QueryUserByAdvancedFilter(grade int, class []string, name string) (*gorm.DB, error) {
	year, err := timeConverter.GetEnrollmentYear(grade)
	if err != nil {
		zap.L().Error("QueryClassByUsername() dao.mysql.sql_user_nzx err=", zap.Error(err))
		return nil, err
	}
	query := DB.Model(&model.User{}).
		Where("plus_time BETWEEN ? AND ? AND class IN ? AND name LIKE ?",
			fmt.Sprintf("%d-01-01", year.Year()), fmt.Sprintf("%d-12-31", year.Year()), class, fmt.Sprintf("%%%s%%", name))
	if query.Error != nil {
		zap.L().Error("QueryArticleByAdvancedFilter() dao.mysql.sql_user_nzx err=", zap.Error(err))
		return nil, err
	}
	return query, nil
}

//func QueryUserIdByUsername(username string) (int, error) {
//	var user model.User
//	if err := DB.Model(&model.User{}).Select("id").Where("username = ?", username).First(&user).Error; err != nil {
//		zap.L().Error("QueryUserIdByUsername() dao.mysql.sql_user_nzx err=", zap.Error(err))
//		return -1, err
//	}
//	return int(user.ID), nil
//}

// QueryAllUserId 查询所有用户的id
func QueryAllUserId() ([]uint, error) {
	var ids []uint
	if err := DB.Model(&model.User{}).Pluck("id", &ids).Error; err != nil {
		zap.L().Error("AddManagerMsg() dao.mysql.sql_msg err=", zap.Error(err))
		return nil, err
	}

	return ids, nil
}

func QueryUserByUserId(uid int) (*model.User, error) {
	var user model.User
	if err := DB.Where("id = ?", uid).First(&user).Error; err != nil {
		zap.L().Error("QueryUserByUserId() dao.mysql.sql_msg err=", zap.Error(err))
		return nil, err
	}
	return &user, nil
}
